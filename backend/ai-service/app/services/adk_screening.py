"""Run Google ADK screening pipeline and map result to Go contract.

ADK (google.adk) is optional. If not installed, screening falls back to Ollama/heuristic.
Flow: resume_text + jd_text -> (Resume -> JD -> NLP -> ATS -> SKM -> RuleBasedScoring) -> final_score in state.
"""
import asyncio
import uuid
from typing import Any, Dict, Optional

from app.core.logging import get_logger

logger = get_logger(__name__)

APP_NAME = "ScreeningApp"
USER_ID = "screen"


def _run_adk_pipeline_async(resume_text: str, jd_text: str) -> Optional[Dict[str, Any]]:
    """
    Run the ADK pipeline via Runner: create session with state, run_async, drain events, read final_score.
    """
    try:
        from google.adk.apps import App
        from google.adk.runners import Runner
        from google.adk.sessions import InMemorySessionService
        from google.genai.types import Content, Part
    except ImportError as e:
        logger.debug("ADK imports not available: %s", e)
        return None

    try:
        from app.agents.pipeline_agent import screening_pipeline
    except Exception as e:
        logger.debug("Pipeline import failed: %s", e)
        return None

    async def _run() -> Optional[Dict[str, Any]]:
        session_service = InMemorySessionService()
        session = await session_service.create_session(
            app_name=APP_NAME,
            user_id=USER_ID,
            state={
                "resume_text": resume_text,
                "jd_text": jd_text,
            },
            session_id=f"screen_{uuid.uuid4().hex[:12]}",
        )
        app = App(name=APP_NAME, root_agent=screening_pipeline)
        runner = Runner(app=app, session_service=session_service)
        new_message = Content(role="user", parts=[Part(text="Screen this resume against the job description.")])
        try:
            async for _ in runner.run_async(
                user_id=USER_ID,
                session_id=session.id,
                new_message=new_message,
            ):
                pass
        except Exception as e:
            logger.warning("ADK run_async failed: %s", e)
            return None
        updated = await session_service.get_session(
            app_name=APP_NAME, user_id=USER_ID, session_id=session.id
        )
        if not updated:
            return None
        final_score = updated.state.get("final_score")
        if final_score is None:
            return None
        return {"final_score": float(final_score)}

    try:
        return asyncio.run(_run())
    except Exception as e:
        logger.warning("ADK pipeline async run failed: %s", e)
        return None


def run_adk_pipeline(resume_text: str, jd_text: str) -> Optional[Dict[str, Any]]:
    """
    Run the ADK screening pipeline (resume -> JD -> NLP -> ATS -> SKM -> scoring).
    Returns pipeline result dict with 'final_score' if successful, else None.
    """
    return _run_adk_pipeline_async(resume_text, jd_text)


def adk_result_to_scores(result: Dict[str, Any]) -> tuple[float, float, bool]:
    """
    Map ADK pipeline result to (skill_match_score, ranking_score, qualified).
    final_score is 0-100; we normalize to 0-1 and set qualified >= 50.
    """
    final = float(result.get("final_score") or 0)
    normalized = max(0.0, min(1.0, final / 100.0))
    qualified = final >= 50
    return normalized, normalized, qualified
