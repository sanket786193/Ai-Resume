"""Orchestrates resume screening: ADK pipeline first, then Ollama/heuristic fallback."""
import json
import re
from dataclasses import dataclass
from typing import List, Optional

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
    """Runs resume screening: ADK pipeline first, then Ollama, then heuristic."""

    def __init__(self, ollama_client: Optional[OllamaClient] = None) -> None:
        self.ollama = ollama_client or OllamaClient()

    def screen(
        self,
        resume_content_or_path: str,
        job_description: str,
        vector_similarity: Optional[float] = None,
    ) -> ScreeningResult:
        """
        Screen resume against job description.
        - resume_content_or_path: raw text or path (path: read file in production).
        - job_description: job description text.
        Returns scores in [0, 1]; qualified is a recommendation.
        """
        resume_text = resume_content_or_path.strip()
        if not resume_text:
            return self._fallback_result()

        job_text = (job_description or "").strip()
        if not job_text:
            job_text = "General role"

        # 1) Try Google ADK pipeline (resume -> jd -> nlp -> ats -> skm -> rule-based scoring)
        try:
            from app.services.adk_screening import adk_result_to_scores, run_adk_pipeline

            adk_result = run_adk_pipeline(resume_text, job_text)
            if adk_result and "final_score" in adk_result:
                skill, rank, qual = adk_result_to_scores(adk_result)
                return ScreeningResult(
                    skill_match_score=skill,
                    ranking_score=rank,
                    qualified=qual,
                )
        except Exception as e:
            logger.warning("ADK screening failed, using fallback: %s", e)

        # 2) Ollama-only scoring (full ATS evaluation: ats_score, missing_skills, summary)
        try:
            if self.ollama.is_available():
                return self._screen_with_ollama(resume_text, job_text, vector_similarity)
        except OllamaUnavailableError:
            logger.warning("Ollama unavailable (timeout or unreachable), using fallback heuristic scoring")

        # 3) Heuristic
        return self._fallback_heuristic(resume_text, job_text)

    def _screen_with_ollama(
        self, resume_text: str, job_text: str, vector_similarity: Optional[float] = None
    ) -> ScreeningResult:
        model = get_settings().ollama_model
        sim_note = f" Vector similarity (resume-job): {vector_similarity:.2f}." if vector_similarity is not None else ""
        system = (
            "You are an ATS assistant. Evaluate the candidate resume against the job description. "
            "Reply with ONLY a single JSON object, no markdown or extra text. Keys: ats_score (0-100), "
            "skill_match_pct (0-100), missing_skills (array of required job skills/terms not found or weak in resume), "
            "experience_match (short string, e.g. Good/Fair/Low), experience_warnings (array of strings: specific experience mismatches, e.g. 'Years in X below requirement'), "
            "keyword_matches (array of job skills/terms that appear verbatim or near-verbatim in the resume), "
            "semantic_matches (array of job skills/requirements that are satisfied by resume meaning even if different words used), "
            "summary (2-3 sentences), qualified (boolean)."
        )
        prompt = (
            f"Resume excerpt:\n{resume_text[:4000]}\n\nJob description:\n{job_text[:2500]}\n\n{sim_note}\n\n"
            "Output JSON only:"
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

    def _parse_ollama_response(self, response: str) -> ScreeningResult:
        try:
            parts = response.split()
            if len(parts) >= 3:
                skill = float(parts[0])
                rank = float(parts[1])
                qual = parts[2] in ("1", "true", "yes")
                return ScreeningResult(
                    skill_match_score=max(0.0, min(1.0, skill)),
                    ranking_score=max(0.0, min(1.0, rank)),
                    qualified=qual,
                )
        except (ValueError, IndexError):
            pass
        return self._fallback_result()

    def _fallback_heuristic(self, resume_text: str, job_text: str) -> ScreeningResult:
        """Simple keyword overlap when Ollama is not used."""
        resume_clean = extract_text_from_content(resume_text)
        job_clean = extract_text_from_content(job_text)
        resume_skills = set(s.lower() for s in extract_skills_stub(resume_clean))
        job_words = set(re.findall(r"[a-zA-Z]{3,}", job_clean.lower()))
        keyword_matches_list = sorted(resume_skills & job_words)[:30]
        overlap = len(resume_skills & job_words) / max(len(job_words), 1)
        score = min(1.0, overlap * 2.0)
        return ScreeningResult(
            skill_match_score=score,
            ranking_score=score * 0.9,
            qualified=score >= 0.3,
            keyword_matches=keyword_matches_list if keyword_matches_list else None,
        )

    def _fallback_result(self) -> ScreeningResult:
        return ScreeningResult(
            skill_match_score=0.5,
            ranking_score=0.5,
            qualified=False,
        )
