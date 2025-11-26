package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"trips-api/internal/service"
)

// ChatController handles chat-related HTTP requests
type ChatController struct {
	chatService service.ChatService
}

// NewChatController creates a new chat controller instance
func NewChatController(chatService service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// SendMessage handles POST /trips/:id/messages
// Requires authentication - user info extracted from JWT middleware
func (c *ChatController) SendMessage(ctx *gin.Context) {
	tripID := ctx.Param("id")

	// Get user info from JWT context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		log.Warn().Msg("Unauthorized chat message attempt - no user_id in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized - authentication required",
		})
		return
	}

	userName, _ := ctx.Get("user_name")
	userNameStr, ok := userName.(string)
	if !ok || userNameStr == "" {
		userNameStr = "Anonymous" // Fallback if name not available
	}

	// Parse request body
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for send message")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "message is required",
		})
		return
	}

	// Convert userID to int64
	userIDInt64, ok := userID.(int64)
	if !ok {
		// Try float64 conversion (JSON numbers are sometimes parsed as float64)
		if userIDFloat, ok := userID.(float64); ok {
			userIDInt64 = int64(userIDFloat)
		} else {
			log.Error().Interface("user_id", userID).Msg("Invalid user_id type in context")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "internal server error",
			})
			return
		}
	}

	// Send message using chat service (with concurrent processing)
	message, err := c.chatService.SendMessage(
		ctx.Request.Context(),
		tripID,
		userIDInt64,
		userNameStr,
		req.Message,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Int64("user_id", userIDInt64).
			Msg("Failed to send chat message")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Info().
		Str("message_id", message.ID.Hex()).
		Str("trip_id", tripID).
		Int64("user_id", userIDInt64).
		Msg("Chat message sent successfully")

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    message,
	})
}

// GetMessages handles GET /trips/:id/messages
// Retrieves chat messages for a trip
func (c *ChatController) GetMessages(ctx *gin.Context) {
	tripID := ctx.Param("id")

	// Get messages from chat service
	messages, err := c.chatService.GetMessages(ctx.Request.Context(), tripID)
	if err != nil {
		log.Error().
			Err(err).
			Str("trip_id", tripID).
			Msg("Failed to get chat messages")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	log.Debug().
		Str("trip_id", tripID).
		Int("count", len(messages)).
		Msg("Chat messages retrieved successfully")

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"messages": messages,
		"count":    len(messages),
	})
}
