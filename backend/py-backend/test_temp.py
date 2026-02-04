from interview_email_agent.email_controller import process_interview_email

candidate_data = {
    "name": "Test User",
    "email": "mihirkhode90@gmail.com",  # use your own email
    "score": 20,                      # above cutoff
    "job_role": "Backend Engineer",
    "email_sent": False
}

result = process_interview_email(candidate_data)
print(result)
