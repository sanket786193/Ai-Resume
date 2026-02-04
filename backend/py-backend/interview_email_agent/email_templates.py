def interview_email_subject(job_role: str) -> str:
    return f"Interview Invitation – {job_role}"


def interview_email_body(name: str, job_role: str) -> str:
    return f"""
Dear {name},

Congratulations! 🎉

You have been shortlisted for the {job_role} position based on your profile evaluation.

Our HR team will contact you shortly with interview details.

Best Regards,
HR Team
"""
