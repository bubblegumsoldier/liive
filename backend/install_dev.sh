#!/bin/bash
set -e

# Create and activate virtual environment if it doesn't exist
if [ ! -d ".venv" ]; then
    python -m venv .venv
fi
source .venv/bin/activate

# Install uv if not already installed
if ! command -v uv &> /dev/null; then
    curl -LsSf https://astral.sh/uv/install.sh | sh
fi

# Install build dependencies
uv pip install hatch

# Install shared packages first
cd shared/liive-sql-models
uv pip install -e .
cd ../..

# Install app packages
for app in apps/*; do
    if [ -d "$app" ]; then
        cd "$app"
        uv pip install -e .
        cd ../..
    fi
done

# Install development dependencies
uv pip install -e ".[dev]" 