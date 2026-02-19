"""
Debug script to verify tool registration and test direct tool calls.
"""
import logging
from dotenv import load_dotenv
from interview_email_agent.agent import agent

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

print("=" * 60)
print("DEBUG: Checking registered tools")
print("=" * 60)

# Check what tools are registered
if hasattr(agent, 'tools'):
    print(f"\nAgent tools: {agent.tools}")
    print(f"Number of tools: {len(agent.tools) if agent.tools else 0}")
    
    for i, tool in enumerate(agent.tools):
        print(f"\nTool {i+1}:")
        print(f"  Name: {getattr(tool, '__name__', 'Unknown')}")
        print(f"  Type: {type(tool)}")
        if hasattr(tool, '__doc__'):
            print(f"  Docstring: {tool.__doc__[:100] if tool.__doc__ else 'None'}...")
else:
    print("Agent has no 'tools' attribute")

print("\n" + "=" * 60)
print("Agent name:", getattr(agent, 'name', 'Unknown'))
print("=" * 60)
