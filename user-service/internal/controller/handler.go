package controller

import (
	"user-service/internal/controller/responder"
	_ "user-service/internal/models"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
	UploadAvatar(ctx *gin.Context)
	ListUsers(ctx *gin.Context)
}

type HandlerObj struct {
	resp *responder.ResponderObj
}

func NewHandlerObj() *HandlerObj {
	return &HandlerObj{resp: responder.NewResponderObj()}
}

// @Summary 		Регистрация нового пользователя
// @Description 	Добавление  пользователя в базу
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			Users	body	models.User true "Новый пользователь"
// @Success 		200		{object}	models.Response
// @Failure			400		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/register		[post]
func (h *HandlerObj) RegisterUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Аутентификация пользователя
// @Description 	Вход в систему по почте и паролю
// @Tags 			Users
// @Accept			mpfd
// @Produce 		json
// @Param			email	query	string true "email"
// @Param			pass	query	string true "password"
// @Success 		200		{object}	models.LoginResponse
// @Failure			400		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/login		[post]
func (h *HandlerObj) LoginUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Получение информации о пользователе
// @Description 	Получение информации о пользователе по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			id	path	int true "ID"
// @Success 		200		{object}	models.Response
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}		[get]
func (h *HandlerObj) GetUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Обновление информации о пользователе
// @Description 	Обновление информации о пользователе по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			id	path	int true "ID"
// @Success 		200		{object}	models.User
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}		[put]
func (h *HandlerObj) UpdateUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Удаление информации о пользователе
// @Description 	Удаление информации о пользователе по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			id	path	int true "ID"
// @Success 		200		{object}	models.Response
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}		[delete]
func (h *HandlerObj) DeleteUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Получение профиля пользователя
// @Description 	Получение профиля пользователя по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			id	path	int true "ID"
// @Success 		200		{object}	models.UserProfile
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}/profile		[get]
func (h *HandlerObj) GetUserProfile(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Обновление профиля пользователя
// @Description 	Обновление профиля пользователя по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			id	path	int true "ID"
// @Success 		200		{object}	models.Response
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}/profile		[put]
func (h *HandlerObj) UpdateUserProfile(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Загрузка аватара
// @Description 	Загрузка аватара по id
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			file	formData	file true "Avatar"
// @Success 		200		{object}	models.Response
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/{id}/avatar	[post]
func (h *HandlerObj) UploadAvatar(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Получение списка всех пользователей
// @Description 	Получение списка всех пользователей
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Success 		200		{array}	models.UserProfile
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/	[get]
func (h *HandlerObj) ListUsers(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}

// @Summary 		Поиск пользователя по username
// @Description 	Поиск пользователя по username
// @Tags 			Users
// @Accept			json
// @Produce 		json
// @Param			username	query	string true "username"
// @Success 		200		{array}	models.UserProfile
// @Failure			400		{object}	models.Response
// @Failure			403		{object}	models.Response
// @Failure			500		{object}	models.Response
// @Router 			/users/search	[get]
func (h *HandlerObj) SearchUser(ctx *gin.Context) {
	h.resp.SuccessJSON(ctx, "Not implemented!")
}
