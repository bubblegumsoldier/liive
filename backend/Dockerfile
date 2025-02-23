# syntax=docker/dockerfile:1
FROM ghcr.io/astral-sh/uv:python3.11-bookworm-slim as builder

# Install build essentials for any packages that need compilation
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        build-essential \
    && rm -rf /var/lib/apt/lists/*

# Enable bytecode compilation and copy mode
ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    PYTHONPATH=/app

# Set up workspace structure
WORKDIR /app
COPY pyproject.toml ./
COPY shared ./shared

ARG APP_NAME
RUN mkdir -p apps
COPY apps/$APP_NAME ./apps/$APP_NAME

# Install shared packages first
RUN --mount=type=cache,target=/root/.cache/uv \
    cd /app/shared/liive-sql-models && \
    uv pip install --system -e .

RUN --mount=type=cache,target=/root/.cache/uv \
    cd /app/shared/liive-common-api && \
    uv pip install --system -e .

# Install service package
RUN --mount=type=cache,target=/root/.cache/uv \
    cd /app/apps/${APP_NAME} && \
    uv pip install --system -e .

# Runtime image
FROM ghcr.io/astral-sh/uv:python3.11-bookworm-slim as runtime

# Create and switch to a non-root user
RUN useradd --create-home appuser
WORKDIR /app
USER appuser

# Set Python path for absolute imports
ENV PYTHONPATH=/app

# Copy installed packages and source code
COPY --from=builder --chown=appuser:appuser /app /app
COPY --from=builder --chown=appuser:appuser /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages

# Set the entrypoint to the app's run.sh
ARG APP_NAME
ENV APP_NAME=${APP_NAME}
RUN ls -la /app/apps/${APP_NAME}
WORKDIR /app/apps/${APP_NAME}
ENTRYPOINT ["./run.sh"] 