from uuid import UUID

from pydantic import BaseModel, ConfigDict, EmailStr


class TokenResponse(BaseModel):
    """Token response schema."""

    access_token: str
    token_type: str = "bearer"


class UserBase(BaseModel):
    """Base user schema."""

    email: EmailStr
    username: str
    full_name: str | None = None


class UserCreate(UserBase):
    """User creation schema."""

    password: str


class UserResponse(UserBase):
    """User response schema."""

    id: UUID
    is_active: bool

    model_config = ConfigDict(from_attributes=True) 