from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings."""

    # API Configuration
    api_title: str = "Liive Auth API"
    api_description: str = "Authentication service for Liive platform"
    api_version: str = "0.1.0"

    # Security
    jwt_secret_key: str = "dev_secret_key"  # Change in production
    jwt_algorithm: str = "HS256"
    access_token_expire_minutes: int = 30

    # CORS
    cors_origins: list[str] = ["http://localhost:3000"]

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=True,
    )


settings = Settings() 