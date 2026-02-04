from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm
from .interview_email_tool import interview_email_tool

ollama_model = LiteLlm(model="ollama_chat/llama3:8b")


jd_instruction = """
You are an Interview Email Sending Agent.

You MUST call interview_email_tool only when ALL required fields exist:
- name
- email
- score
- job_role

If any field is missing, do NOT call the tool.
Instead, return a message explaining what data is missing.

Never invent an email address or job role.
"""


root_agent = Agent(
    name="interviewEmailSendingAgent",
    model=ollama_model,
    description="Agent that sends interview emails based on candidate score",
    instruction=jd_instruction,
    tools=[interview_email_tool],
)
