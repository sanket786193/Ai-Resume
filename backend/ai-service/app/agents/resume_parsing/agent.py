from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm

from app.config import get_settings

_settings = get_settings()
ollama_model = LiteLlm(model=f"ollama_chat/{_settings.ollama_model}")

resume_instruction = """
You are a Resume Analysis AI.
Analyze resumes using ATS and hiring standards.
Provide structured feedback, skill extraction, and improvement suggestions.
"""

root_agent = Agent(
    name="ResumeParsingAgent",
    model=ollama_model,
    description="Resume parsing and analysis agent",
    instruction=resume_instruction,
    tools=[],
)
