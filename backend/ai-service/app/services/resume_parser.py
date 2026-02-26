"""Resume text extraction and structured parsing for PDF/plain text. Output is JSON-serializable for resume_parsed_data."""
import re
from typing import Any, Dict, List

from app.core.logging import get_logger

logger = get_logger(__name__)

# Email and phone patterns
EMAIL_RE = re.compile(r"[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+")
PHONE_RE = re.compile(
    r"[\+]?[(]?[0-9]{1,3}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}"
)

# Common skills for better extraction (subset; extend as needed)
COMMON_SKILLS = frozenset(
    {
        "python", "java", "javascript", "typescript", "react", "node", "node.js",
        "sql", "aws", "docker", "kubernetes", "git", "rest", "api", "machine learning",
        "ml", "data analysis", "agile", "scrum", "leadership", "communication",
        "problem solving", "go", "golang", "rust", "c++", "c#", ".net", "angular",
        "vue", "html", "css", "mongodb", "postgresql", "redis", "linux", "ci/cd",
        "terraform", "graphql", "microservices", "testing", "tdd", "figma",
        "excel", "tableau", "power bi", "jira", "confluence", "azure", "gcp",
    }
)

# Section headers (lowercase) that start experience / education / projects
EXP_HEADERS = ("experience", "work experience", "employment", "professional experience", "career")
EDU_HEADERS = ("education", "academic", "qualification", "degrees", "certifications")
PROJECT_HEADERS = ("projects", "key projects", "side projects", "open source")


def extract_text_from_content(content: str) -> str:
    """Normalize raw resume content: strip, collapse internal whitespace to single space."""
    if not content or not content.strip():
        return ""
    text = content.strip()
    text = re.sub(r"\s+", " ", text)
    return text


def extract_skills_stub(text: str) -> List[str]:
    """
    Extract skill-like tokens: ALL CAPS words, CamelCase terms, and known common skills.
    Returns deduplicated list (max 50) for JSON storage.
    """
    if not text:
        return []
    seen: set = set()
    result: List[str] = []
    # 1) Known common skills (from normalized text)
    lower = text.lower()
    for skill in COMMON_SKILLS:
        if skill in lower and skill not in seen:
            seen.add(skill)
            result.append(skill.title() if len(skill) > 2 else skill.upper())
    # 2) CamelCase tokens (e.g. JavaScript, TypeScript)
    for m in re.finditer(r"\b([A-Z][a-z]+(?:[A-Z][a-z]+)+)\b", text):
        tok = m.group(1)
        if len(tok) >= 2 and tok.lower() not in seen:
            seen.add(tok.lower())
            result.append(tok)
    # 3) Consecutive ALL CAPS words (e.g. MACHINE LEARNING)
    tokens = text.upper().split()
    for t in tokens:
        t_clean = re.sub(r"[^A-Z]", "", t)
        if len(t_clean) >= 2 and t_clean.isalpha() and t_clean.lower() not in seen:
            seen.add(t_clean.lower())
            result.append(t)
    return result[:50]


def _section_lines(lines: List[str], start_headers: tuple, stop_headers: tuple) -> List[str]:
    """Return lines belonging to a section from first start_header to next stop or end."""
    section: List[str] = []
    in_section = False
    for line in lines:
        lower = line.lower().strip()
        if any(lower.startswith(h) for h in start_headers):
            in_section = True
            continue
        if in_section and any(lower.startswith(h) for h in stop_headers):
            break
        if in_section and line.strip():
            section.append(line.strip())
    return section


def _parse_entries(lines: List[str], max_entries: int = 20) -> List[Dict[str, str]]:
    """
    Heuristic: split by lines that look like job title (short, often company below).
    Each entry: title, company, duration, description (all strings for JSON).
    """
    entries: List[Dict[str, str]] = []
    current: Dict[str, str] = {}
    desc_lines: List[str] = []
    for line in lines:
        line = line.strip()
        if not line:
            continue
        # Likely title line: short, no date pattern at start
        date_at_start = re.match(r"^(?:\d{4}|\d{1,2}/\d{1,2}|\w+\s+\d{4})", line)
        if not date_at_start and len(line) < 80 and " at " not in line.lower():
            if current and (current.get("title") or current.get("description")):
                current["description"] = " ".join(desc_lines)[:2000]
                entries.append(current)
                if len(entries) >= max_entries:
                    break
            current = {"title": line[:300], "company": "", "duration": "", "description": ""}
            desc_lines = []
            continue
        # Duration pattern (e.g. "2020 - 2022" or "Jan 2020 – Present")
        dur = re.match(
            r"^([\d\w\s\-–—/]+(?:\d{4}|present|current))(?:\s*[·\-]\s*(.*))?$",
            line,
            re.IGNORECASE,
        )
        if dur and current and not current.get("duration"):
            current["duration"] = dur.group(1).strip()[:100]
            if dur.lastindex and dur.lastindex >= 2 and dur.group(2):
                current["company"] = dur.group(2).strip()[:300]
            continue
        if current:
            if not current.get("company") and len(line) < 120:
                current["company"] = line[:300]
            else:
                desc_lines.append(line)
    if current and (current.get("title") or current.get("description")):
        current["description"] = " ".join(desc_lines)[:2000]
        entries.append(current)
    return entries[:max_entries]


def parse_resume_structured(raw_text: str) -> Dict[str, Any]:
    """
    Parse raw resume text into structured fields for resume_parsed_data.
    All values are JSON-serializable: str, list of str, or list of dict with str values.
    """
    if not raw_text or not raw_text.strip():
        return _empty_parsed()

    text = raw_text.strip()
    lines = [ln.strip() for ln in text.split("\n") if ln.strip()]

    # Contact
    emails = EMAIL_RE.findall(text)
    phones = PHONE_RE.findall(text)
    name = ""
    if lines:
        first = lines[0]
        if not (EMAIL_RE.search(first) or PHONE_RE.search(first)) and len(first) < 120:
            name = first[:200]

    # Sections (use line structure)
    exp_lines = _section_lines(lines, EXP_HEADERS, EDU_HEADERS + PROJECT_HEADERS)
    edu_lines = _section_lines(lines, EDU_HEADERS, EXP_HEADERS + PROJECT_HEADERS)
    proj_lines = _section_lines(lines, PROJECT_HEADERS, EXP_HEADERS + EDU_HEADERS)

    experience_entries = _parse_entries(exp_lines, 15)
    education_entries = _parse_entries(edu_lines, 10)
    # Projects as simple list of strings (description per project)
    projects = []
    for ln in proj_lines[:15]:
        if ln and len(ln) > 2:
            projects.append(ln[:500])

    skills = extract_skills_stub(text)

    # Single-string fallbacks for backward compatibility with screening/pipeline
    experience_str = " ".join(exp_lines)[:4000] if exp_lines else ""
    education_str = " ".join(edu_lines)[:2000] if edu_lines else ""

    return {
        "name": name,
        "email": emails[0] if emails else "",
        "phone": phones[0] if phones else "",
        "skills": skills,
        "experience": experience_str,
        "education": education_str,
        "experience_entries": [dict(e) for e in experience_entries],
        "education_entries": [dict(e) for e in education_entries],
        "projects": projects,
    }


def _empty_parsed() -> Dict[str, Any]:
    return {
        "name": "",
        "email": "",
        "phone": "",
        "skills": [],
        "experience": "",
        "education": "",
        "experience_entries": [],
        "education_entries": [],
        "projects": [],
    }


def cleaned_text_for_embedding(raw_text: str, parsed: Dict[str, Any]) -> str:
    """Build normalized, deduplicated text for embedding (stored in resume_parsed_data.cleaned_text)."""
    parts = []
    if parsed.get("name"):
        parts.append(parsed["name"].lower())
    if parsed.get("skills"):
        parts.append(" ".join(s.lower() for s in parsed["skills"]))
    if parsed.get("experience"):
        parts.append(extract_text_from_content(parsed["experience"]).lower())
    for e in parsed.get("experience_entries") or []:
        for v in (e.get("title"), e.get("company"), e.get("description")):
            if v:
                parts.append(extract_text_from_content(v).lower())
    if parsed.get("education"):
        parts.append(extract_text_from_content(parsed["education"]).lower())
    for e in parsed.get("education_entries") or []:
        for v in (e.get("title"), e.get("company"), e.get("description")):
            if v:
                parts.append(extract_text_from_content(v).lower())
    for p in parsed.get("projects") or []:
        if p:
            parts.append(extract_text_from_content(p).lower())
    if not parts:
        parts.append(extract_text_from_content(raw_text).lower())
    combined = " ".join(parts)
    combined = re.sub(r"\s+", " ", combined)
    return combined[:8000]
