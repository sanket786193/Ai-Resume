"""Launcher: run from repo root with `python main.py` or `uvicorn app.main:app`."""
import uvicorn

from app.config import get_settings

if __name__ == "__main__":
    s = get_settings()
    uvicorn.run(
        "app.main:app",
        host=s.host,
        port=s.port,
        reload=True,
    )
