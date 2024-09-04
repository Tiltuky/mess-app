package handler

import (
	_ "messaging-service/docs" // Swagger files
	responder "messaging-service/internal/modules/handler/responder"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	responder responder.Responder
}

func NewHandler(responder responder.Responder) *Handler {
	return &Handler{
		responder: responder,
	}
}

// @Summary Get Messages
// @Description Получение сообщений в чате по ID чата
// @Tags Messages
// @Param chat_id path string true "Chat ID"
// @Success 200 {array} models.Message
// @Failure 404 {object} responder.ErrorResponse
// @Router /messages/{chat_id} [get]
func (h *Handler) getMessages(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Send Message
// @Description Отправка сообщения в чат по ID чата
// @Tags Messages
// @Param chat_id path string true "Chat ID"
// @Param message body models.Message true "Message Input"
// @Success 200 {object} models.Message
// @Failure 404 {object} responder.ErrorResponse
// @Router /messages/{chat_id} [post]
func (h *Handler) sendMessage(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Delete Message
// @Description Удаление сообщения по ID сообщения в чате
// @Tags Messages
// @Param chat_id path string true "Chat ID"
// @Param msg_id path string true "Message ID"
// @Success 	200		{string} string "ok"
// @Failure		400		{object} responder.ErrorResponse
// @Failure		404		{object} responder.ErrorResponse
// @Failure		500		{object} responder.ErrorResponse
// @Router /messages/{chat_id}/{msg_id} [delete]
func (h *Handler) deleteMessage(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Get Chats
// @Description Получение списка чатов текущего пользователя
// @Tags Chats
// @Success 200 {array} models.Chat
// @Failure 404 {object} responder.ErrorResponse
// @Router /chats [get]
func (h *Handler) getChats(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Create Chat
// @Description Создание нового чата
// @Tags Chats
// @Param chat body models.Chat true "Chat Input"
// @Success 201 {object} models.Chat
// @Failure 400 {object} responder.ErrorResponse
// @Router /chats [post]
func (h *Handler) createChat(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Get Chat Information
// @Description Получение информации о чате по ID
// @Tags Chats
// @Param id path string true "Chat ID"
// @Success 200 {object} models.Chat
// @Failure 404 {object} responder.ErrorResponse
// @Router /chats/{id} [get]
func (h *Handler) getChatInfo(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}

// @Summary Delete Chat
// @Description Удаление чата по ID
// @Tags Chats
// @Param id path string true "Chat ID"
// @Success 200 {string} string "ok"
// @Failure 404 {object} responder.ErrorResponse
// @Failure 500 {object} responder.ErrorResponse
// @Router /chats/{id} [delete]
func (h *Handler) deleteChat(c *gin.Context) {
	h.responder.SuccessJSON(c, "Not implemented!")
}
