package controller

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicitapa/cooking-service/internal/models"
	"github.com/nicitapa/cooking-service/internal/service"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	svc *service.AuthService
}

func NewAuthController(svc *service.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

// @Summary Login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.AuthRequest true "Credentials"
// @Success 200 {object} models.TokenResponse
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user := &models.User{ID: 1, Username: req.Username}

	tokens, err := a.svc.GenerateTokens(c.Request.Context(), user)
	if err != nil {
		log.Error().Err(err).Msg("token generation failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// @Summary Refresh
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body models.RefreshRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Router /auth/refresh [post]

func (a *AuthController) Refresh(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	token, err := a.svc.ParseToken(payload.RefreshToken)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	user := &models.User{ID: int64(claims["sub"].(float64)), Username: claims["username"].(string)}

	tokens, err := a.svc.GenerateTokens(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}
