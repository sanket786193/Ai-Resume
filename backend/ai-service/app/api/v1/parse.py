"""Resume parse endpoint: raw text, structured fields, cleaned text (best practice: keep separately)."""
from urllib.parse import quote, unquote, urlparse, urlunparse
from fastapi import APIRouter, HTTPException
import httpx
from pypdf import PdfReader
from io import BytesIO

from app.core.logging import get_logger
from app.services.resume_parser import (
    cleaned_text_for_embedding,
    extract_text_from_content,
    parse_resume_structured,
)

from .schemas import ParseRequest, ParseResponse

logger = get_logger(__name__)

router = APIRouter()


def _normalize_resume_url(url: str) -> str:
    """
    Normalize URL path so path segments are properly percent-encoded.
    Cloudinary (and some CDNs) can return 400 for paths with unencoded
    spaces or parentheses; encoding the path often fixes the request.
    """
    try:
        parsed = urlparse(url)
        if not parsed.path:
            return url
        # Decode then re-encode each path segment (avoids double-encoding).
        segments = [quote(unquote(s), safe="") for s in parsed.path.split("/")]
        new_path = "/".join(segments)
        return urlunparse((
            parsed.scheme,
            parsed.netloc,
            new_path,
            parsed.params,
            parsed.query,
            parsed.fragment,
        ))
    except Exception:
        return url


def _fetch_resume_text(url_or_content: str) -> str:
    """If input looks like a URL, fetch and extract PDF text; else return as-is."""
    s = (url_or_content or "").strip()
    if not s.startswith(("http://", "https://")):
        return s
    fetch_url = _normalize_resume_url(s)
    try:
        with httpx.Client(timeout=30) as client:
            r = client.get(fetch_url)
            r.raise_for_status()
            content_type = (r.headers.get("content-type") or "").lower()
            if "pdf" in content_type or s.lower().endswith(".pdf"):
                reader = PdfReader(BytesIO(r.content))
                parts = []
                for page in reader.pages:
                    parts.append(page.extract_text() or "")
                return "\n".join(parts)
            return r.text[:100_000]
    except Exception as e:
        logger.warning("Fetch resume from URL failed: %s", e)
        return ""


@router.post("", response_model=ParseResponse, summary="Parse resume")
@router.post("/", response_model=ParseResponse)
def parse_resume(req: ParseRequest) -> ParseResponse:
    """
    Extract raw text (from URL or inline), parse structured fields (name, email, phone, skills, experience, education),
    and build cleaned text for embedding.
    """
    raw = (req.resume_path_or_content or "").strip()
    if raw.startswith(("http://", "https://")):
        raw = _fetch_resume_text(raw)
    if not raw:
        return ParseResponse(
            raw_text="",
            parsed={
                "name": "",
                "email": "",
                "phone": "",
                "skills": [],
                "experience": "",
                "education": "",
            },
            cleaned_text="",
        )
    raw_text = extract_text_from_content(raw)
    parsed = parse_resume_structured(raw_text or raw)
    cleaned_text = cleaned_text_for_embedding(raw_text or raw, parsed)
    return ParseResponse(raw_text=raw_text or raw, parsed=parsed, cleaned_text=cleaned_text)
