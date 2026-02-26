"""HTTP client for the Go backend API."""
import httpx
from typing import Any, Optional

from app.config import get_settings
from app.core.exceptions import GoBackendError
from app.core.logging import get_logger

logger = get_logger(__name__)


class GoBackendClient:
    """Calls Go backend REST APIs (jobs, health, etc.)."""

    def __init__(
        self,
        base_url: Optional[str] = None,
        timeout_sec: Optional[int] = None,
    ) -> None:
        s = get_settings()
        self.base_url = (base_url or s.go_backend_url).rstrip("/")
        self.timeout = timeout_sec or s.go_backend_timeout_sec

    def _request(
        self,
        method: str,
        path: str,
        *,
        json: Optional[dict] = None,
        params: Optional[dict] = None,
    ) -> Any:
        url = f"{self.base_url}{path}"
        with httpx.Client(timeout=self.timeout) as client:
            try:
                r = client.request(method, url, json=json, params=params)
                r.raise_for_status()
                return r.json() if r.content else None
            except httpx.HTTPStatusError as e:
                logger.warning("Go backend error: %s %s -> %s", method, path, e.response.status_code)
                raise GoBackendError(
                    f"Go backend returned {e.response.status_code}",
                    details={"path": path, "status": e.response.status_code},
                ) from e
            except httpx.RequestError as e:
                logger.warning("Go backend request failed: %s", e)
                raise GoBackendError(f"Go backend unreachable: {e}") from e

    def health(self) -> dict:
        """GET /health."""
        return self._request("GET", "/health") or {}

    def get_job(self, job_id: str) -> Optional[dict]:
        """GET /jobs/:id - public job details (e.g. for screening context)."""
        return self._request("GET", f"/jobs/{job_id}")
