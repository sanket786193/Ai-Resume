import sys
sys.path.append('backend/py-backend')

try:
    from pipeline_agent.agent import root_agent
    print('Import successful')
except Exception as e:
    print(f'Import failed: {e}')
