import asyncio
import logging
from dotenv import load_dotenv
from interview_email_agent.agent import agent

# Load environment variables BEFORE importing agent
load_dotenv()

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# IMPORTANT: Use score < HR_CUTOFF_SCORE (default 20) to test email sending
# Score >= 20 will be blocked by HR cutoff logic
prompt = """
Send an interview email to the following candidate.

Candidate information:
- name: Test User
- email: mihirkhode90@gmail.com
- score: 15
- job_role: Backend Engineer
- email_sent: false

Use the interview_email_tool function to send the email.
"""

async def main():
    logger.info("Starting Interview Email Agent...")
    logger.info(f"Prompt: {prompt}")
    
    try:
        async for event in agent.run_async(prompt):
            logger.info(f"Agent event: {event}")
            print(event)
    except Exception as e:
        logger.error(f"Error running agent: {e}")
        logger.exception("Full exception traceback:")
        raise

if __name__ == "__main__":
    asyncio.run(main())

