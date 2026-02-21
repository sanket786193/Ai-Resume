-- Optional: run this only after pgvector is installed (e.g. CREATE EXTENSION vector; as superuser).
-- Not run by goose. Apply manually: psql $DSN -f scripts/optional_pgvector_apply.sql

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS resume_embeddings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
  embedding vector(768),
  model_version VARCHAR(64),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(resume_id, job_id)
);
CREATE INDEX IF NOT EXISTS idx_resume_embeddings_job_id ON resume_embeddings(job_id);
CREATE INDEX IF NOT EXISTS idx_resume_embeddings_candidate_id ON resume_embeddings(candidate_id);
CREATE INDEX IF NOT EXISTS idx_resume_embeddings_vector ON resume_embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 1);

CREATE TABLE IF NOT EXISTS job_embeddings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE UNIQUE,
  embedding vector(768),
  model_version VARCHAR(64),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_job_embeddings_vector ON job_embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 1);
