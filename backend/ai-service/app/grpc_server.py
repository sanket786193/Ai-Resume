"""gRPC server for AIScreeningService (Go backend calls this)."""
import sys
from pathlib import Path

# Ensure app.grpc_gen.proto is importable as proto
_GRPC_GEN = Path(__file__).resolve().parent / "grpc_gen"
if str(_GRPC_GEN) not in sys.path:
    sys.path.insert(0, str(_GRPC_GEN))

import grpc
from concurrent import futures

from app.config import get_settings
from app.core.logging import get_logger
from app.services.screening import ScreeningService

logger = get_logger(__name__)


def _get_stub():
    from proto.ats import screening_pb2, screening_pb2_grpc
    return screening_pb2, screening_pb2_grpc


class AIScreeningServicer:
    """Implements AIScreeningService.Screen using ScreeningService."""

    def __init__(self):
        self._screening = ScreeningService()

    def Screen(self, request, context):
        from proto.ats import screening_pb2
        result = self._screening.screen(
            request.resume_path_or_content,
            request.job_description,
        )
        return screening_pb2.ScreenResponse(
            skill_match_score=result.skill_match_score,
            ranking_score=result.ranking_score,
            qualified=result.qualified,
        )


def serve(port: int | None = None):
    """Run the gRPC server (blocking)."""
    port = port or get_settings().grpc_port
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    screening_pb2, screening_pb2_grpc = _get_stub()
    screening_pb2_grpc.add_AIScreeningServiceServicer_to_server(
        AIScreeningServicer(), server
    )
    server.add_insecure_port(f"[::]:{port}")
    server.start()
    logger.info("gRPC server listening on port %s", port)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
