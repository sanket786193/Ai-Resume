"""Generate Python gRPC code from proto. Run from ai-service: python scripts/gen_grpc.py"""
import subprocess
import sys
from pathlib import Path

# ai-service/scripts/gen_grpc.py -> ai-service
ROOT = Path(__file__).resolve().parent.parent
BACKEND = ROOT.parent  # backend
PROTO_ATS = BACKEND / "proto" / "ats"
OUT = ROOT / "app" / "grpc_gen" / "ats"

def main():
    OUT.mkdir(parents=True, exist_ok=True)
    # -I backend so that proto/ats/screening.proto is the path
    cmd = [
        sys.executable, "-m", "grpc_tools.protoc",
        f"-I{BACKEND}",
        f"--python_out={ROOT / 'app' / 'grpc_gen'}",
        f"--grpc_python_out={ROOT / 'app' / 'grpc_gen'}",
        "proto/ats/screening.proto",
    ]
    r = subprocess.run(cmd, cwd=BACKEND)
    if r.returncode != 0:
        sys.exit(r.returncode)
    print("Generated in app/grpc_gen/proto/ats/")

if __name__ == "__main__":
    main()
