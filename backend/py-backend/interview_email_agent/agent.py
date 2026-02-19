import logging
from google.adk.agents import LlmAgent
from google.adk.models.lite_llm import LiteLlm
from .interview_email_tool import interview_email_tool

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

ollama_model = LiteLlm(model="ollama_chat/llama3:8b")


jd_instruction = """
You are an Interview Email Sending Agent.

CRITICAL: You have ONE tool available. Its EXACT name is: interview_email_tool

When you receive candidate information, extract these required fields:
- name (string) - Candidate's full name
- email (string) - Valid email address
- score (integer) - MUST be a whole number (0, 15, 20, 25, etc.). NOT a string, NOT empty, NOT null
- job_role (string) - Job position title

IMPORTANT TYPE REQUIREMENTS:
- score MUST be an integer number (e.g., 15, 20, 25)
- score MUST NOT be an empty string ""
- score MUST NOT be null or missing
- score MUST NOT be a string like "15" or "twenty"
- If score is not a valid integer, do NOT call the tool

If ALL fields are present AND score is a valid integer, call the function "interview_email_tool" with a dictionary parameter named "candidate".

If ANY field is missing or invalid, do NOT call the tool. Instead, respond: "Missing or invalid required field: [field_name]"

Never invent email addresses or job roles.

CORRECT FUNCTION CALL:
Function name: interview_email_tool
Parameter name: candidate
Parameter value (dictionary):
{
  "name": "John Doe",
  "email": "john@example.com",
  "score": 15,
  "job_role": "Backend Engineer",
  "email_sent": false
}

WRONG EXAMPLES (DO NOT USE):
- "score": ""  (empty string - WRONG)
- "score": "15"  (string - WRONG)
- "score": null  (null - WRONG)
- "score": missing  (missing field - WRONG)

The function name is exactly: interview_email_tool
"""


root_agent = LlmAgent(
    name="InterviewEmailAgent",  # Simplified name to avoid confusion
    model=ollama_model,
    description="Agent that sends interview emails based on candidate score. Available tool: interview_email_tool",
    instruction=jd_instruction,
    tools=[interview_email_tool],
)

# Log registered tools for debugging
logger.info(f"Agent initialized with {len(root_agent.tools)} tool(s)")
for tool in root_agent.tools:
    logger.info(f"Registered tool: {getattr(tool, '__name__', 'Unknown')}")

# Export as 'agent' for compatibility
agent = root_agent

logger.info("Interview Email Agent initialized successfully")
