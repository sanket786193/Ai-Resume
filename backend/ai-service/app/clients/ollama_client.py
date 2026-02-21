"""Client for local Ollama API."""
import httpx
from typing import Any, Optional

from app.config import get_settings
from app.core.exceptions import OllamaUnavailableError
from app.core.logging import get_logger

logger = get_logger(__name__)


class OllamaClient:
    """Calls local Ollama for completion/embedding."""

    def __init__(
        self,
        base_url: Optional[str] = None,
        model: Optional[str] = None,
        timeout_sec: Optional[int] = None,
    ) -> None:
        s = get_settings()
        self.base_url = (base_url or s.ollama_base_url).rstrip("/")
        self.model = model or s.ollama_model
        self.timeout = timeout_sec or s.ollama_timeout_sec

    def generate(self, prompt: str, system: Optional[str] = None) -> str:
        """Run completion; returns generated text."""
        url = f"{self.base_url}/api/generate"
        body: dict[str, Any] = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
        }
        if system:
            body["system"] = system
        try:
            with httpx.Client(timeout=self.timeout) as client:
                r = client.post(url, json=body)
                r.raise_for_status()
                data = r.json()
                return data.get("response", "")
        except httpx.RequestError as e:
            logger.warning("Ollama request failed: %s", e)
            raise OllamaUnavailableError(f"Ollama unreachable: {e}") from e
        except httpx.HTTPStatusError as e:
            logger.warning("Ollama error: %s", e.response.status_code)
            raise OllamaUnavailableError(f"Ollama returned {e.response.status_code}") from e

    def embed(self, text: str) -> list[float]:
        """Return embedding vector for text (e.g. nomic-embed-text, 768-dim)."""
        url = f"{self.base_url}/api/embeddings"
        s = get_settings()
        model = s.ollama_embed_model
        dim = s.ollama_embed_dim
        try:
            with httpx.Client(timeout=self.timeout) as client:
                r = client.post(
                    url,
                    json={"model": model, "prompt": text[:8192]},
                )
                r.raise_for_status()
                data = r.json()
                emb = data.get("embedding") or []
                if len(emb) != dim:
                    emb = (emb + [0.0] * dim)[:dim]
                return emb
        except httpx.RequestError as e:
            logger.warning("Ollama embed request failed: %s", e)
            raise OllamaUnavailableError(f"Ollama unreachable: {e}") from e
        except httpx.HTTPStatusError as e:
            logger.warning("Ollama embed error: %s", e.response.status_code)
            raise OllamaUnavailableError(f"Ollama returned {e.response.status_code}") from e

    def is_available(self) -> bool:
        """Check if Ollama is reachable."""
        try:
            with httpx.Client(timeout=5) as client:
                r = client.get(f"{self.base_url}/api/tags")
                return r.status_code == 200
        except Exception:
            return False
