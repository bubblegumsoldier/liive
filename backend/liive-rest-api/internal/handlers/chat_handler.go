package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/liive/backend/liive-rest-api/internal/service"
	"github.com/liive/backend/liive-rest-api/internal/types"
	"github.com/liive/backend/shared/pkg/models"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// CreateChat godoc
// @Summary Create a new chat
// @Description Create a new chat with specified members
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param chat body types.CreateChatRequest true "Chat creation details"
// @Success 201 {object} types.ChatResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats [post]
func (h *ChatHandler) CreateChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req types.CreateChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	// Add creator to members if not already included
	found := false
	for _, id := range req.MemberIDs {
		if id == userID {
			found = true
			break
		}
	}
	if !found {
		req.MemberIDs = append(req.MemberIDs, userID)
	}

	chat, err := h.chatService.CreateChat(c.Request().Context(), userID, req.Title, req.MemberIDs)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "One or more users not found"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create chat"})
		}
	}

	return c.JSON(http.StatusCreated, toChatResponse(chat))
}

// GetUserChats godoc
// @Summary Get user's chats
// @Description Get all chats that the user is a member of
// @Tags chats
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} types.ChatResponse
// @Failure 401 {object} types.ErrorResponse
// @Router /api/chats [get]
func (h *ChatHandler) GetUserChats(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	chats, err := h.chatService.GetUserChats(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get chats"})
	}

	response := make([]types.ChatResponse, len(chats))
	for i, chat := range chats {
		response[i] = toChatResponse(&chat)
	}

	return c.JSON(http.StatusOK, response)
}

// GetChat godoc
// @Summary Get a specific chat
// @Description Get details of a specific chat by ID
// @Tags chats
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Chat ID"
// @Success 200 {object} types.ChatResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats/{id} [get]
func (h *ChatHandler) GetChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid chat ID"})
	}

	chat, err := h.chatService.GetChat(c.Request().Context(), uint(chatID), userID)
	if err != nil {
		switch err {
		case service.ErrChatNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Chat not found"})
		case service.ErrNotChatMember:
			return c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Not a member of this chat"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get chat"})
		}
	}

	return c.JSON(http.StatusOK, toChatResponse(chat))
}

// UpdateChatTitle godoc
// @Summary Update chat title
// @Description Update the title of a group chat
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Chat ID"
// @Param chat body types.UpdateChatTitleRequest true "New chat title"
// @Success 200 {object} types.ChatResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats/{id} [put]
func (h *ChatHandler) UpdateChatTitle(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid chat ID"})
	}

	var req types.UpdateChatTitleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	chat, err := h.chatService.UpdateChatTitle(c.Request().Context(), uint(chatID), userID, req.Title)
	if err != nil {
		switch err {
		case service.ErrChatNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Chat not found"})
		case service.ErrNotChatMember:
			return c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Not a member of this chat"})
		case service.ErrNotGroupChat:
			return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Cannot update title of non-group chat"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update chat title"})
		}
	}

	return c.JSON(http.StatusOK, toChatResponse(chat))
}

// LeaveChat godoc
// @Summary Leave a chat
// @Description Leave a chat you are a member of
// @Tags chats
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Chat ID"
// @Success 204
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats/{id}/leave [post]
func (h *ChatHandler) LeaveChat(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid chat ID"})
	}

	err = h.chatService.LeaveChat(c.Request().Context(), uint(chatID), userID)
	if err != nil {
		switch err {
		case service.ErrChatNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Chat not found"})
		case service.ErrNotChatMember:
			return c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Not a member of this chat"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to leave chat"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

// AddMembers godoc
// @Summary Add members to a group chat
// @Description Add new members to an existing group chat
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Chat ID"
// @Param members body types.AddMembersRequest true "New members to add"
// @Success 200 {object} types.ChatResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats/{id}/members [post]
func (h *ChatHandler) AddMembers(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid chat ID"})
	}

	var req types.AddMembersRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	chat, err := h.chatService.AddMembers(c.Request().Context(), uint(chatID), userID, req.MemberIDs)
	if err != nil {
		switch err {
		case service.ErrChatNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Chat not found"})
		case service.ErrNotChatMember:
			return c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Not a member of this chat"})
		case service.ErrNotGroupChat:
			return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Cannot add members to non-group chat"})
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "One or more users not found"})
		case service.ErrAlreadyMember:
			return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "One or more users are already members"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to add members"})
		}
	}

	return c.JSON(http.StatusOK, toChatResponse(chat))
}

// RemoveMember godoc
// @Summary Remove a member from a group chat
// @Description Remove a member from a group chat
// @Tags chats
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Chat ID"
// @Param userId path int true "User ID to remove"
// @Success 204
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /api/chats/{id}/members/{userId} [delete]
func (h *ChatHandler) RemoveMember(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	chatID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid chat ID"})
	}

	memberID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid user ID"})
	}

	err = h.chatService.RemoveMember(c.Request().Context(), uint(chatID), userID, uint(memberID))
	if err != nil {
		switch err {
		case service.ErrChatNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Chat not found"})
		case service.ErrNotChatMember:
			return c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Not a member of this chat"})
		case service.ErrNotGroupChat:
			return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Cannot remove members from non-group chat"})
		case service.ErrLastMember:
			return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Cannot remove the last member"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to remove member"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func toChatResponse(chat *models.Chat) types.ChatResponse {
	members := make([]types.ChatMemberResponse, len(chat.Members))
	for i, member := range chat.Members {
		var leftAt string
		if member.LeftAt != nil {
			leftAt = member.LeftAt.Format(time.RFC3339)
		}
		members[i] = types.ChatMemberResponse{
			ID:       member.ID,
			UserID:   member.UserID,
			Username: member.User.Username,
			JoinedAt: member.JoinedAt.Format(time.RFC3339),
			LeftAt:   leftAt,
		}
	}

	return types.ChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		IsGroup:   chat.IsGroup,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		Members:   members,
	}
} 