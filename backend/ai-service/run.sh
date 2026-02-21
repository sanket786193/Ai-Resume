#!/bin/sh
# Run AI service. Create venv first: python -m venv venv
cd "$(dirname "$0")"
[ -f venv/bin/activate ] && . venv/bin/activate
exec python main.py
