#!/bin/bash
set -e

# Run the service with uvicorn
exec python -m uvicorn src.liive_auth_api.main:app --host ${HOST:-0.0.0.0} --port ${PORT:-8000} ${UVICORN_EXTRA_ARGS:-} 