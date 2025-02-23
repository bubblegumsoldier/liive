from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field

from liive_sql_models.models.chat import ParticipantRole


class ChatParticipantBase(BaseModel):
    user_id: int
    role: ParticipantRole


class ChatParticipantCreate(ChatParticipantBase):
    pass


class ChatParticipantUpdate(BaseModel):
    role: ParticipantRole = Field(description="New role for the participant")


class ChatParticipantResponse(ChatParticipantBase):
    id: int
    chat_id: int
    message: Optional[str] = None
    joined_at: datetime
    message_updated_at: Optional[datetime] = None

    class Config:
        from_attributes = True


class ChatBase(BaseModel):
    is_one_on_one: bool = Field(description="Whether this is a 1-on-1 chat or a group chat")


class ChatCreate(ChatBase):
    participants: List[ChatParticipantCreate]


class ChatResponse(ChatBase):
    id: int
    created_at: datetime
    updated_at: datetime
    participants: List[ChatParticipantResponse]

    class Config:
        from_attributes = True


class ChatList(BaseModel):
    chats: List[ChatResponse] 