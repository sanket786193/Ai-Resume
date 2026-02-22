from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm

from app.config import get_settings

_settings = get_settings()
ollama_model = LiteLlm(model=f"ollama_chat/{_settings.ollama_model}")

resume_instruction = """
You are a Resume Analysis AI. Analyze the resume below and output structured JSON only.

Resume text:
{resume_text}

Output MUST be valid JSON only (no markdown, no explanation). Use this schema:
{
  "name": "",
  "role": "",
  "summary": "",
  "skills": ["skill1", "skill2"],
  "experience": [{"company": "", "duration": "", "responsibilities": []}],
  "projects": [],
  "education": [],
  "certifications": []
}
"""

root_agent = Agent(
    name="ResumeParsingAgent",
    model=ollama_model,
    description="Parse resume into structured JSON for downstream agents",
    instruction=resume_instruction,
    tools=[],
    output_key="parsed_resume",
)
