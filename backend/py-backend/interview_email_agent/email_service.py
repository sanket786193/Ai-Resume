import os
import smtplib
from email.message import EmailMessage
from dotenv import load_dotenv

from interview_email_agent.email_templates import (
    interview_email_subject,
    interview_email_body
)

load_dotenv()


def send_interview_email(name: str, to_email: str, job_role: str) -> bool:
    try:
        smtp_host = os.getenv("SMTP_HOST")
        smtp_port = int(os.getenv("SMTP_PORT"))
        smtp_email = os.getenv("SMTP_EMAIL")
        smtp_password = os.getenv("SMTP_PASSWORD")

        msg = EmailMessage()
        msg["From"] = smtp_email
        msg["To"] = to_email
        msg["Subject"] = interview_email_subject(job_role)
        msg.set_content(interview_email_body(name, job_role))

        with smtplib.SMTP(smtp_host, smtp_port) as server:
            server.starttls()
            server.login(smtp_email, smtp_password)
            server.send_message(msg)

        return True

    except Exception as e:
        print("Email sending failed:", e)
        return False
