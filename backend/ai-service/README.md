# ATS AI Service (Python)

Resume screening and ranking for the Go ATS backend. Uses **Google ADK** (Agent Development Kit) with Ollama when available; modular, readable, and scalable.

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
- **app/services/** – Screening (ADK pipeline → Ollama → heuristic), adk_screening, resume_parser
- **app/agents/** – Google ADK agents: resume_parsing, jd_agent, nlp_agent, ats_agent, skm_agent, pipeline_agent (rule-based scoring)
- **app/clients/** – Ollama client, Go backend client
- **app/core/** – Logging, exceptions

## Screening flow

1. **ADK pipeline** (if available): Runs via Runner + session with initial state `resume_text`, `jd_text`. Sequential agents: **Resume** (→ `parsed_resume`) → **JD** (→ `parsed_jd`) → **NLP** (→ `nlp_result`: grammar/language) → **ATS** (→ `ats_result`: structure/keywords) → **SKM** (→ `skm_result`) → **Rule-based scoring** (reads state, writes `final_score` 0–100). Result mapped to Go contract.
2. **Ollama only**: single LLM call for scores (fallback if ADK fails).
3. **Heuristic**: keyword overlap (fallback if Ollama unavailable).

## Ollama

Run Ollama and pull a model, e.g. `ollama run llama3:8b`. Set `OLLAMA_MODEL=llama3:8b` in `.env` if needed (default `llama3.2`). ADK uses LiteLlm with `ollama_chat/<model>`.

**Timeouts:** The service waits up to `OLLAMA_TIMEOUT_SEC` (default 50s) for Ollama generate. This is kept below the Go backend’s AI client timeout (~60s) so that when Ollama is slow or times out, the service can return the heuristic fallback and the client still gets a response. For slower models, set `OLLAMA_TIMEOUT_SEC` higher and increase `AI_TIMEOUT_SEC` in the Go backend accordingly.

## gRPC (Go ↔ Python)

- **Proto**: Shared definitions live in `backend/proto/ats/screening.proto`.
- **Regenerate Go**: From `backend/`: `make proto-go` (needs `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`).
- **Regenerate Python**: From `backend/ai-service/` with venv: `python -m grpc_tools.protoc -I.. --python_out=app/grpc_gen --grpc_python_out=app/grpc_gen ../proto/ats/screening.proto`.
- **Run gRPC server**: Set `GRPC_ENABLED=true`; the same process will listen on HTTP (PORT) and gRPC (GRPC_PORT=50051). Or run `python -m app.grpc_server` in a separate process.
- **Go backend**: Set `AI_USE_GRPC=true` and `AI_GRPC_TARGET=localhost:50051` to call the AI service via gRPC instead of HTTP.
