package types

type CreateChatRequest struct {
	Title     string `json:"title"`
	MemberIDs []uint `json:"member_ids" validate:"required,min=1"`
}

type UpdateChatTitleRequest struct {
	Title string `json:"title" validate:"required,min=1"`
}

type AddMembersRequest struct {
	MemberIDs []uint `json:"member_ids" validate:"required,min=1"`
}

type ChatMemberResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	JoinedAt  string `json:"joined_at"`
	LeftAt    string `json:"left_at,omitempty"`
}

type ChatResponse struct {
	ID        uint               `json:"id"`
	Title     string            `json:"title,omitempty"`
	IsGroup   bool              `json:"is_group"`
	CreatedAt string            `json:"created_at"`
	Members   []ChatMemberResponse `json:"members"`
}

type ErrorResponse struct {
	Error string `json:"error"`
} 