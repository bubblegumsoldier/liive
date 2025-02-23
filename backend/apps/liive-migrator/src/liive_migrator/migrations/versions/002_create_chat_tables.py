"""Create chat tables

Revision ID: 002
Revises: 001
Create Date: 2024-02-23

"""

from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

# revision identifiers, used by Alembic.
revision: str = "002"
down_revision: Union[str, None] = "001"
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Create chat table
    op.create_table(
        "chats",
        sa.Column("id", postgresql.UUID(), nullable=False),
        sa.Column("is_one_on_one", sa.Boolean(), nullable=False),
        sa.Column("created_at", sa.DateTime(), nullable=False),
        sa.Column("updated_at", sa.DateTime(), nullable=False),
        sa.PrimaryKeyConstraint("id"),
    )

    # Create chat_participants table
    op.create_table(
        "chat_participants",
        sa.Column("id", postgresql.UUID(), nullable=False),
        sa.Column("chat_id", postgresql.UUID(), nullable=False),
        sa.Column("user_id", postgresql.UUID(), nullable=False),
        sa.Column("role", sa.Enum("admin", "member", name="participantrole"), nullable=False),
        sa.Column("message", sa.Text(), nullable=True),
        sa.Column("joined_at", sa.DateTime(), nullable=False),
        sa.Column("message_updated_at", sa.DateTime(), nullable=True),
        sa.ForeignKeyConstraint(["chat_id"], ["chats.id"], ondelete="CASCADE"),
        sa.ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("id"),
    )

    # Create indexes
    op.create_index("ix_chat_participants_chat_id", "chat_participants", ["chat_id"])
    op.create_index("ix_chat_participants_user_id", "chat_participants", ["user_id"])


def downgrade() -> None:
    # Drop indexes
    op.drop_index("ix_chat_participants_user_id")
    op.drop_index("ix_chat_participants_chat_id")

    # Drop tables
    op.drop_table("chat_participants")
    op.drop_table("chats")
