package handlers

import (
	"jwt-viewer/models"
	"jwt-viewer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTHandler struct {
	jwtService *services.JWTService
}

func NewJWTHandler(jwtService *services.JWTService) *JWTHandler {
	return &JWTHandler{
		jwtService: jwtService,
	}
}

// DecodeHandler handles POST /api/decode
func (h *JWTHandler) DecodeHandler(c *gin.Context) {
	var req models.DecodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Decode the token
	response, err := h.jwtService.DecodeToken(req.Token)
	if err != nil {
		c.JSON(http.StatusOK, response) // Still return 200 with error in body
		return
	}

	// Extract claim info for convenience
	claimInfo := h.jwtService.ExtractClaimInfo(response.Payload)

	c.JSON(http.StatusOK, gin.H{
		"header":     response.Header,
		"payload":    response.Payload,
		"signature":  response.Signature,
		"claim_info": claimInfo,
	})
}

// EncodeHandler handles POST /api/encode
func (h *JWTHandler) EncodeHandler(c *gin.Context) {
	var req models.EncodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Generate the token
	response, err := h.jwtService.EncodeToken(&req)
	if err != nil {
		c.JSON(http.StatusOK, response) // Still return 200 with error in body
		return
	}

	c.JSON(http.StatusOK, response)
}

// VerifyHandler handles POST /api/verify
func (h *JWTHandler) VerifyHandler(c *gin.Context) {
	var req models.VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Verify the token
	response, err := h.jwtService.VerifyToken(&req)
	if err != nil && response == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Extract claim info if claims are available
	var claimInfo *models.ClaimInfo
	if response.Claims != nil {
		claimInfo = h.jwtService.ExtractClaimInfo(response.Claims)
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      response.Valid,
		"message":    response.Message,
		"claims":     response.Claims,
		"claim_info": claimInfo,
		"error":      response.Error,
	})
}
