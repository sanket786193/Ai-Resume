import os
from typing import Dict
from .email_controller import process_interview_email

HR_CUTOFF_SCORE = int(os.getenv("HR_CUTOFF_SCORE", 20))
def interview_email_tool(candidate: dict) -> dict:
    if candidate.get("score", 0) < HR_CUTOFF_SCORE:
        return {
            "email_sent": False,
            "reason": "Score below HR cutoff"
        }

    if not candidate.get("email"):
        return {
            "email_sent": False,
            "reason": "Missing candidate email"
        }

    if not candidate.get("job_role"):
        return {
            "email_sent": False,
            "reason": "Missing job role"
        }

    result = process_interview_email(candidate)

    return {
        "email_sent": candidate.get("email_sent", False),
        "result": result
    }
