package service

import (
	"context"
	"errors"
	"time"

	"github.com/liive/backend/shared/pkg/models"
	"gorm.io/gorm"
)

var (
	ErrChatNotFound     = errors.New("chat not found")
	ErrNotGroupChat     = errors.New("not a group chat")
	ErrNotChatMember    = errors.New("not a chat member")
	ErrLastMember       = errors.New("cannot remove last member")
	ErrAlreadyMember    = errors.New("user is already a member")
	ErrUserNotFound     = errors.New("user not found")
)

type ChatService struct {
	db *gorm.DB
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{
		db: db,
	}
}

func (s *ChatService) CreateChat(ctx context.Context, creatorID uint, title string, memberIDs []uint) (*models.Chat, error) {
	// Check if all users exist
	var count int64
	if err := s.db.Model(&models.User{}).Where("id IN ?", memberIDs).Count(&count).Error; err != nil {
		return nil, err
	}
	if int(count) != len(memberIDs) {
		return nil, ErrUserNotFound
	}

	// Create chat
	chat := &models.Chat{
		Title:   title,
		IsGroup: len(memberIDs) > 2,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(chat).Error; err != nil {
			return err
		}

		// Create chat members
		members := make([]models.ChatMember, len(memberIDs))
		for i, userID := range memberIDs {
			members[i] = models.ChatMember{
				ChatID:   chat.ID,
				UserID:   userID,
				JoinedAt: time.Now(),
			}
		}

		if err := tx.Create(&members).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load members
	if err := s.db.Preload("Members.User").First(chat, chat.ID).Error; err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uint) ([]models.Chat, error) {
	var chats []models.Chat
	err := s.db.Joins("JOIN chat_members ON chat_members.chat_id = chats.id").
		Where("chat_members.user_id = ? AND chat_members.left_at IS NULL", userID).
		Preload("Members.User").
		Find(&chats).Error
	return chats, err
}

func (s *ChatService) GetChat(ctx context.Context, chatID, userID uint) (*models.Chat, error) {
	var chat models.Chat
	err := s.db.Preload("Members.User").First(&chat, chatID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}

	// Check if user is a member
	isMember := false
	for _, member := range chat.Members {
		if member.UserID == userID && member.LeftAt == nil {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, ErrNotChatMember
	}

	return &chat, nil
}

func (s *ChatService) UpdateChatTitle(ctx context.Context, chatID, userID uint, title string) (*models.Chat, error) {
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	if !chat.IsGroup {
		return nil, ErrNotGroupChat
	}

	chat.Title = title
	if err := s.db.Save(chat).Error; err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) LeaveChat(ctx context.Context, chatID, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var member models.ChatMember
		err := tx.Where("chat_id = ? AND user_id = ? AND left_at IS NULL", chatID, userID).
			First(&member).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotChatMember
			}
			return err
		}

		// Check if last member
		var activeMembers int64
		if err := tx.Model(&models.ChatMember{}).
			Where("chat_id = ? AND left_at IS NULL", chatID).
			Count(&activeMembers).Error; err != nil {
			return err
		}

		if activeMembers == 1 {
			// Last member, soft delete the chat
			if err := tx.Delete(&models.Chat{}, chatID).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		member.LeftAt = &now
		return tx.Save(&member).Error
	})
}

func (s *ChatService) AddMembers(ctx context.Context, chatID, userID uint, newMemberIDs []uint) (*models.Chat, error) {
	var resultChat *models.Chat
	err := s.db.Transaction(func(tx *gorm.DB) error {
		chat, err := s.GetChat(ctx, chatID, userID)
		if err != nil {
			return err
		}

		if !chat.IsGroup {
			return ErrNotGroupChat
		}

		// Check if users exist
		var count int64
		if err := tx.Model(&models.User{}).Where("id IN ?", newMemberIDs).Count(&count).Error; err != nil {
			return err
		}
		if int(count) != len(newMemberIDs) {
			return ErrUserNotFound
		}

		// Check for existing members
		var existingCount int64
		if err := tx.Model(&models.ChatMember{}).
			Where("chat_id = ? AND user_id IN ? AND left_at IS NULL", chatID, newMemberIDs).
			Count(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			return ErrAlreadyMember
		}

		// Add new members
		members := make([]models.ChatMember, len(newMemberIDs))
		for i, userID := range newMemberIDs {
			members[i] = models.ChatMember{
				ChatID:   chatID,
				UserID:   userID,
				JoinedAt: time.Now(),
			}
		}

		if err := tx.Create(&members).Error; err != nil {
			return err
		}

		// Reload chat with members
		if err := tx.Preload("Members.User").First(chat, chatID).Error; err != nil {
			return err
		}

		resultChat = chat
		return nil
	})

	return resultChat, err
}

func (s *ChatService) RemoveMember(ctx context.Context, chatID, userID, memberToRemoveID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		chat, err := s.GetChat(ctx, chatID, userID)
		if err != nil {
			return err
		}

		if !chat.IsGroup {
			return ErrNotGroupChat
		}

		// Check if member exists
		var member models.ChatMember
		err = tx.Where("chat_id = ? AND user_id = ? AND left_at IS NULL", chatID, memberToRemoveID).
			First(&member).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotChatMember
			}
			return err
		}

		// Check if last member
		var activeMembers int64
		if err := tx.Model(&models.ChatMember{}).
			Where("chat_id = ? AND left_at IS NULL", chatID).
			Count(&activeMembers).Error; err != nil {
			return err
		}

		if activeMembers == 1 {
			return ErrLastMember
		}

		now := time.Now()
		member.LeftAt = &now
		return tx.Save(&member).Error
	})
} 