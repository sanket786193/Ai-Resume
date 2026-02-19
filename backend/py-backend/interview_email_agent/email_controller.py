import os
import logging
from typing import Union, Optional
from dotenv import load_dotenv
from interview_email_agent.email_service import send_interview_email

# Load environment variables
load_dotenv()

# Configure logging
logger = logging.getLogger(__name__)

HR_CUTOFF_SCORE = int(os.getenv("HR_CUTOFF_SCORE", "20"))
AUTO_EMAIL_ENABLED = os.getenv("AUTO_EMAIL_ENABLED", "true").lower() == "true"

logger.info(f"Email Controller initialized - HR_CUTOFF_SCORE: {HR_CUTOFF_SCORE}, AUTO_EMAIL_ENABLED: {AUTO_EMAIL_ENABLED}")


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


def process_interview_email(candidate: dict) -> str:
    """
    Process and send interview email for a candidate.
    
    Args:
        candidate: Dictionary containing:
            - name: str (required)
            - email: str (required)
            - score: int (required)
            - job_role: str (required)
            - email_sent: bool (optional)
    
    Returns:
        str: Status message indicating success or failure reason
    """
    logger.info(f"process_interview_email called with candidate: {candidate}")
    
    candidate_name = candidate.get("name", "")
    candidate_email = candidate.get("email", "")
    raw_score = candidate.get("score")
    candidate_job_role = candidate.get("job_role", "")
    
    # Safely parse and validate score
    candidate_score, score_error = safe_parse_score(raw_score)
    
    logger.info(f"Raw score value: {raw_score!r} (type: {type(raw_score).__name__})")
    logger.info(f"Parsed score: {candidate_score}")
    if score_error:
        logger.warning(f"Score parsing issue: {score_error}")
    
    logger.info(f"Processing email for: {candidate_name} ({candidate_email}), Score: {candidate_score}, Role: {candidate_job_role}")

    if not AUTO_EMAIL_ENABLED:
        reason = "Auto email disabled by HR"
        logger.warning(reason)
        return reason

    # Validate score before comparison
    if score_error:
        reason = f"Invalid candidate score: {score_error}. Cannot process email."
        logger.error(reason)
        return reason
    
    # Type safety check before comparison
    if not isinstance(candidate_score, int):
        error_msg = f"Type safety violation: candidate_score is {type(candidate_score).__name__}, expected int"
        logger.error(error_msg)
        return f"Internal error: {error_msg}"

    # Note: This check is redundant since interview_email_tool already checks this,
    # but keeping it as a safety measure
    if candidate_score >= HR_CUTOFF_SCORE:
        reason = f"Candidate selected by HR score ({candidate_score} >= {HR_CUTOFF_SCORE}); skip auto email"
        logger.info(reason)
        return reason

    if candidate.get("email_sent"):
        reason = "Email already sent"
        logger.info(reason)
        return reason

    # Validate required fields
    if not candidate_name:
        reason = "Missing candidate name"
        logger.error(reason)
        return reason
    
    if not candidate_email:
        reason = "Missing candidate email"
        logger.error(reason)
        return reason
    
    if not candidate_job_role:
        reason = "Missing job role"
        logger.error(reason)
        return reason

    logger.info(f"Attempting to send interview email to {candidate_email}...")
    
    try:
        success = send_interview_email(
            name=candidate_name,
            to_email=candidate_email,
            job_role=candidate_job_role
        )

        if success:
            candidate["email_sent"] = True
            success_msg = f"Interview email sent successfully to {candidate_email}"
            logger.info(success_msg)
            return success_msg
        else:
            error_msg = f"Failed to send interview email to {candidate_email}"
            logger.error(error_msg)
            return error_msg
            
    except Exception as e:
        error_msg = f"Exception occurred while sending email: {str(e)}"
        logger.exception(error_msg)
        return error_msg
