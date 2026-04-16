// Package users предоставляет HTTP-обработчики для управления пользователями и аутентификацией.
// Включает регистрацию, вход, обновление профиля и администрирование пользователей.
package users

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserClient UserClientInterface
}

func NewHandler(uc UserClientInterface) *Handler {
	return &Handler{
		UserClient: uc,
	}
}

func maxAgeFromUnix(ts int64) int {
	exp := time.Unix(ts, 0)
	secs := int(time.Until(exp).Seconds())
	if secs < 0 {
		return 0
	}
	return secs
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя с email, паролем и телефоном
// @Tags users
// @Accept json
// @Produce json
// @Param username body string true "Имя пользователя"
// @Param email body string true "Email"
// @Param password body string true "Пароль"
// @Param phone body string true "Телефон в формате E.164"
// @Success 200 {object} object{access_token=string} "JWT токен доступа"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Failure 409 {object} object{error=string} "Пользователь уже существует"
// @Router /api/users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
		return
	}
	resp, err := h.UserClient.Register(c.Request.Context(), body.Username, body.Email, body.Password, body.Phone)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "уже существует"):
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		case strings.Contains(errMsg, "ошибка валидации"):
			// Извлекаем сообщение об ошибке валидации
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		case strings.Contains(errMsg, "телефон"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		case strings.Contains(errMsg, "email"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации"})
		}
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: resp.RefreshToken, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, MaxAge: maxAgeFromUnix(resp.RefreshExpiresAt)})
	c.JSON(http.StatusOK, gin.H{"access_token": resp.AccessToken})
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя по email и паролю
// @Tags users
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Пароль"
// @Success 200 {object} object{access_token=string} "JWT токен доступа"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Failure 401 {object} object{error=string} "Неверный email или пароль"
// @Router /api/users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса"})
		return
	}
	resp, err := h.UserClient.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: resp.RefreshToken, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, MaxAge: maxAgeFromUnix(resp.RefreshExpiresAt)})
	c.JSON(http.StatusOK, gin.H{"access_token": resp.AccessToken})
}

// Revoke godoc
// @Summary Выход пользователя
// @Description Отзыв токена обновления и выход из системы
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /api/users/revoke [post]
func (h *Handler) Revoke(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: "", Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, MaxAge: -1})
	if refreshToken != "" {
		_, err := h.UserClient.RevokeRefreshToken(c.Request.Context(), refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
	c.Status(http.StatusOK)
}

// Refresh godoc
// @Summary Обновление токена доступа
// @Description Обновление токена доступа с использованием токена обновления
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} object{access_token=string} "Новый JWT токен доступа"
// @Failure 401 {object} object{error=string} "Неверный или истекший токен обновления"
// @Router /api/users/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	resp, err := h.UserClient.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или истекший токен обновления"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: resp.RefreshToken, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, MaxAge: maxAgeFromUnix(resp.RefreshExpiresAt)})
	c.JSON(http.StatusOK, gin.H{"access_token": resp.AccessToken})
}

// GetProfile godoc
// @Summary Получить профиль пользователя
// @Description Получить информацию о профиле текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Информация о пользователе"
// @Failure 401 {object} object{error=string} "Требуется аутентификация"
// @Router /api/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена. Используйте: Bearer <token>"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.UserClient.GetProfile(ctx, parts[1])
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "требуется аутентификация") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		} else if strings.Contains(errMsg, "пользователь не найден") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка при получении профиля"})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateProfile godoc
// @Summary Обновить профиль пользователя
// @Description Обновить информацию о профиле текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username body string false "Имя пользователя"
// @Param email body string false "Email"
// @Param phone body string false "Телефон"
// @Success 200 {object} object "Обновленный профиль"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Failure 401 {object} object{error=string} "Требуется аутентификация"
// @Failure 404 {object} object{error=string} "Пользователь не найден"
// @Failure 409 {object} object{error=string} "Email уже используется"
// @Router /api/profile [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена. Используйте: Bearer <token>"})
		return
	}
	token := parts[1]
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	updatedProfile, err := h.UserClient.UpdateProfile(ctx, token, req.Username, req.Email, req.Phone)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "уже используется"):
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
		case strings.Contains(errMsg, "пользователь не найден"):
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
		}
		return
	}
	c.JSON(http.StatusOK, updatedProfile)
}

// UpdateRole godoc
// @Summary Обновить роль пользователя
// @Description Обновить роль пользователя (только для администраторов)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id body uint64 true "ID пользователя"
// @Param role body string true "Роль (admin или user)"
// @Success 200 {object} object "Обновленная роль"
// @Failure 400 {object} object{error=string} "Некорректные данные"
// @Failure 401 {object} object{error=string} "Требуется аутентификация"
// @Failure 403 {object} object{error=string} "Доступ запрещен"
// @Failure 404 {object} object{error=string} "Пользователь не найден"
// @Router /api/admin/users/role [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется аутентификация"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена. Используйте: Bearer <token>"})
		return
	}
	token := parts[1]
	var req struct {
		UserID uint64 `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	resp, err := h.UserClient.UpdateRole(ctx, token, req.Role, req.UserID)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "доступ запрещен"):
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен: требуется роль администратора"})
		case strings.Contains(errMsg, "нельзя изменить собственную роль"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя изменить собственную роль"})
		case strings.Contains(errMsg, "пользователь не найден"):
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		case strings.Contains(errMsg, "недопустимая роль"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимая роль. Разрешены только 'admin' или 'user'"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления роли"})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена. Используйте: Bearer <token>"})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный id пользователя"})
		return
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	resp, err := h.UserClient.DeleteUser(ctx, parts[1], userID)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "доступ запрещен"):
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен: требуется роль администратора"})
		case strings.Contains(errMsg, "собственную учетную запись"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Нельзя удалить собственную учетную запись"})
		case strings.Contains(errMsg, "пользователь не найден"):
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления пользователя"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUsers godoc
// @Summary Получить список пользователей
// @Description Получить список всех пользователей с пагинацией (только для администраторов)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int32 false "Лимит записей" default(20)
// @Param offset query int32 false "Смещение" default(0)
// @Success 200 {object} object "Список пользователей"
// @Failure 401 {object} object{error=string} "Требуется аутентификация"
// @Router /api/admin/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр limit"})
		return
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр offset"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
		return
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	resp, err := h.UserClient.GetUsers(ctx, parts[1], int32(limit), int32(offset))
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "доступ запрещен"):
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		case strings.Contains(errMsg, "требуется аутентификация"):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка пользователей"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": resp.Users,
		"total": resp.Total,
	})
}

func requestContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), 10*time.Second)
}
