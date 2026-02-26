"""Screening pipeline: resume -> jd -> nlp -> ats -> skm -> scoring (no email/suggestion)."""
from google.adk.agents import SequentialAgent

from app.agents.ats_agent import root_agent as ats_agent
from app.agents.jd_agent import root_agent as jd_agent
from app.agents.nlp_agent import root_agent as nlp_agent
from app.agents.pipeline_agent.pipeline_runner import ScoringAgent
from app.agents.resume_parsing import root_agent as resume_agent
from app.agents.skm_agent import root_agent as skm_agent

# ScoringAgent is a singleton instance (rule-based, no LLM)
scoring_agent = ScoringAgent

screening_pipeline = SequentialAgent(
    name="ScreeningPipeline",
    description="Resume screening for ATS: parse, JD, NLP, ATS, SKM, rule-based scoring",
    sub_agents=[
        resume_agent,
        jd_agent,
        nlp_agent,
        ats_agent,
        skm_agent,
        scoring_agent,
    ],
)
