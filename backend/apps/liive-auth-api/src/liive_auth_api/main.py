import datetime
from uuid import uuid4
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

from liive_auth_api.config import settings
from liive_auth_api.schemas import TokenResponse, UserCreate, UserResponse
from liive_auth_api.security import (
    create_access_token,
    get_current_user,
    get_db,
    get_password_hash,
    verify_password,
)
from liive_common_api.app import create_app
from liive_sql_models.models.user import User

# Create FastAPI application with common settings
app = create_app(
    title=settings.api_title,
    description=settings.api_description,
    version=settings.api_version,
    cors_origins=settings.cors_origins,
    tags_metadata=[
        {
            "name": "auth",
            "description": "Authentication endpoints for user login and registration",
        },
        {
            "name": "users",
            "description": "User management endpoints",
        },
    ],
)


@app.post("/token", response_model=TokenResponse, tags=["auth"])
async def login(
    form_data: OAuth2PasswordRequestForm = Depends(),
    db: Session = Depends(get_db),
) -> TokenResponse:
    """Login to get an access token.

    Args:
        form_data: OAuth2 form containing username (email) and password
        db: Database session

    Returns:
        Access token for the authenticated user

    Raises:
        HTTPException: If credentials are invalid
    """
    # Use username field for email (OAuth2 form standard)
    user = db.query(User).filter(User.email == form_data.username).first()
    if not user or not verify_password(form_data.password, user.password_hash):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )

    access_token = create_access_token(data={"sub": str(user.id)})
    return TokenResponse(access_token=access_token)


@app.post("/register", response_model=UserResponse, status_code=201, tags=["auth"])
async def register(user_data: UserCreate, db: Session = Depends(get_db)) -> User:
    """Register a new user.

    Args:
        user_data: User registration data
        db: Database session

    Returns:
        The created user

    Raises:
        HTTPException: If email or username is already taken
    """
    try:
        user = User(
            id=uuid4(),
            email=user_data.email,
            username=user_data.username,
            password_hash=get_password_hash(user_data.password),
            full_name=user_data.full_name,
            created_at=datetime.datetime.now(),
            updated_at=datetime.datetime.now(),
        )
        db.add(user)
        db.commit()
        db.refresh(user)
        return user
    except IntegrityError:
        db.rollback()
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Email or username already registered",
        )


@app.get("/users/me", response_model=UserResponse, tags=["users"])
async def read_users_me(current_user: User = Depends(get_current_user)) -> User:
    """Get current user information.

    Args:
        current_user: The authenticated user (injected by dependency)

    Returns:
        The current user's information
    """
    return current_user
