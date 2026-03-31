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

func (h *Handler) Register(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

func (h *Handler) GetUsers(c *gin.Context) {
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
	limit := int32(20)
	offset := int32(0)
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = int32(l)
		}
	}
	if oStr := c.Query("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o > 0 {
			offset = int32(o)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := h.UserClient.GetUsers(ctx, token, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка пользователей"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
