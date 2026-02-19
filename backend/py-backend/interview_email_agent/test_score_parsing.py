"""
Test script to verify safe score parsing handles all edge cases.
"""
import logging
from dotenv import load_dotenv
from interview_email_agent.interview_email_tool import interview_email_tool, safe_parse_score

load_dotenv()

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def test_safe_parse_score():
    """Test safe_parse_score function with various inputs."""
    print("\n" + "=" * 60)
    print("TESTING safe_parse_score() FUNCTION")
    print("=" * 60)
    
    test_cases = [
        # (input, expected_score, should_have_error)
        (15, 15, False),           # Valid int
        ("15", 15, False),          # Valid string int
        ("20", 20, False),          # Valid string int
        (20.5, 20, False),          # Float (should round down)
        ("20.7", 20, False),        # String float
        ("", 0, True),              # Empty string - ERROR
        (None, 0, True),             # None - ERROR
        ("abc", 0, True),           # Non-numeric string - ERROR
        ("  ", 0, True),            # Whitespace only - ERROR
        (0, 0, False),              # Zero (valid)
        (-5, -5, False),            # Negative (valid, though unusual)
    ]
    
    print("\nTest Case | Input | Expected Score | Has Error | Result")
    print("-" * 70)
    
    all_passed = True
    for i, (input_val, expected_score, should_have_error) in enumerate(test_cases, 1):
        score, error = safe_parse_score(input_val)
        has_error = error is not None
        
        # Check if result matches expectation
        score_match = score == expected_score
        error_match = has_error == should_have_error
        
        passed = score_match and error_match
        
        if not passed:
            all_passed = False
        
        status = "✓ PASS" if passed else "✗ FAIL"
        print(f"{i:2d} | {repr(input_val):15s} | {expected_score:14d} | {str(should_have_error):9s} | {status}")
        if error:
            print(f"    Error message: {error}")
    
    print("\n" + "=" * 60)
    if all_passed:
        print("ALL TESTS PASSED ✓")
    else:
        print("SOME TESTS FAILED ✗")
    print("=" * 60)


def test_interview_email_tool_with_edge_cases():
    """Test interview_email_tool with various score edge cases."""
    print("\n" + "=" * 60)
    print("TESTING interview_email_tool() WITH EDGE CASES")
    print("=" * 60)
    
    base_candidate = {
        "name": "Test User",
        "email": "test@example.com",
        "job_role": "Backend Engineer",
        "email_sent": False
    }
    
    test_cases = [
        # (score_value, description)
        (15, "Valid integer - should process"),
        ("15", "Valid string integer - should process"),
        ("", "Empty string - should return error"),
        (None, "None - should return error"),
        ("abc", "Non-numeric string - should return error"),
        (20, "Valid integer at cutoff - should block"),
        (25, "Valid integer above cutoff - should block"),
    ]
    
    for score_value, description in test_cases:
        print(f"\n--- Test: {description} ---")
        print(f"Score value: {repr(score_value)} (type: {type(score_value).__name__})")
        
        candidate = base_candidate.copy()
        candidate["score"] = score_value
        
        try:
            result = interview_email_tool(candidate)
            print(f"Result: {result}")
            print(f"Email sent: {result.get('email_sent', False)}")
            if 'error' in result:
                print(f"Error type: {result.get('error')}")
        except Exception as e:
            print(f"EXCEPTION (should not happen): {type(e).__name__}: {e}")
            import traceback
            traceback.print_exc()


if __name__ == "__main__":
    print("\n" + "=" * 60)
    print("SCORE PARSING TEST SUITE")
    print("=" * 60)
    
    # Test the parsing function
    test_safe_parse_score()
    
    # Test the tool function
    test_interview_email_tool_with_edge_cases()
    
    print("\n" + "=" * 60)
    print("TEST SUITE COMPLETE")
    print("=" * 60)
