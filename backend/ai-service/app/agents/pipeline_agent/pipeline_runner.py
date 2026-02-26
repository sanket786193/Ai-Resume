"""Pipeline runner: aggregates agent outputs and runs rule-based scoring."""
from typing import Any, Dict, List

from app.agents.pipeline_agent.rule_based_scoring import RuleBasedScoringEngine


def _skill_names(skills: Any) -> List[str]:
    """Extract skill name strings from technical_skills (list of dicts) or plain list."""
    if not skills:
        return []
    if isinstance(skills, list):
        if not skills:
            return []
        first = skills[0]
        if isinstance(first, dict):
            return [str(s.get("skill_name", "")) for s in skills if s.get("skill_name")]
        return [str(s) for s in skills]
    return []


class PipelineRunner:
    def __init__(self) -> None:
        self.scoring_engine = RuleBasedScoringEngine()

    def run(self, input_data: Dict[str, Any]) -> Dict[str, float]:
        parsed_resume = input_data.get("parsed_resume") or {}
        parsed_jd = input_data.get("parsed_jd") or {}
        ats_result = input_data.get("ats_result") or {}
        nlp_result = input_data.get("nlp_result") or {}
        if isinstance(parsed_resume, str):
            try:
                import json
                parsed_resume = json.loads(parsed_resume) if parsed_resume.strip() else {}
            except Exception:
                parsed_resume = {}
        if isinstance(parsed_jd, str):
            try:
                import json
                parsed_jd = json.loads(parsed_jd) if parsed_jd.strip() else {}
            except Exception:
                parsed_jd = {}
        if isinstance(ats_result, str):
            try:
                import json
                ats_result = json.loads(ats_result) if ats_result.strip() else {}
            except Exception:
                ats_result = {}
        if isinstance(nlp_result, str):
            try:
                import json
                nlp_result = json.loads(nlp_result) if nlp_result.strip() else {}
            except Exception:
                nlp_result = {}
        required_skills = _skill_names(parsed_jd.get("technical_skills"))
        candidate_skills = parsed_resume.get("skills")
        if isinstance(candidate_skills, list) and candidate_skills and not isinstance(candidate_skills[0], str):
            candidate_skills = [str(s.get("skill_name", s) if isinstance(s, dict) else s) for s in candidate_skills]
        candidate_skills = candidate_skills or []
        scoring_input = {
            "required_skills": required_skills,
            "candidate_skills": list(candidate_skills),
            "required_experience": _extract_experience_years(parsed_jd),
            "candidate_experience": _extract_candidate_experience(parsed_resume),
            "ats": ats_result,
            "achievement_count": len(parsed_resume.get("experience") or []) * 2
            + len(parsed_resume.get("projects") or []),
            "language": nlp_result if isinstance(nlp_result, dict) else {},
            "role": {
                "required_keywords": list(parsed_jd.get("ats_keywords") or []),
                "matched_keywords": list(parsed_resume.get("skills") or []),
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


def _create_scoring_agent():
    """Create the rule-based scoring agent for the pipeline.
    Uses ADK BaseAgent so it runs in sequence and writes final_score to state.
    """
    from typing import AsyncGenerator
    from google.adk.agents import BaseAgent
    from google.adk.agents.invocation_context import InvocationContext
    from google.adk.events import Event
    from google.adk.events.event_actions import EventActions

    class RuleBasedScoringAgent(BaseAgent):
        """Workflow agent: reads pipeline state, runs rule-based scoring, writes final_score to state."""

        async def _run_async_impl(self, ctx: InvocationContext) -> AsyncGenerator[Event, None]:
            state = dict(ctx.session.state)
            result = PipelineRunner().run(state)
            score = result.get("final_score", 0.0)
            yield Event(
                author=self.name,
                invocation_id=ctx.invocation_id,
                actions=EventActions(state_delta={"final_score": score}),
            )

    return RuleBasedScoringAgent(
        name="RuleBasedScoringAgent",
        description="Rule-based ATS scoring from parsed resume, JD, ATS and NLP results.",
    )


ScoringAgent = _create_scoring_agent()
