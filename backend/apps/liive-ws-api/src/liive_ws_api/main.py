"""WebSocket API service for Liive platform."""
from liive_common_api import create_app

app = create_app(
    title="Liive WebSocket API",
    description="WebSocket service for Liive platform",
    version="0.1.0",
) 