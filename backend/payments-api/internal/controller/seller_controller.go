package controller

import (
	"net/http"

	"payments-api/internal/domain"
	"payments-api/internal/middleware"
	"payments-api/internal/service"

	"github.com/gin-gonic/gin"
)

// SellerController handles seller account endpoints
type SellerController struct {
	sellerService service.SellerService
}

// NewSellerController creates a new seller controller
func NewSellerController(sellerService service.SellerService) *SellerController {
	return &SellerController{
		sellerService: sellerService,
	}
}

// LinkManualTokenRequest represents manual token linking request
type LinkManualTokenRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// LinkManualToken links a seller account using manually provided access token
// @Summary Link account with manual token
// @Description Links seller account using manually provided MercadoPago access token
// @Tags seller
// @Accept json
// @Produce json
// @Param request body LinkManualTokenRequest true "Access token"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/seller/link-manual [post]
func (c *SellerController) LinkManualToken(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	var req LinkManualTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "access_token es requerido",
		})
		return
	}

	err := c.sellerService.LinkManualToken(userID, req.AccessToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrSellerAlreadyLinked {
			status = http.StatusConflict
		} else if err == domain.ErrInvalidToken {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cuenta de Mercado Pago vinculada exitosamente",
	})
}

// GetOAuthURL returns the URL to link MP account
// @Summary Get OAuth URL
// @Description Returns the Mercado Pago OAuth URL for account linking
// @Tags seller
// @Produce json
// @Param redirect_uri query string true "Redirect URI after OAuth"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/seller/oauth/url [get]
func (c *SellerController) GetOAuthURL(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "redirect_uri es requerido",
		})
		return
	}

	url, err := c.sellerService.GetOAuthURL(userID, redirectURI)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrSellerAlreadyLinked {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"oauth_url": url,
	})
}

// OAuthCallbackRequest represents the OAuth callback request
type OAuthCallbackRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
}

// OAuthCallback processes the OAuth callback
// @Summary Process OAuth callback
// @Description Processes the Mercado Pago OAuth callback
// @Tags seller
// @Accept json
// @Produce json
// @Param request body OAuthCallbackRequest true "OAuth callback data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/seller/oauth/callback [post]
func (c *SellerController) OAuthCallback(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	var req OAuthCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	err := c.sellerService.ProcessOAuthCallback(userID, req.Code, req.RedirectURI)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrOAuthFailed {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cuenta de Mercado Pago vinculada exitosamente",
	})
}

// GetStatus returns the seller account status
// @Summary Get seller status
// @Description Returns the status of the user's seller account
// @Tags seller
// @Produce json
// @Success 200 {object} service.SellerStatusResponse
// @Router /api/v1/seller/status [get]
func (c *SellerController) GetStatus(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	resp, err := c.sellerService.GetStatus(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// Disconnect removes the MP account link
// @Summary Disconnect seller account
// @Description Removes the Mercado Pago account link
// @Tags seller
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/seller/disconnect [delete]
func (c *SellerController) Disconnect(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	err := c.sellerService.Disconnect(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrSellerNotFound {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cuenta de Mercado Pago desvinculada",
	})
}
