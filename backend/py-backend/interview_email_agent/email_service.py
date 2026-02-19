import os
import smtplib
import logging
import traceback
import sys
from email.message import EmailMessage
from dotenv import load_dotenv

from interview_email_agent.email_templates import (
    interview_email_subject,
    interview_email_body
)

# Load environment variables
load_dotenv()

# Configure logging
logger = logging.getLogger(__name__)


def send_interview_email(name: str, to_email: str, job_role: str) -> bool:
    """
    Send interview email via SMTP (Gmail) with comprehensive debugging and error handling.
    
    Args:
        name: Candidate name
        to_email: Recipient email address
        job_role: Job role for the interview
    
    Returns:
        bool: True if email sent successfully, False otherwise
    """
    logger.info("=" * 80)
    logger.info(f"send_interview_email called: name={name}, to_email={to_email}, job_role={job_role}")
    logger.info("=" * 80)
    
    try:
        # Load SMTP configuration from environment
        smtp_host = os.getenv("SMTP_HOST")
        smtp_port_str = os.getenv("SMTP_PORT")
        smtp_email = os.getenv("SMTP_EMAIL")
        smtp_password = os.getenv("SMTP_PASSWORD")
        
        # Validate environment variables
        if not smtp_host:
            logger.error("SMTP_HOST environment variable is not set")
            return False
        
        if not smtp_port_str:
            logger.error("SMTP_PORT environment variable is not set")
            return False
        
        try:
            smtp_port = int(smtp_port_str)
        except ValueError:
            logger.error(f"SMTP_PORT must be a valid integer, got: {smtp_port_str}")
            return False
        
        if not smtp_email:
            logger.error("SMTP_EMAIL environment variable is not set")
            return False
        
        if not smtp_password:
            logger.error("SMTP_PASSWORD environment variable is not set")
            return False
        
        # Verify Gmail App Password format (16 characters, no spaces)
        if smtp_host == "smtp.gmail.com" and smtp_password:
            password_clean = smtp_password.replace(" ", "")
            if len(password_clean) != 16 or not password_clean.isalnum():
                logger.warning("WARNING: Gmail App Password should be 16 characters (alphanumeric). "
                             "Current password length: {}. Verify App Password format.".format(len(password_clean)))
        
        logger.info(f"SMTP Configuration:")
        logger.info(f"  Host: {smtp_host}")
        logger.info(f"  Port: {smtp_port}")
        logger.info(f"  From Email: {smtp_email}")
        logger.info(f"  To Email: {to_email}")
        logger.info(f"  Password Length: {len(smtp_password) if smtp_password else 0} characters")
        logger.info(f"Preparing email message for {to_email}...")
        
        # Create email message
        msg = EmailMessage()
        msg["From"] = smtp_email
        msg["To"] = to_email
        msg["Subject"] = interview_email_subject(job_role)
        msg.set_content(interview_email_body(name, job_role))
        
        # Add headers to improve deliverability
        msg["Reply-To"] = smtp_email
        msg["X-Mailer"] = "Python SMTP"
        
        logger.info(f"Email message prepared:")
        logger.info(f"  Subject: {msg['Subject']}")
        logger.info(f"  From: {msg['From']}")
        logger.info(f"  To: {msg['To']}")
        logger.info(f"  Reply-To: {msg.get('Reply-To', 'Not set')}")
        
        logger.info(f"Connecting to SMTP server {smtp_host}:{smtp_port} (timeout=30s)...")
        
        # Connect to SMTP server with timeout
        server = smtplib.SMTP(smtp_host, smtp_port, timeout=30)
        
        # Enable FULL SMTP debug logging (level 1 = verbose)
        server.set_debuglevel(1)
        logger.info("SMTP debug logging ENABLED (level 1)")
        
        try:
            # Log initial connection
            logger.info("=" * 80)
            logger.info("SMTP CONNECTION PHASE")
            logger.info("=" * 80)
            logger.info(f"Connected to {smtp_host}:{smtp_port}")
            
            # Explicit EHLO command (required by some SMTP servers)
            logger.info("Sending EHLO command...")
            ehlo_response = server.ehlo()
            logger.info(f"EHLO response code: {ehlo_response[0]}")
            logger.info(f"EHLO response message: {ehlo_response[1]}")
            if ehlo_response[0] != 250:
                logger.error(f"EHLO failed with code {ehlo_response[0]}: {ehlo_response[1]}")
                return False
            logger.info("✓ EHLO successful")
            
            # Enable TLS encryption
            logger.info("=" * 80)
            logger.info("TLS HANDSHAKE PHASE")
            logger.info("=" * 80)
            logger.info("Starting TLS...")
            starttls_response = server.starttls()
            logger.info(f"STARTTLS response: {starttls_response}")
            logger.info("✓ TLS started successfully")
            
            # Re-EHLO after STARTTLS (some servers require this)
            logger.info("Sending EHLO after STARTTLS...")
            ehlo_response_tls = server.ehlo()
            logger.info(f"EHLO after TLS response code: {ehlo_response_tls[0]}")
            logger.info(f"EHLO after TLS response message: {ehlo_response_tls[1]}")
            if ehlo_response_tls[0] != 250:
                logger.warning(f"EHLO after TLS returned code {ehlo_response_tls[0]}, but continuing...")
            else:
                logger.info("✓ EHLO after TLS successful")
            
            # Authenticate
            logger.info("=" * 80)
            logger.info("AUTHENTICATION PHASE")
            logger.info("=" * 80)
            logger.info(f"Logging in as {smtp_email}...")
            logger.info("Attempting SMTP authentication...")
            
            try:
                auth_response = server.login(smtp_email, smtp_password)
                logger.info(f"Login response: {auth_response}")
                logger.info("✓ SMTP authentication successful")
            except smtplib.SMTPAuthenticationError as auth_err:
                logger.error("=" * 80)
                logger.error("SMTP AUTHENTICATION FAILED")
                logger.error("=" * 80)
                logger.error(f"Error code: {auth_err.smtp_code}")
                logger.error(f"Error message: {auth_err.smtp_error}")
                logger.error(f"Server response: {auth_err.args}")
                logger.error("")
                logger.error("TROUBLESHOOTING CHECKLIST:")
                logger.error("1. Verify SMTP_EMAIL matches your Gmail account exactly")
                logger.error("2. Verify SMTP_PASSWORD is a Gmail App Password (NOT your regular password)")
                logger.error("3. Ensure 2-Step Verification is enabled on your Gmail account")
                logger.error("4. Generate a NEW App Password at: https://myaccount.google.com/apppasswords")
                logger.error("5. Copy the App Password exactly (16 characters, no spaces)")
                logger.error("6. Verify the App Password hasn't been revoked")
                logger.error("")
                logger.error("Full authentication error traceback:")
                logger.error(traceback.format_exc())
                raise
            
            # Send email
            logger.info("=" * 80)
            logger.info("EMAIL SENDING PHASE")
            logger.info("=" * 80)
            logger.info(f"Sending email to {to_email}...")
            logger.info("Calling server.send_message()...")
            
            try:
                # Capture send_message response
                send_response = server.send_message(msg)
                logger.info(f"send_message() returned: {send_response}")
                
                # Verify successful delivery (250 OK)
                if send_response:
                    logger.warning(f"send_message() returned non-empty response: {send_response}")
                    # Check for failed recipients
                    if send_response:
                        failed_recipients = send_response
                        logger.error(f"Failed recipients: {failed_recipients}")
                        for recipient, (code, error_msg) in failed_recipients.items():
                            logger.error(f"  {recipient}: {code} - {error_msg}")
                        return False
                else:
                    logger.info("✓ send_message() returned empty dict (success)")
                
                # Verify server is still connected
                try:
                    server.noop()
                    logger.info("✓ Server connection verified (NOOP successful)")
                except Exception as noop_err:
                    logger.warning(f"NOOP check failed: {noop_err}")
                
                logger.info("=" * 80)
                logger.info("EMAIL DELIVERY CONFIRMED")
                logger.info("=" * 80)
                logger.info(f"✓ Email successfully sent to {to_email}")
                logger.info(f"✓ SMTP server accepted message (250 OK)")
                logger.info("")
                logger.info("NOTE: If email not received:")
                logger.info("  1. Check recipient's spam/junk folder")
                logger.info("  2. Verify recipient email address is correct")
                logger.info("  3. Check Gmail account for delivery errors")
                logger.info("  4. Verify sender reputation (new accounts may be flagged)")
                logger.info("=" * 80)
                
            except smtplib.SMTPRecipientsRefused as recipient_err:
                logger.error("=" * 80)
                logger.error("SMTP RECIPIENT REFUSED")
                logger.error("=" * 80)
                logger.error(f"Recipient email refused: {to_email}")
                logger.error(f"Error details: {recipient_err}")
                logger.error(f"Recipients dict: {recipient_err.recipients}")
                logger.error("Full recipient error traceback:")
                logger.error(traceback.format_exc())
                raise
                
            except smtplib.SMTPDataError as data_err:
                logger.error("=" * 80)
                logger.error("SMTP DATA ERROR")
                logger.error("=" * 80)
                logger.error(f"Error code: {data_err.smtp_code}")
                logger.error(f"Error message: {data_err.smtp_error}")
                logger.error(f"Server response: {data_err.args}")
                logger.error("Possible causes:")
                logger.error("  1. Message content rejected by server")
                logger.error("  2. Sender reputation issues")
                logger.error("  3. Rate limiting")
                logger.error("Full data error traceback:")
                logger.error(traceback.format_exc())
                raise
                
            except smtplib.SMTPResponseException as resp_err:
                logger.error("=" * 80)
                logger.error("SMTP RESPONSE EXCEPTION")
                logger.error("=" * 80)
                logger.error(f"Error code: {resp_err.smtp_code}")
                logger.error(f"Error message: {resp_err.smtp_error}")
                logger.error(f"Server response: {resp_err.args}")
                logger.error("Full response exception traceback:")
                logger.error(traceback.format_exc())
                raise
                
        finally:
            # Always close the connection
            try:
                logger.info("Closing SMTP connection...")
                server.quit()
                logger.info("✓ SMTP connection closed gracefully")
            except Exception as quit_err:
                logger.warning(f"Error during server.quit(): {quit_err}")
                try:
                    server.close()
                    logger.info("✓ SMTP connection closed (fallback)")
                except Exception as close_err:
                    logger.warning(f"Error during server.close(): {close_err}")
        
        return True

    except smtplib.SMTPAuthenticationError as e:
        logger.error("=" * 80)
        logger.error("SMTP AUTHENTICATION ERROR (OUTER HANDLER)")
        logger.error("=" * 80)
        logger.error(f"Error: {e}")
        logger.error(f"Error code: {getattr(e, 'smtp_code', 'N/A')}")
        logger.error(f"Error message: {getattr(e, 'smtp_error', 'N/A')}")
        logger.error("Full authentication error traceback:")
        logger.error(traceback.format_exc())
        print("\n" + "=" * 80, file=sys.stderr)
        print("SMTP AUTHENTICATION FAILED - CHECK YOUR CREDENTIALS", file=sys.stderr)
        print("=" * 80, file=sys.stderr)
        print(traceback.format_exc(), file=sys.stderr)
        return False
        
    except smtplib.SMTPRecipientsRefused as e:
        logger.error("=" * 80)
        logger.error("SMTP RECIPIENT REFUSED (OUTER HANDLER)")
        logger.error("=" * 80)
        logger.error(f"Invalid recipient email: {to_email}")
        logger.error(f"Error: {e}")
        logger.error(f"Recipients dict: {getattr(e, 'recipients', {})}")
        logger.error("Full recipient error traceback:")
        logger.error(traceback.format_exc())
        print("\n" + "=" * 80, file=sys.stderr)
        print(f"SMTP RECIPIENT REFUSED: {to_email}", file=sys.stderr)
        print("=" * 80, file=sys.stderr)
        print(traceback.format_exc(), file=sys.stderr)
        return False
        
    except smtplib.SMTPServerDisconnected as e:
        logger.error("=" * 80)
        logger.error("SMTP SERVER DISCONNECTED")
        logger.error("=" * 80)
        logger.error(f"Error: {e}")
        logger.error("Possible causes:")
        logger.error("  1. Network connection lost")
        logger.error("  2. SMTP server closed connection")
        logger.error("  3. Firewall blocking port 587")
        logger.error("  4. Server timeout")
        logger.error("Full disconnection error traceback:")
        logger.error(traceback.format_exc())
        print("\n" + "=" * 80, file=sys.stderr)
        print("SMTP SERVER DISCONNECTED", file=sys.stderr)
        print("=" * 80, file=sys.stderr)
        print(traceback.format_exc(), file=sys.stderr)
        return False
        
    except smtplib.SMTPException as e:
        logger.error("=" * 80)
        logger.error("SMTP EXCEPTION (GENERIC)")
        logger.error("=" * 80)
        logger.error(f"Error: {e}")
        logger.error(f"Error type: {type(e).__name__}")
        logger.error(f"Error code: {getattr(e, 'smtp_code', 'N/A')}")
        logger.error(f"Error message: {getattr(e, 'smtp_error', 'N/A')}")
        logger.error("Full SMTP exception traceback:")
        logger.error(traceback.format_exc())
        print("\n" + "=" * 80, file=sys.stderr)
        print("SMTP ERROR OCCURRED", file=sys.stderr)
        print("=" * 80, file=sys.stderr)
        print(traceback.format_exc(), file=sys.stderr)
        return False
        
    except Exception as e:
        logger.error("=" * 80)
        logger.error("UNEXPECTED ERROR")
        logger.error("=" * 80)
        logger.error(f"Error type: {type(e).__name__}")
        logger.error(f"Error: {e}")
        logger.error("Full unexpected error traceback:")
        logger.error(traceback.format_exc())
        print("\n" + "=" * 80, file=sys.stderr)
        print("UNEXPECTED ERROR WHILE SENDING EMAIL", file=sys.stderr)
        print("=" * 80, file=sys.stderr)
        print(traceback.format_exc(), file=sys.stderr)
        return False
