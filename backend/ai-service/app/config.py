"""Application configuration from environment."""
import os
from functools import lru_cache
from typing import Optional


@lru_cache
def get_settings() -> "Settings":
    return Settings()


class Settings:
    """Settings loaded from environment; use get_settings() in app."""

    # Server
    host: str = "0.0.0.0"
    port: int = 8000
    grpc_port: int = 50051
    grpc_enabled: bool = False

    # Go backend (for fetching job details, etc.)
    go_backend_url: str = "http://localhost:8080"
    go_backend_timeout_sec: int = 30

    # Ollama (local LLM)
    ollama_base_url: str = "http://localhost:11434"
    ollama_model: str = "llama3:8b"
    ollama_embed_model: str = "nomic-embed-text"
    ollama_embed_dim: int = 768
    ollama_timeout_sec: int = 50  # keep under Go client timeout (~60s) so fallback returns in sync; override with OLLAMA_TIMEOUT_SEC for slow models

    # Logging
    log_level: str = "INFO"

    def __init__(self) -> None:
        self.host = os.environ.get("HOST", self.host)
        self.port = int(os.environ.get("PORT", str(self.port)))
        self.grpc_port = int(os.environ.get("GRPC_PORT", str(self.grpc_port)))
        self.grpc_enabled = os.environ.get("GRPC_ENABLED", "false").lower() == "true"
        self.go_backend_url = os.environ.get("GO_BACKEND_URL", self.go_backend_url).rstrip("/")
        self.go_backend_timeout_sec = int(
            os.environ.get("GO_BACKEND_TIMEOUT_SEC", str(self.go_backend_timeout_sec))
        )
        self.ollama_base_url = os.environ.get("OLLAMA_BASE_URL", self.ollama_base_url).rstrip("/")
        self.ollama_model = os.environ.get("OLLAMA_MODEL", self.ollama_model)
        self.ollama_embed_model = os.environ.get("OLLAMA_EMBED_MODEL", self.ollama_embed_model)
        self.ollama_embed_dim = int(os.environ.get("OLLAMA_EMBED_DIM", str(self.ollama_embed_dim)))
        self.ollama_timeout_sec = int(
            os.environ.get("OLLAMA_TIMEOUT_SEC", str(self.ollama_timeout_sec))
        )
        self.log_level = os.environ.get("LOG_LEVEL", self.log_level)
