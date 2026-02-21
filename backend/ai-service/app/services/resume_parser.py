"""Resume text extraction and skill parsing (stub; extend with PDF/OCR)."""
import re
from typing import Any, Dict, List

from app.core.logging import get_logger

logger = get_logger(__name__)

# Email and phone patterns
EMAIL_RE = re.compile(r"[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+")
PHONE_RE = re.compile(r"[\+]?[(]?[0-9]{1,3}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,4}[-\s\.]?[0-9]{1,9}")


def extract_text_from_content(content: str) -> str:
    """Normalize raw resume content (strip, collapse whitespace)."""
    if not content or not content.strip():
        return ""
    text = content.strip()
    text = re.sub(r"\s+", " ", text)
    return text


def extract_skills_stub(text: str) -> List[str]:
    """
    Stub: extract skill-like tokens (words in ALL CAPS or known patterns).
    Replace with NER or LLM-based extraction in production.
    """
    if not text:
        return []
    # Simple heuristic: consecutive caps words (e.g. "MACHINE LEARNING")
    tokens = text.upper().split()
    skills = []
    for i, t in enumerate(tokens):
        t_clean = re.sub(r"[^A-Z]", "", t)
        if len(t_clean) >= 2 and t_clean.isalpha():
            skills.append(t)
    return list(dict.fromkeys(skills))[:30]  # dedupe, cap


def parse_resume_structured(raw_text: str) -> Dict[str, Any]:
    """
    Parse raw resume text into structured fields: name, email, phone, skills, experience, education.
    Best practice: keep raw_text, structured, and cleaned_text separately.
    """
    if not raw_text or not raw_text.strip():
        return {
            "name": "",
            "email": "",
            "phone": "",
            "skills": [],
            "experience": "",
            "education": "",
        }
    text = raw_text.strip()
    emails = EMAIL_RE.findall(text)
    phones = PHONE_RE.findall(text)
    # First line often has name (no email/phone)
    lines = [l.strip() for l in text.split("\n") if l.strip()]
    name = lines[0] if lines else ""
    if name and (EMAIL_RE.search(name) or PHONE_RE.search(name)):
        name = ""
    skills = extract_skills_stub(text)
    # Heuristic: experience/education sections (look for headers)
    experience_parts = []
    education_parts = []
    in_exp = in_edu = False
    for line in lines[1:]:
        lower = line.lower()
        if "experience" in lower or "work" in lower or "employment" in lower:
            in_exp, in_edu = True, False
            continue
        if "education" in lower or "academic" in lower:
            in_edu, in_exp = True, False
            continue
        if in_exp:
            experience_parts.append(line)
        if in_edu:
            education_parts.append(line)
    return {
        "name": name[:200],
        "email": emails[0] if emails else "",
        "phone": phones[0] if phones else "",
        "skills": skills[:50],
        "experience": " ".join(experience_parts)[:4000],
        "education": " ".join(education_parts)[:2000],
    }


def cleaned_text_for_embedding(raw_text: str, parsed: Dict[str, Any]) -> str:
    """Build normalized, deduplicated text for embedding (lowercase, key fields)."""
    parts = []
    if parsed.get("name"):
        parts.append(parsed["name"].lower())
    if parsed.get("skills"):
        parts.append(" ".join(s.lower() for s in parsed["skills"]))
    if parsed.get("experience"):
        parts.append(extract_text_from_content(parsed["experience"]).lower())
    if parsed.get("education"):
        parts.append(extract_text_from_content(parsed["education"]).lower())
    if not parts:
        parts.append(extract_text_from_content(raw_text).lower())
    combined = " ".join(parts)
    combined = re.sub(r"\s+", " ", combined)
    return combined[:8000]
