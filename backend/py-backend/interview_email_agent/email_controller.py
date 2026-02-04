import os
from interview_email_agent.email_service import send_interview_email


# HR_CUTOFF_SCORE = int(os.getenv("HR_CUTOFF_SCORE", 20))
AUTO_EMAIL_ENABLED = os.getenv("AUTO_EMAIL_ENABLED", "true").lower() == "true"


def process_interview_email(candidate: dict) -> str:
    """
    candidate = {
        "name": "Rahul",
        "email": "mihirkhode90@gmail.com",
        "score": 20,
        "job_role": "Backend Engineer",
        "email_sent": False
    }
    """

    if not AUTO_EMAIL_ENABLED:
        return "Auto email disabled by HR"

    # if candidate["score"] < HR_CUTOFF_SCORE:
    #     return "Candidate score below cutoff"

    if candidate.get("email_sent"):
        return "Email already sent"

    success = send_interview_email(
        name=candidate["name"],
        to_email=candidate["email"],
        job_role=candidate["job_role"]
    )

    if success:
        candidate["email_sent"] = True
        return "Interview email sent successfully"

    return "Failed to send interview email"
