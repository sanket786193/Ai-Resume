"""Embedding endpoint for pgvector: text -> vector. Best-effort when Ollama is down: returns empty embedding instead of 503."""
from fastapi import APIRouter

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
    """Return embedding vector for text (e.g. cleaned resume or JD). When Ollama is unreachable, returns 200 with empty embedding so parse/screen flow continues."""
    text = (req.text or "").strip()
    if not text:
        return EmbedResponse(embedding=[], model_version=None)
    try:
        embedding = ollama.embed(text)
        settings = get_settings()
        return EmbedResponse(
            embedding=embedding,
            model_version=settings.ollama_embed_model,
        )
    except OllamaUnavailableError as e:
        logger.warning("Embed skipped (Ollama unreachable): %s", e)
        return EmbedResponse(embedding=[], model_version=None)
