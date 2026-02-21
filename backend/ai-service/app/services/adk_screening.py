"""Run Google ADK screening pipeline and map result to Go contract."""
from typing import Any, Dict, Optional

from app.core.logging import get_logger

logger = get_logger(__name__)


def run_adk_pipeline(resume_text: str, jd_text: str) -> Optional[Dict[str, Any]]:
    """
    Run the ADK screening pipeline (resume -> jd -> nlp -> ats -> skm -> scoring).
    Returns pipeline result dict with 'final_score' if successful, else None.
    """
    try:
        from app.agents.pipeline_agent import screening_pipeline
    except Exception as e:
        logger.warning("ADK pipeline import failed: %s", e)
        return None

    try:
        result = screening_pipeline.run(
            {"resume_text": resume_text, "jd_text": jd_text}
        )
        return result if isinstance(result, dict) else None
    except Exception as e:
        logger.warning("ADK pipeline run failed: %s", e)
        return None


def adk_result_to_scores(result: Dict[str, Any]) -> tuple[float, float, bool]:
    """
    Map ADK pipeline result to (skill_match_score, ranking_score, qualified).
    final_score is 0-100; we normalize to 0-1 and set qualified >= 50.
    """
    final = float(result.get("final_score") or 0)
    normalized = max(0.0, min(1.0, final / 100.0))
    qualified = final >= 50
    return normalized, normalized, qualified
