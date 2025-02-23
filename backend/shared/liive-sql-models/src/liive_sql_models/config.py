import os
from urllib.parse import urlparse

from sqlalchemy import URL, create_engine
from sqlalchemy.orm import Session, sessionmaker


def get_database_url() -> URL:
    """Get database URL from environment variables."""
    database_url = os.getenv("DATABASE_URL")
    if database_url:
        # Parse the DATABASE_URL into components
        parsed = urlparse(database_url)
        return URL.create(
            drivername="postgresql+psycopg2",
            username=parsed.username,
            password=parsed.password,
            host=parsed.hostname,
            port=parsed.port or 5432,
            database=parsed.path.lstrip("/"),
        )
    
    # Fall back to individual environment variables
    return URL.create(
        drivername="postgresql+psycopg2",
        username=os.getenv("POSTGRES_USER", "liive"),
        password=os.getenv("POSTGRES_PASSWORD", "liive"),
        host=os.getenv("POSTGRES_HOST", "localhost"),
        port=int(os.getenv("POSTGRES_PORT", "5432")),
        database=os.getenv("POSTGRES_DB", "liivedb"),
    )


def create_session_factory() -> sessionmaker[Session]:
    """Create a session factory for database connections."""
    engine = create_engine(
        get_database_url(),
        echo=bool(os.getenv("SQL_ECHO", "")),
        pool_pre_ping=True,
    )
    return sessionmaker(
        bind=engine,
        expire_on_commit=False,
        autoflush=False,
    ) 