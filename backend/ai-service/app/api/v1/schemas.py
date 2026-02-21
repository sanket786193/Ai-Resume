"""Request/response schemas for API (match Go backend contract)."""
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class ScreenRequest(BaseModel):
    """POST /screen body - aligned with Go internal/ai/client."""

    resume_path_or_content: str = Field(..., description="Resume text or path")
    job_description: str = Field(..., description="Job description text")
    vector_similarity: Optional[float] = Field(None, description="Optional cosine similarity for context")


class ScreenResponse(BaseModel):
    """Response for /screen - full AI evaluation for ATS."""

    skill_match_score: float = Field(..., ge=0, le=1)
    ranking_score: float = Field(..., ge=0, le=1)
    qualified: bool = Field(...)
    ats_score: Optional[int] = Field(None, ge=0, le=100)
    skill_match_pct: Optional[int] = Field(None, ge=0, le=100)
    missing_skills: Optional[List[str]] = None
    experience_match: Optional[str] = None
    summary: Optional[str] = None
    model_version: Optional[str] = None


class ParseRequest(BaseModel):
    """POST /parse body."""

    resume_path_or_content: str = Field(..., description="Resume raw text or URL/path to fetch")


class ParseResponse(BaseModel):
    """Parsed resume: raw, structured, cleaned (best practice: keep separately)."""

    raw_text: str = Field(..., description="Raw extracted text")
    parsed: Dict[str, Any] = Field(..., description="Structured: name, email, phone, skills, experience, education")
    cleaned_text: str = Field(..., description="Normalized text for embedding")


class EmbedRequest(BaseModel):
    """POST /embed body."""

    text: str = Field(..., description="Text to embed (e.g. cleaned resume or JD)")


class EmbedResponse(BaseModel):
    """Embedding vector for pgvector."""

    embedding: List[float] = Field(..., description="Vector (dim from model, e.g. 768)")
    model_version: Optional[str] = None
