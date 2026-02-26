# ATS AI Service (Python)

Resume screening and ranking for the Go ATS backend. Uses **Ollama** for LLM-based screening when available, with heuristic fallback; modular and scalable.

## Contract with Go backend

**HTTP**
- **POST /screen** – Body: `{ "resume_path_or_content": "...", "job_description": "...", "job_requirements": { "skills": ["Golang", "Python"], "experience_level": "EXPERIENCED", "qualification": "B.Tech" } }` (job_requirements optional).  
  Response: `{ "skill_match_score": 0.0–1.0, "ranking_score": 0.0–1.0, "qualified": true|false, ... }`. When job_requirements is provided, matching uses required skills, experience level, and qualification.
- **GET /health** – `{ "status": "ok" }`
- **GET /health/ready** – `{ "status": "ok", "ollama_available": true|false }`

**gRPC** (optional)
- **AIScreeningService.Screen** – Same request/response as POST /screen.  
  Run with `GRPC_ENABLED=true` so the gRPC server listens on `GRPC_PORT` (default 50051).  
  Go backend can use gRPC by setting `AI_USE_GRPC=true` and `AI_GRPC_TARGET=localhost:50051`.

## Virtual environment and run

```bash
cd backend/ai-service
python -m venv venv

# Windows
venv\Scripts\activate

# macOS/Linux
source venv/bin/activate

pip install -r requirements.txt
cp .env.example .env   # optional: edit .env (OLLAMA_MODEL, GO_BACKEND_URL)
python main.py
# or: uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

Service listens on port 8000. Set `AI_SERVICE_URL=http://localhost:8000` and `AI_ENABLED=true` in the Go backend so it calls this service for screening.

## Structure

- **app/main.py** – FastAPI app and lifespan
- **app/config.py** – Settings from env
- **app/api/** – Routes (health, screen); schemas match Go
- **app/services/** – Screening (Ollama → heuristic), resume_parser
- **app/clients/** – Ollama client, Go backend client
- **app/core/** – Logging, exceptions

## Screening flow

1. **Ollama**: Single LLM call (shorter resume/JD excerpts → JSON). Timeout 55s default; set `OLLAMA_TIMEOUT_SEC=60` if the model still times out.
2. **Heuristic**: Keyword/skill overlap (fallback when Ollama is unavailable or times out).

## Ollama

Run Ollama and pull a model, e.g. `ollama run llama3:8b`. Default in config: `llama3:8b`. For faster response use a smaller model: `OLLAMA_MODEL=llama3.2:3b` or `phi3:mini`.

**Timeouts:** The service waits up to `OLLAMA_TIMEOUT_SEC` (default 50s) for Ollama generate. This is kept below the Go backend’s AI client timeout (~60s) so that when Ollama is slow or times out, the service can return the heuristic fallback and the client still gets a response. For slower models, set `OLLAMA_TIMEOUT_SEC` higher and increase `AI_TIMEOUT_SEC` in the Go backend accordingly.

## gRPC (Go ↔ Python)

- **Proto**: Shared definitions live in `backend/proto/ats/screening.proto`.
- **Regenerate Go**: From `backend/`: `make proto-go` (needs `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`).
- **Regenerate Python**: From `backend/ai-service/` with venv: `python -m grpc_tools.protoc -I.. --python_out=app/grpc_gen --grpc_python_out=app/grpc_gen ../proto/ats/screening.proto`.
- **Run gRPC server**: Set `GRPC_ENABLED=true`; the same process will listen on HTTP (PORT) and gRPC (GRPC_PORT=50051). Or run `python -m app.grpc_server` in a separate process.
- **Go backend**: Set `AI_USE_GRPC=true` and `AI_GRPC_TARGET=localhost:50051` to call the AI service via gRPC instead of HTTP.
