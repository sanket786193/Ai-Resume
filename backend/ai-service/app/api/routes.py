"""Register all API routers."""
from fastapi import APIRouter

from app.api.v1 import embed, health, parse, screen

api_router = APIRouter(prefix="/api/v1")

api_router.include_router(health.router, prefix="/health", tags=["health"])
api_router.include_router(screen.router, prefix="/screen", tags=["screen"])
api_router.include_router(parse.router, prefix="/parse", tags=["parse"])
api_router.include_router(embed.router, prefix="/embed", tags=["embed"])

# Root-level routes for Go backend compatibility (Go calls POST /screen, /parse, /embed; GET /health)
root_router = APIRouter()
root_router.include_router(health.router, prefix="/health", tags=["health"])
root_router.include_router(screen.router, prefix="/screen", tags=["screen"])
root_router.include_router(parse.router, prefix="/parse", tags=["parse"])
root_router.include_router(embed.router, prefix="/embed", tags=["embed"])
