"""Pipeline runner: aggregates agent outputs and runs rule-based scoring."""
from typing import Any, Dict

from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm

from app.agents.pipeline_agent.rule_based_scoring import RuleBasedScoringEngine
from app.config import get_settings

_settings = get_settings()
dummy_model = LiteLlm(model=f"ollama_chat/{_settings.ollama_model}")


class PipelineRunner:
    def __init__(self) -> None:
        self.scoring_engine = RuleBasedScoringEngine()

    def run(self, input_data: Dict[str, Any]) -> Dict[str, float]:
        parsed_resume = input_data.get("parsed_resume") or {}
        parsed_jd = input_data.get("parsed_jd") or {}
        ats_result = input_data.get("ats_result") or {}
        nlp_result = input_data.get("nlp_result") or {}
        scoring_input = {
            "required_skills": parsed_jd.get("technical_skills") and [
                s.get("skill_name") for s in parsed_jd["technical_skills"] if s.get("skill_name")
            ] or [],
            "candidate_skills": parsed_resume.get("skills") or [],
            "required_experience": _extract_experience_years(parsed_jd),
            "candidate_experience": _extract_candidate_experience(parsed_resume),
            "ats": ats_result,
            "achievement_count": len(parsed_resume.get("experience") or []) * 2
            + len(parsed_resume.get("projects") or []),
            "language": nlp_result if isinstance(nlp_result, dict) else {},
            "role": {
                "required_keywords": parsed_jd.get("ats_keywords") or [],
                "matched_keywords": parsed_resume.get("skills") or [],
            },
        }
        return self.scoring_engine.calculate_final_score(scoring_input)


def _extract_experience_years(parsed_jd: Dict) -> float:
    exp = parsed_jd.get("experience_required_years") or {}
    if isinstance(exp, dict):
        min_y = exp.get("min")
        if min_y is not None:
            try:
                return float(min_y)
            except (TypeError, ValueError):
                pass
    return 0.0


def _extract_candidate_experience(parsed_resume: Dict) -> float:
    exp_list = parsed_resume.get("experience") or []
    total = 0.0
    for e in exp_list:
        d = e.get("duration") or ""
        # Heuristic: try to parse "2 years" etc.
        for part in d.replace(",", " ").split():
            try:
                total += float(part)
                break
            except ValueError:
                continue
    return total if total > 0 else 1.0


class ScoringAgent(Agent):
    def __init__(self) -> None:
        super().__init__(
            name="RuleBasedScoringAgent",
            model=dummy_model,
            description="Rule-based ATS scoring agent",
            instruction="Calculate the final ATS score based on inputs.",
            tools=[],
        )

    def run(self, input_data: Dict[str, Any]) -> Dict[str, float]:
        pipeline_runner = PipelineRunner()
        return pipeline_runner.run(input_data)
