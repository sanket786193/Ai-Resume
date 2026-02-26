"""Orchestrates resume screening: Ollama (55s timeout), then heuristic."""
import json
import re
from dataclasses import dataclass
from typing import List, Optional

from app.api.v1.schemas import JobRequirements
from app.clients.ollama_client import OllamaClient
from app.config import get_settings
from app.core.exceptions import OllamaUnavailableError
from app.core.logging import get_logger
from app.services.resume_parser import extract_skills_stub, extract_text_from_content

logger = get_logger(__name__)


@dataclass
class ScreeningResult:
    """Result of screening a resume against a job description (full AI evaluation)."""

    skill_match_score: float
    ranking_score: float
    qualified: bool
    ats_score: Optional[int] = None
    skill_match_pct: Optional[int] = None
    missing_skills: Optional[List[str]] = None
    experience_match: Optional[str] = None
    experience_warnings: Optional[List[str]] = None
    keyword_matches: Optional[List[str]] = None
    semantic_matches: Optional[List[str]] = None
    summary: Optional[str] = None
    model_version: Optional[str] = None


class ScreeningService:
    """Runs resume screening: Ollama (55s timeout), then heuristic fallback."""

    def __init__(self, ollama_client: Optional[OllamaClient] = None) -> None:
        self.ollama = ollama_client or OllamaClient()

    def screen(
        self,
        resume_content_or_path: str,
        job_description: str,
        vector_similarity: Optional[float] = None,
        job_requirements: Optional[JobRequirements] = None,
    ) -> ScreeningResult:
        """
        Screen resume against job description and optional job requirements.
        - resume_content_or_path: raw text or path (path: read file in production).
        - job_description: job description text.
        - job_requirements: optional skills, experience_level, qualification for matching.
        Returns scores in [0, 1]; qualified is a recommendation.
        """
        resume_text = resume_content_or_path.strip()
        if not resume_text:
            return self._fallback_result()

        job_text = (job_description or "").strip()
        if not job_text:
            job_text = "General role"
        # Enrich job text with requirements for Ollama when provided
        if job_requirements:
            parts = [job_text]
            if job_requirements.skills:
                parts.append("Required skills: " + ", ".join(job_requirements.skills))
            if job_requirements.experience_level and job_requirements.experience_level != "ANY":
                parts.append("Experience level: " + job_requirements.experience_level)
            if job_requirements.qualification:
                parts.append("Qualification: " + job_requirements.qualification)
            job_text = "\n\n".join(parts)

        # 1) Ollama (55s timeout; ats_score, missing_skills, summary)
        try:
            if self.ollama.is_available():
                return self._screen_with_ollama(resume_text, job_text, vector_similarity)
        except OllamaUnavailableError:
            logger.warning("Ollama unavailable (timeout or unreachable), using heuristic fallback")

        # 2) Heuristic (use required skills for overlap when provided)
        return self._fallback_heuristic(resume_text, job_text, job_requirements)

    def _screen_with_ollama(
        self, resume_text: str, job_text: str, vector_similarity: Optional[float] = None
    ) -> ScreeningResult:
        model = get_settings().ollama_model
        sim_note = f" Vector similarity: {vector_similarity:.2f}." if vector_similarity is not None else ""
        # Shorter system prompt and inputs for faster response (<60s total)
        system = (
            "ATS assistant. Output ONLY one JSON object, no markdown. Keys: ats_score (0-100), skill_match_pct (0-100), "
            "missing_skills (array), experience_match (Good/Fair/Low), experience_warnings (array), "
            "keyword_matches (array), semantic_matches (array), summary (3-5 sentences: experience, skills, fit), qualified (boolean)."
        )
        prompt = (
            f"Resume:\n{resume_text[:2500]}\n\nJD:\n{job_text[:1500]}{sim_note}\n\nJSON only:"
        )
        response = self.ollama.generate(prompt, system=system).strip()
        return self._parse_ollama_json_response(response, model)

    def _parse_ollama_json_response(self, response: str, model_version: str) -> ScreeningResult:
        # Strip markdown code block if present
        text = response.strip()
        if text.startswith("```"):
            text = text.split("\n", 1)[-1].rsplit("```", 1)[0].strip()
        try:
            data = json.loads(text)
            ats = int(data.get("ats_score", 0))
            skill_pct = int(data.get("skill_match_pct", 0))
            qual = bool(data.get("qualified", False))
            missing = data.get("missing_skills")
            if not isinstance(missing, list):
                missing = []
            missing = [str(x) for x in missing][:20]
            exp = str(data.get("experience_match", "")) or None
            exp_warnings = data.get("experience_warnings")
            if not isinstance(exp_warnings, list):
                exp_warnings = []
            exp_warnings = [str(x) for x in exp_warnings][:15]
            keyword_m = data.get("keyword_matches")
            if not isinstance(keyword_m, list):
                keyword_m = []
            keyword_m = [str(x) for x in keyword_m][:30]
            semantic_m = data.get("semantic_matches")
            if not isinstance(semantic_m, list):
                semantic_m = []
            semantic_m = [str(x) for x in semantic_m][:30]
            summary = str(data.get("summary", "")) or None
            return ScreeningResult(
                skill_match_score=max(0.0, min(1.0, ats / 100.0)),
                ranking_score=max(0.0, min(1.0, skill_pct / 100.0)),
                qualified=qual,
                ats_score=max(0, min(100, ats)),
                skill_match_pct=max(0, min(100, skill_pct)),
                missing_skills=missing or None,
                experience_match=exp,
                experience_warnings=exp_warnings or None,
                keyword_matches=keyword_m or None,
                semantic_matches=semantic_m or None,
                summary=summary,
                model_version=model_version,
            )
        except (json.JSONDecodeError, ValueError, TypeError) as e:
            logger.warning("Ollama JSON parse failed: %s", e)
        return self._fallback_result()

    def _fallback_heuristic(
        self,
        resume_text: str,
        job_text: str,
        job_requirements: Optional[JobRequirements] = None,
    ) -> ScreeningResult:
        """Simple keyword overlap when Ollama is not used. Uses required skills when provided."""
        resume_clean = extract_text_from_content(resume_text)
        job_clean = extract_text_from_content(job_text)
        resume_skills = set(s.lower() for s in extract_skills_stub(resume_clean))
        if job_requirements and job_requirements.skills:
            job_words = set(s.lower().strip() for s in job_requirements.skills if s and len(s) >= 2)
        else:
            job_words = set(re.findall(r"[a-zA-Z]{3,}", job_clean.lower()))
        keyword_matches_list = sorted(resume_skills & job_words)[:30]
        missing = sorted(job_words - resume_skills)[:20] if job_words else []
        overlap = len(resume_skills & job_words) / max(len(job_words), 1)
        score = min(1.0, overlap * 2.0)
        return ScreeningResult(
            skill_match_score=score,
            ranking_score=score * 0.9,
            qualified=score >= 0.3,
            keyword_matches=keyword_matches_list if keyword_matches_list else None,
            missing_skills=missing if missing else None,
        )

    def _fallback_result(self) -> ScreeningResult:
        return ScreeningResult(
            skill_match_score=0.5,
            ranking_score=0.5,
            qualified=False,
        )
