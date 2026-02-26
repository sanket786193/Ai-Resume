"""FastAPI application entry - modular, scalable AI service for Go backend."""
import threading
from contextlib import asynccontextmanager

from fastapi import FastAPI

from app.api.routes import api_router, root_router
from app.config import get_settings
from app.core.logging import get_logger, setup_logging

logger = get_logger(__name__)


def _start_grpc_server():
    """Run gRPC server in a background thread (for Go backend to call)."""
    try:
        from app.grpc_server import serve
        s = get_settings()
        serve(s.grpc_port)
    except Exception as e:
        logger.exception("gRPC server failed: %s", e)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Startup: init logging, optionally start gRPC server. Shutdown: cleanup."""
    setup_logging()
    logger.info("AI service starting")
    s = get_settings()
    if getattr(s, "grpc_enabled", False):
        t = threading.Thread(target=_start_grpc_server, daemon=True)
        t.start()
        logger.info("gRPC server thread started on port %s", s.grpc_port)
    yield
    logger.info("AI service shutting down")


def create_app() -> FastAPI:
    """Application factory for testability and clarity."""
    app = FastAPI(
        title="ATS AI Service",
        description="Resume screening and ranking; called by Go backend.",
        version="1.0.0",
        lifespan=lifespan,
    )
    app.include_router(root_router)
    app.include_router(api_router)
    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn
    s = get_settings()
    uvicorn.run(
        "app.main:app",
        host=s.host,
        port=s.port,
        reload=True,
    )
