import os
import logging
from typing import Dict, Optional, Union
from dotenv import load_dotenv
from .email_controller import process_interview_email

# Load environment variables
load_dotenv()

# Configure logging
logger = logging.getLogger(__name__)

HR_CUTOFF_SCORE = int(os.getenv("HR_CUTOFF_SCORE", "20"))


def safe_parse_score(score_value: Union[str, int, float, None]) -> tuple[int, Optional[str]]:
    """
    Safely parse score value to integer with comprehensive error handling.
    
    Args:
        score_value: Score value that could be int, str, None, or empty string
    
    Returns:
        tuple: (parsed_score: int, error_message: Optional[str])
        - If parsing succeeds: (score_int, None)
        - If parsing fails: (0, error_message)
    """
    # Handle None
    if score_value is None:
        logger.warning("Score is None, defaulting to 0")
        return 0, "Score is missing or None"
    
    # Handle empty string
    if isinstance(score_value, str):
        score_value = score_value.strip()
        if not score_value:
            logger.warning("Score is empty string, defaulting to 0")
            return 0, "Score is empty or not provided"
    
    # If already an integer, return it
    if isinstance(score_value, int):
        return score_value, None
    
    # If float, convert to int (round down)
    if isinstance(score_value, float):
        logger.info(f"Score is float ({score_value}), converting to int: {int(score_value)}")
        return int(score_value), None
    
    # Try to parse string to integer
    if isinstance(score_value, str):
        try:
            # Try parsing as float first (handles "15.5" -> 15)
            parsed = float(score_value)
            parsed_int = int(parsed)
            logger.info(f"Successfully parsed score string '{score_value}' to int: {parsed_int}")
            return parsed_int, None
        except ValueError:
            error_msg = f"Score '{score_value}' is not a valid number"
            logger.error(error_msg)
            return 0, error_msg
    
    # Unknown type
    error_msg = f"Score has unsupported type: {type(score_value).__name__}"
    logger.error(error_msg)
    return 0, error_msg


def interview_email_tool(candidate: dict) -> dict:
    """
    Send interview email to a candidate.
    
    This tool sends an interview invitation email to candidates who meet the criteria.
    Only candidates with score below HR_CUTOFF_SCORE will receive emails.
    
    Args:
        candidate: Dictionary containing candidate information
            - name: str (required) - Candidate's full name
            - email: str (required) - Candidate's email address
            - score: int (required) - Candidate's evaluation score
            - job_role: str (required) - Job role/position name
            - email_sent: bool (optional, defaults to False) - Whether email was already sent
    
    Returns:
        dict with email_sent status (bool) and result/reason (str)
        
    Example:
        interview_email_tool({
            "name": "John Doe",
            "email": "john@example.com",
            "score": 15,
            "job_role": "Backend Engineer",
            "email_sent": False
        })
    """
    logger.info(f"interview_email_tool called with candidate: {candidate}")
    
    # Extract candidate data
    raw_score = candidate.get("score")
    candidate_email = candidate.get("email", "")
    candidate_name = candidate.get("name", "")
    candidate_job_role = candidate.get("job_role", "")
    
    # Safely parse and validate score
    candidate_score, score_error = safe_parse_score(raw_score)
    
    # Log score parsing result
    logger.info(f"Raw score value: {raw_score!r} (type: {type(raw_score).__name__})")
    logger.info(f"Parsed score: {candidate_score}")
    if score_error:
        logger.warning(f"Score parsing issue: {score_error}")
    
    logger.info(f"Processing candidate: {candidate_name}, email: {candidate_email}, score: {candidate_score}, role: {candidate_job_role}")
    logger.info(f"HR_CUTOFF_SCORE: {HR_CUTOFF_SCORE} (type: {type(HR_CUTOFF_SCORE).__name__})")
    
    # Validate score before comparison
    if score_error:
        reason = f"Invalid candidate score: {score_error}. Cannot determine email eligibility."
        logger.warning(reason)
        return {
            "email_sent": False,
            "reason": reason,
            "error": "invalid_score",
            "raw_score": str(raw_score) if raw_score is not None else "None"
        }
    
    # Type safety check before comparison
    if not isinstance(candidate_score, int):
        error_msg = f"Type safety violation: candidate_score is {type(candidate_score).__name__}, expected int"
        logger.error(error_msg)
        return {
            "email_sent": False,
            "reason": f"Internal error: {error_msg}",
            "error": "type_mismatch"
        }
    
    # Do not send email to candidates selected by HR score (they are handled by HR)
    if candidate_score >= HR_CUTOFF_SCORE:
        reason = f"Candidate selected by HR score ({candidate_score} >= {HR_CUTOFF_SCORE}); HR will contact directly"
        logger.info(reason)
        return {
            "email_sent": False,
            "reason": reason
        }

    if not candidate_email:
        reason = "Missing candidate email"
        logger.warning(reason)
        return {
            "email_sent": False,
            "reason": reason
        }

    if not candidate_job_role:
        reason = "Missing job role"
        logger.warning(reason)
        return {
            "email_sent": False,
            "reason": reason
        }
    
    if not candidate_name:
        reason = "Missing candidate name"
        logger.warning(reason)
        return {
            "email_sent": False,
            "reason": reason
        }

    logger.info(f"All validations passed. Processing interview email for {candidate_name}")
    
    # Update candidate dict with parsed score to ensure type safety downstream
    candidate["score"] = candidate_score
    
    result = process_interview_email(candidate)

    # Check if email was actually sent by examining the result
    # Success message format: "Interview email sent successfully to {email}"
    # Failure messages: "Failed to send...", "Exception occurred...", "Missing...", etc.
    result_lower = result.lower()
    
    # The success message contains "sent successfully" which is unique
    # Failure messages contain "failed", "exception", "error", "missing", etc.
    # Check for the exact success pattern
    email_sent = "sent successfully" in result_lower
    
    logger.info(f"Email processing result: {result}")
    logger.info(f"Email sent status: {email_sent}")
    
    if not email_sent:
        logger.warning(f"Email was not sent. Result: {result}")
    
    return {
        "email_sent": email_sent,
        "result": result,
        "candidate_name": candidate_name,
        "candidate_email": candidate_email
    }
