"""Health and readiness endpoints."""
from fastapi import APIRouter

from app.clients.ollama_client import OllamaClient
from app.core.logging import get_logger

logger = get_logger(__name__)

router = APIRouter()


@router.get("", summary="Liveness")
@router.get("/")
def health() -> dict:
    """Liveness: service is up."""
    return {"status": "ok", "service": "ats-ai-service"}


@router.get("/ready", summary="Readiness")
def ready() -> dict:
    """Readiness: optional Ollama check."""
    ollama_ok = OllamaClient().is_available()
    return {
        "status": "ok",
        "ollama_available": ollama_ok,
    }
