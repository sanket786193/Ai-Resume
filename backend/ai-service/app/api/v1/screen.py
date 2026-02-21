"""Resume screening endpoint - contract with Go backend."""
from fastapi import APIRouter, HTTPException

from app.core.exceptions import AIServiceError
from app.core.logging import get_logger
from app.services.screening import ScreeningService

from .parse import _fetch_resume_text
from .schemas import ScreenRequest, ScreenResponse

logger = get_logger(__name__)

router = APIRouter()
screening_service = ScreeningService()


@router.post("", response_model=ScreenResponse, summary="Screen resume")
@router.post("/", response_model=ScreenResponse)
def screen_resume(req: ScreenRequest) -> ScreenResponse:
    """
    Screen resume against job description. Returns full ATS evaluation:
    ats_score (0-100), skill_match_pct, missing_skills, experience_match, summary, model_version.
    If resume_path_or_content is a URL (e.g. Supabase), fetches and extracts PDF text first.
    """
    resume_input = (req.resume_path_or_content or "").strip()
    if resume_input.startswith(("http://", "https://")):
        resume_input = _fetch_resume_text(resume_input)
    try:
        result = screening_service.screen(
            resume_content_or_path=resume_input,
            job_description=req.job_description,
            vector_similarity=req.vector_similarity,
        )
        return ScreenResponse(
            skill_match_score=result.skill_match_score,
            ranking_score=result.ranking_score,
            qualified=result.qualified,
            ats_score=result.ats_score,
            skill_match_pct=result.skill_match_pct,
            missing_skills=result.missing_skills,
            experience_match=result.experience_match,
            experience_warnings=result.experience_warnings,
            keyword_matches=result.keyword_matches,
            semantic_matches=result.semantic_matches,
            summary=result.summary,
            model_version=result.model_version,
        )
    except AIServiceError as e:
        logger.exception("Screening error: %s", e.message)
        raise HTTPException(status_code=503, detail=e.message) from e
