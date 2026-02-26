"""Embedding endpoint for pgvector: text -> vector."""
from fastapi import APIRouter, HTTPException

from app.clients.ollama_client import OllamaClient
from app.config import get_settings
from app.core.exceptions import OllamaUnavailableError
from app.core.logging import get_logger

from .schemas import EmbedRequest, EmbedResponse

logger = get_logger(__name__)

router = APIRouter()
ollama = OllamaClient()


@router.post("", response_model=EmbedResponse, summary="Get embedding")
@router.post("/", response_model=EmbedResponse)
def embed_text(req: EmbedRequest) -> EmbedResponse:
    """Return embedding vector for text (e.g. cleaned resume or JD). Dimension from OLLAMA_EMBED_DIM (768 for nomic-embed-text)."""
    try:
        embedding = ollama.embed(req.text or "")
        settings = get_settings()
        return EmbedResponse(
            embedding=embedding,
            model_version=settings.ollama_embed_model,
        )
    except OllamaUnavailableError as e:
        logger.warning("Embed failed: %s", e)
        raise HTTPException(status_code=503, detail="Embedding service unavailable") from e
