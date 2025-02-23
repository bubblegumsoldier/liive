from datetime import datetime
from enum import Enum
from typing import Optional
from uuid import UUID, uuid4

from sqlalchemy import ForeignKey, Boolean, Text, Enum as SQLEnum
from sqlalchemy.orm import Mapped, mapped_column, relationship

from ..base import Base
from .user import User


class ParticipantRole(str, Enum):
    ADMIN = "admin"
    MEMBER = "member"


class Chat(Base):
    __tablename__ = "chats"

    # Relationships
    participants: Mapped[list["ChatParticipant"]] = relationship(
        back_populates="chat", cascade="all, delete-orphan"
    )

    # Default values
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)

    created_at: Mapped[datetime] = mapped_column(nullable=False, default=datetime.now)
    updated_at: Mapped[datetime] = mapped_column(
        nullable=False, default=datetime.now, onupdate=datetime.now
    )
    is_one_on_one: Mapped[bool] = mapped_column(Boolean, nullable=False, default=False)


class ChatParticipant(Base):
    __tablename__ = "chat_participants"

    chat_id: Mapped[UUID] = mapped_column(
        ForeignKey("chats.id", ondelete="CASCADE"), nullable=False
    )
    user_id: Mapped[UUID] = mapped_column(
        ForeignKey("users.id", ondelete="CASCADE"), nullable=False
    )
    role: Mapped[ParticipantRole] = mapped_column(SQLEnum(ParticipantRole), nullable=False)
    message: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    message_updated_at: Mapped[Optional[datetime]] = mapped_column(nullable=True)

    # Relationships
    chat: Mapped[Chat] = relationship(back_populates="participants")
    user: Mapped[User] = relationship()

    # Default values
    joined_at: Mapped[datetime] = mapped_column(nullable=False, default=datetime.now)
    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
