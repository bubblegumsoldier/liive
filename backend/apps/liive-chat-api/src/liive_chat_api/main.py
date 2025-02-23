"""Chat API service for Liive platform."""

from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy import select
from sqlalchemy.orm import Session

from liive_common_api.auth import get_current_user
from liive_common_api.db import get_db
from liive_sql_models.models.chat import Chat, ChatParticipant, ParticipantRole
from liive_sql_models.models.user import User

from .schemas import (
    ChatCreate,
    ChatResponse,
    ChatList,
    ChatParticipantCreate,
    ChatParticipantUpdate,
    ChatParticipantResponse,
)

app = FastAPI(title="Liive Chat API")


def get_chat_participant(db: Session, chat_id: int, user_id: int) -> ChatParticipant:
    participant = db.scalar(
        select(ChatParticipant).where(
            ChatParticipant.chat_id == chat_id,
            ChatParticipant.user_id == user_id,
        )
    )
    if not participant:
        raise HTTPException(status_code=404, detail="Participant not found")
    return participant


def verify_admin_access(participant: ChatParticipant) -> None:
    if participant.role != ParticipantRole.ADMIN:
        raise HTTPException(status_code=403, detail="Only admins can perform this action")


@app.post("/chats", response_model=ChatResponse)
def create_chat(
    chat_data: ChatCreate,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user),
) -> Chat:
    # Ensure at least one participant is provided
    if not chat_data.participants:
        raise HTTPException(status_code=400, detail="At least one participant must be provided")

    # For 1-on-1 chats, ensure exactly 2 participants
    if chat_data.is_one_on_one and len(chat_data.participants) != 2:
        raise HTTPException(
            status_code=400, detail="One-on-one chats must have exactly 2 participants"
        )

    # Create chat
    chat = Chat(is_one_on_one=chat_data.is_one_on_one)
    db.add(chat)

    # Add participants
    for participant_data in chat_data.participants:
        participant = ChatParticipant(
            chat=chat,
            user_id=participant_data.user_id,
            role=participant_data.role,
        )
        db.add(participant)

    db.commit()
    db.refresh(chat)
    return chat


@app.get("/chats", response_model=ChatList)
def list_chats(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user),
) -> ChatList:
    chats = db.scalars(
        select(Chat).join(ChatParticipant).where(ChatParticipant.user_id == current_user.id)
    ).all()
    return ChatList(chats=chats)


@app.post("/chats/{chat_id}/participants", response_model=ChatParticipantResponse)
def add_participant(
    chat_id: int,
    participant_data: ChatParticipantCreate,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user),
) -> ChatParticipant:
    # Get chat and verify it exists
    chat = db.get(Chat, chat_id)
    if not chat:
        raise HTTPException(status_code=404, detail="Chat not found")

    # Check if chat is one-on-one
    if chat.is_one_on_one:
        raise HTTPException(status_code=400, detail="Cannot add participants to one-on-one chats")

    # Verify current user is an admin
    current_participant = get_chat_participant(db, chat_id, current_user.id)
    verify_admin_access(current_participant)

    # Check if user is already in chat
    existing_participant = db.scalar(
        select(ChatParticipant).where(
            ChatParticipant.chat_id == chat_id,
            ChatParticipant.user_id == participant_data.user_id,
        )
    )
    if existing_participant:
        raise HTTPException(status_code=400, detail="User is already a participant in this chat")

    # Add new participant
    new_participant = ChatParticipant(
        chat_id=chat_id,
        user_id=participant_data.user_id,
        role=participant_data.role,
    )
    db.add(new_participant)
    db.commit()
    db.refresh(new_participant)
    return new_participant


@app.delete("/chats/{chat_id}/participants/{user_id}")
def remove_participant(
    chat_id: int,
    user_id: int,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user),
) -> dict:
    # Get chat and verify it exists
    chat = db.get(Chat, chat_id)
    if not chat:
        raise HTTPException(status_code=404, detail="Chat not found")

    # Check if chat is one-on-one
    if chat.is_one_on_one:
        raise HTTPException(
            status_code=400, detail="Cannot remove participants from one-on-one chats"
        )

    # Verify current user is an admin
    current_participant = get_chat_participant(db, chat_id, current_user.id)
    verify_admin_access(current_participant)

    # Get participant to remove
    participant = get_chat_participant(db, chat_id, user_id)

    # Remove participant
    db.delete(participant)
    db.commit()

    return {"message": "Participant removed successfully"}


@app.patch("/chats/{chat_id}/participants/{user_id}", response_model=ChatParticipantResponse)
def update_participant_role(
    chat_id: int,
    user_id: int,
    participant_data: ChatParticipantUpdate,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user),
) -> ChatParticipant:
    """Update a participant's role in a chat.
    
    Only admins can change roles. The role field is required.
    Message updates are handled by a separate endpoint.
    """
    # Get chat and verify it exists
    chat = db.get(Chat, chat_id)
    if not chat:
        raise HTTPException(status_code=404, detail="Chat not found")

    # Get participant to update
    participant = get_chat_participant(db, chat_id, user_id)

    # Verify current user is an admin
    current_participant = get_chat_participant(db, chat_id, current_user.id)
    verify_admin_access(current_participant)

    # Update participant role
    participant.role = participant_data.role
    db.commit()
    db.refresh(participant)
    return participant
