import sys
sys.path.append('..')

from google.adk.agents import SequentialAgent

from resume_parsing.agent import root_agent as resume_agent
from jd_agent.agent import root_agent as jd_agent
from nlp_agent.agent import root_agent as nlp_agent
from ats_agent.agent import root_agent as ats_agent
from skm_agent.agent import root_agent as skm_agent
from .pipeline_runner import root_agent as scoring_agent
from suggestion_agent.agent import root_agent as suggestion_agent
from interview_email_agent.agent import root_agent as interview_email_agent


MainResumePipeline = SequentialAgent(
    name="MainResumePipeline",
    description="End-to-End AI Resume Analysis with Deterministic ATS Scoring",
    sub_agents=[
        resume_agent,
        jd_agent,
        nlp_agent,
        ats_agent,
        skm_agent,
        scoring_agent,      
        suggestion_agent,
        interview_email_agent
    ]
)

# REQUIRED BY ADK
root_agent = MainResumePipeline


# -------------------------
# Local execution (optional)
# -------------------------
if __name__ == "__main__":

    resume_text = """
    Rahul Sharma
    Backend Developer
    Skills: Java, Spring Boot, MySQL
    Experience: 1.5 years backend development
    """

    jd_text = """
    Looking for a Backend Engineer with Java, Spring Boot,
    REST APIs, Microservices, Docker, and SQL.
    """

    result = root_agent.run(
        {
            "resume_text": resume_text,
            "jd_text": jd_text
        }
    )

    print("\nFINAL ATS SCORE:\n", result.get("final_score"))
