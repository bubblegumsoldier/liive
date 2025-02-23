"""Common FastAPI application utilities."""

from typing import Any

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.openapi.docs import get_swagger_ui_html
from fastapi.openapi.utils import get_openapi

from liive_common_api.health import health_router


def create_app(
    *,
    title: str,
    description: str,
    version: str = "0.1.0",
    cors_origins: list[str] | None = None,
    tags_metadata: list[dict[str, Any]] | None = None,
    contact: dict[str, str] | None = None,
    license_info: dict[str, str] | None = None,
) -> FastAPI:
    """Create a FastAPI application with common settings.

    Args:
        title: The title of the API
        description: A description of what the API does
        version: The version of the API
        cors_origins: List of allowed CORS origins
        tags_metadata: List of tag metadata for API documentation
        contact: Contact information for the API
        license_info: License information for the API
    """
    # Set default tags if none provided
    if tags_metadata is None:
        tags_metadata = [
            {
                "name": "health",
                "description": "Health check endpoints to verify service status",
            }
        ]

    # Set default contact if none provided
    if contact is None:
        contact = {
            "name": "Liive Team",
            "url": "https://github.com/hmues/liive",
        }

    # Set default license if none provided
    if license_info is None:
        license_info = {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT",
        }

    app = FastAPI(
        title=title,
        description=description,
        version=version,
        openapi_tags=tags_metadata,
        contact=contact,
        license_info=license_info,
        docs_url=None,  # Disable default docs
        redoc_url=None,  # Disable default redoc
    )

    # Add CORS middleware if origins are provided
    if cors_origins:
        app.add_middleware(
            CORSMiddleware,
            allow_origins=cors_origins,
            allow_credentials=True,
            allow_methods=["*"],
            allow_headers=["*"],
        )

    # Include common routers
    app.include_router(health_router, tags=["health"])

    # Custom OpenAPI schema with additional info
    def custom_openapi():
        if app.openapi_schema:
            return app.openapi_schema

        openapi_schema = get_openapi(
            title=app.title,
            version=app.version,
            description=app.description,
            routes=app.routes,
            tags=app.openapi_tags,
            contact=app.contact,
            license_info=app.license_info,
        )

        # Add security schemes if needed
        # openapi_schema["components"]["securitySchemes"] = {...}

        app.openapi_schema = openapi_schema
        return app.openapi_schema

    app.openapi = custom_openapi  # type: ignore

    # Custom docs endpoints
    @app.get("/docs", include_in_schema=False)
    async def custom_swagger_ui_html():
        return get_swagger_ui_html(
            openapi_url=app.openapi_url,  # type: ignore
            title=f"{app.title} - Swagger UI",
            oauth2_redirect_url=app.swagger_ui_oauth2_redirect_url,
            swagger_js_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js",
            swagger_css_url="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css",
        )

    return app
