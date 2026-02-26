# Goose migrations

Run from repo root: `go run ./cmd/migrate -command=up` (or build `migrate` and run `./migrate -command=up`).

**Order (by filename):**

| File | Purpose |
|------|---------|
| `20260126045846_init.sql` | UUID extension |
| `20260126051004_auth.sql` | No-op (legacy) |
| `20260221000001_extensions_and_enums.sql` | `user_role`, `ats_status_enum`, `job_status_enum` |
| `20260221000002_ats_users.sql` | Users (HR/Candidate) |
| `20260221000003_ats_sessions.sql` | JWT refresh sessions |
| `20260221000004_jobs.sql` | Job postings |
| `20260221000005_candidates.sql` | Candidate profiles |
| `20260221000006_resumes.sql` | Resume metadata |
| `20260221000007_ats_records.sql` | Application records (ATS) |
| `20260221000008_interviews.sql` | Interviews |
| `20260221000009_offers.sql` | Offers |
| `20260221100000_drop_legacy_tables.sql` | Drop old `users`, `user_sessions`, etc. if present |
| `20260221100001_ats_pipeline.sql` | `resume_parsed_data`, ATS AI columns on `ats_records` |

**Optional (pgvector):** If your Postgres has the [pgvector](https://github.com/pgvector/pgvector) extension, run `scripts/optional_pgvector_apply.sql` manually to add `resume_embeddings` and `job_embeddings` for vector search. Otherwise omit it; the app works without it.
