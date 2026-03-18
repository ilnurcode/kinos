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
		if strings.Contains(errMsg, "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
		} else if strings.Contains(errMsg, "телефон") {
			// Возвращаем ошибку валидации телефона
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		} else if strings.Contains(errMsg, "Validation") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		} else {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}
	resp, err := h.UserClient.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid refresh token"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := h.UserClient.GetProfile(ctx, parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
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
		if strings.Contains(errMsg, "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email уже используется"})
		} else {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления роли"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат токена"})
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
