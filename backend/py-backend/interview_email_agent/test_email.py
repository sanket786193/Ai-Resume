import logging
from dotenv import load_dotenv
from interview_email_agent.email_controller import process_interview_email
from interview_email_agent.interview_email_tool import interview_email_tool

# Load environment variables
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

print("=" * 60)
print("TEST 1: Direct email sending (bypass LLM)")
print("=" * 60)

# Test with score < 20 (should send email)
candidate_data = {
    "name": "Test User",
    "email": "mihirkhode90@gmail.com",  # Recipient email
    "score": 15,  # Below cutoff (should send email)
    "job_role": "Backend Engineer",
    "email_sent": False
}

print(f"\nTesting with candidate: {candidate_data}")
result = process_interview_email(candidate_data)
print(f"Result: {result}")

print("\n" + "=" * 60)
print("TEST 2: Tool function call (simulating LLM call)")
print("=" * 60)

# Test tool function directly
candidate_data2 = {
    "name": "Test User 2",
    "email": "mihirkhode90@gmail.com",
    "score": 15,
    "job_role": "Backend Engineer",
    "email_sent": False
}

print(f"\nTesting tool with candidate: {candidate_data2}")
tool_result = interview_email_tool(candidate_data2)
print(f"Tool result: {tool_result}")

print("\n" + "=" * 60)
print("TEST 3: Score >= cutoff (should block)")
print("=" * 60)

candidate_data3 = {
    "name": "Test User 3",
    "email": "mihirkhode90@gmail.com",
    "score": 25,  # Above cutoff (should block)
    "job_role": "Backend Engineer",
    "email_sent": False
}

print(f"\nTesting with candidate: {candidate_data3}")
tool_result3 = interview_email_tool(candidate_data3)
print(f"Tool result: {tool_result3}")
