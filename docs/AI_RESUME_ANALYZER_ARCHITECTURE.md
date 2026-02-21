# AI Resume Analyzer – Architecture (Spec vs Current)

This doc maps the **recommended AI Resume Analyzer** spec to the **current codebase** and lists gaps and next steps.

---

## High-Level Architecture (Spec)

```
Candidate → Resume Upload → Resume Parsing → Structured Data + Embeddings
    → Vector DB (Similarity) → LLM Evaluation (ATS Score, Insights) → HR Dashboard
```

---

## 1. Resume Upload & Storage

| Spec | Current | Status |
|------|---------|--------|
| Upload resume securely | Go: `UploadResumePDF` (multipart) → Supabase Storage | Done |
| Validate file type & size | PDF only, max 100 MB | Done (DOCX not yet) |
| Store raw file | Supabase Storage; URL in DB | Done |
| Table: `resumes` (id, candidate_id, file_url, parsed_status, created_at) | `resumes`: id, candidate_id, file_name, storage_path, file_size, mime_type, created_at | Done (no `parsed_status` column) |

**Gaps:** DOCX support; optional `parsed_status` on `resumes` for UI (e.g. "parsed" / "pending").

---

## 2. Resume Parsing Engine

| Spec | Current | Status |
|------|---------|--------|
| Convert unstructured → structured JSON | Python ai-service: `POST /parse` (URL or raw text) | Done |
| Extract: name, email, phone, skills, experience, education, certifications | name, email, phone, skills, experience, education | Done (certifications not extracted) |
| Tools: pdfplumber, docx, rule-based + ML | PyPDF2 (PDF from URL), regex + heuristics in `resume_parser.py` | Done (docx not wired) |
| Output: structured + cleaned text for embedding | `ParseResponse`: raw_text, parsed (dict), cleaned_text | Done |
| Persist parsed data | Go stores in `resume_parsed_data` (raw_text, parsed_json, cleaned_text) after Apply | Done |

**Gaps:** DOCX extraction; certifications in parsed output; optional stronger parsing (e.g. pdfplumber) if needed.

---

## 3. Embedding & Vector Storage

| Spec | Current | Status |
|------|---------|--------|
| Resume text → embeddings | ai-service: `POST /embed` (Ollama, e.g. nomic-embed-text, 768-dim) | Endpoint exists |
| Job description → embeddings | Same `/embed` endpoint; not called from Go on job create/update | Partial |
| Store embeddings | — | Not implemented |
| Vector DB: pgvector on PostgreSQL | No pgvector table or similarity search | Missing |
| Table: `resume_embeddings` (resume_id, embedding vector(1536)) | — | Missing |

**Gaps:**

- **Go backend:** No call to AI `/embed` when resume is parsed or when job is saved.
- **DB:** No `resume_embeddings` (or similar) table with pgvector; no `job_embeddings` (or `jobs` embedding column).
- **Flow:** After parse, run embed on `cleaned_text` and store in vector table; on job create/update, embed JD and store.

---

## 4. Job Description Analyzer

| Spec | Current | Status |
|------|---------|--------|
| Normalize JD, extract must-have vs nice-to-have skills | — | Missing |
| Generate JD embeddings | Could use same `/embed`; not done | Missing |
| Output: mandatory_skills, experience_required | — | Missing |

**Gaps:** Dedicated JD analyzer (or reuse parse + embed) and a place to store JD embedding + optional structured fields (skills, experience_required).

---

## 5. ATS Scoring Engine (AI Brain)

| Spec | Current | Status |
|------|---------|--------|
| Local LLM (Ollama + LLaMA 3) | Screening uses Ollama (configurable model) | Done |
| Prompt: ATS score 0–100, skill match %, missing skills, strengths/summary | `ScreeningService`: Ollama prompt returns ats_score, skill_match_pct, missing_skills, experience_match, summary, qualified | Done |
| Contract: JSON response | ScreenResponse in ai-service; Go `AIScreenResult`; stored on `ats_records` | Done |
| Optional: pass vector_similarity into prompt | ScreenRequest has `vector_similarity`; screening uses it in prompt if provided | Done (caller doesn’t pass it yet) |

**Gaps:** Go does not compute or pass `vector_similarity` into the screen request (no vector search yet).

---

## 6. Ranking & Shortlisting

| Spec | Current | Status |
|------|---------|--------|
| Formula: `final_score = similarity*0.4 + ats*0.4 + experience*0.2` | — | Not implemented |
| Vector similarity (40%) | No vector DB → no similarity | Missing |
| ATS score (40%) | ats_score stored; not used in a combined rank | Partial |
| Experience match (20%) | experience_match stored as text; not normalized to 0–1 | Partial |

**Gaps:**

- Compute similarity (e.g. cosine) from pgvector between resume and job embeddings.
- Pass similarity into screen (or compute after), then compute `final_score` and persist (e.g. `ranking_score` or new column).
- HR list/sort by this final score.

---

## 7. HR Dashboard

| Spec | Current | Status |
|------|---------|--------|
| Job → Candidate list | HR Job Applicants page; ATS pipeline (all jobs) | Done |
| ATS score badges | Applicant cards and detail (Contact / Application / Resume / AI pages) | Done |
| Skill gap highlights | Applicant detail AI tab: missing_skills, summary, experience_match | Done |
| Resume preview | Resume file name/link only (no inline viewer) | Partial |
| Stack: React, shadcn/ui, TanStack Query | Used across HR and applicant flows | Done |

**Gaps:** Optional: resume preview (iframe/PDF viewer or link to storage); charting/analytics (e.g. score distribution, funnel).

---

## Summary: What Exists vs What’s Missing

| Module | Exists | Missing |
|--------|--------|--------|
| 1. Resume upload & storage | Yes (PDF, Supabase) | DOCX; optional parsed_status |
| 2. Resume parsing | Yes (parse + store in resume_parsed_data) | DOCX; certifications |
| 3. Embedding & vector storage | Embed API only | pgvector tables; store resume/JD embeddings; similarity search |
| 4. JD analyzer | — | JD normalization; JD embeddings + storage |
| 5. ATS scoring | Yes (Ollama, full feedback stored) | Pass vector_similarity from backend |
| 6. Ranking | ats_score/skill_match stored | final_score formula; similarity component; sort by rank |
| 7. HR dashboard | Applicants list + detail pages + AI tab | Resume preview; analytics/charts |

---

## Recommended Next Steps (Priority)

1. **Vector storage (pgvector)**  
   - Add extension and tables, e.g. `resume_embeddings (resume_id, embedding vector(768))` and job embeddings (e.g. `job_embeddings` or column on `jobs`).  
   - After parse, call `/embed` and store resume embedding; on job create/update, embed JD and store.

2. **Similarity in screening**  
   - When running AI screening, compute cosine similarity (resume vs job) from DB.  
   - Pass `vector_similarity` into existing `/screen` request so the prompt can use it.

3. **Ranking formula**  
   - Compute `final_score = similarity*0.4 + (ats_score/100)*0.4 + experience_weight*0.2` (define experience_weight from `experience_match` or years).  
   - Persist (e.g. `ranking_score` or new column) and use it for ordering in HR candidate lists and pipeline.

4. **JD analyzer (optional)**  
   - Add a small step (or reuse parse + embed) to extract mandatory_skills and experience_required from JD; store alongside JD embedding for future matching/UI.

5. **DOCX + parsed_status (optional)**  
   - Support DOCX in upload and parsing; add `parsed_status` on `resumes` if you want to show parsing state in the UI.

This keeps your current flow intact and adds the intelligence layer (embeddings, similarity, ranking) as in the spec.
