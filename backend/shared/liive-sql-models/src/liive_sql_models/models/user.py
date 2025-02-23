from datetime import datetime
from typing import Optional
from uuid import UUID, uuid4

from sqlalchemy import String, text
from sqlalchemy.orm import Mapped, mapped_column

from liive_sql_models.base import Base


class User(Base):
    __tablename__ = "users"
    """User model for authentication and profile information."""

    # Required fields without defaults
    email: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    username: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    password_hash: Mapped[str] = mapped_column(String(255), nullable=False)
    created_at: Mapped[datetime] = mapped_column(server_default=text("CURRENT_TIMESTAMP"))
    updated_at: Mapped[datetime] = mapped_column(
        server_default=text("CURRENT_TIMESTAMP"), onupdate=datetime.utcnow
    )

    # Fields with defaults
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    is_active: Mapped[bool] = mapped_column(default=True, server_default=text("true"))

    # Optional fields
    full_name: Mapped[Optional[str]] = mapped_column(String(100), default=None)
