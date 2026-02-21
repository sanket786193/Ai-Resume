"""Custom exceptions for the AI service."""


class AIServiceError(Exception):
    """Base exception for AI service errors."""

    def __init__(self, message: str, details: dict | None = None) -> None:
        super().__init__(message)
        self.message = message
        self.details = details or {}


class OllamaUnavailableError(AIServiceError):
    """Ollama service is not reachable or failed."""


class GoBackendError(AIServiceError):
    """Go backend API call failed."""
