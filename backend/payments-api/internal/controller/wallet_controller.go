package controller

import (
	"net/http"
	"strconv"

	"payments-api/internal/domain"
	"payments-api/internal/middleware"
	"payments-api/internal/service"

	"github.com/gin-gonic/gin"
)

// WalletController handles wallet endpoints
type WalletController struct {
	walletService     service.WalletService
	withdrawalService service.WithdrawalService
}

// NewWalletController creates a new wallet controller
func NewWalletController(
	walletService service.WalletService,
	withdrawalService service.WithdrawalService,
) *WalletController {
	return &WalletController{
		walletService:     walletService,
		withdrawalService: withdrawalService,
	}
}

// GetWallet retrieves the user's wallet
// @Summary Get wallet
// @Description Retrieves the current user's wallet
// @Tags wallet
// @Produce json
// @Success 200 {object} domain.WalletResponse
// @Router /api/v1/wallet [get]
func (c *WalletController) GetWallet(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	resp, err := c.walletService.GetWallet(userID)
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

// GetTransactions retrieves wallet transactions
// @Summary Get transactions
// @Description Retrieves wallet transaction history
// @Tags wallet
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} domain.TransactionsListResponse
// @Router /api/v1/wallet/transactions [get]
func (c *WalletController) GetTransactions(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	resp, err := c.walletService.GetTransactions(userID, page, pageSize)
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

// GetStats retrieves wallet statistics
// @Summary Get wallet stats
// @Description Retrieves wallet statistics
// @Tags wallet
// @Produce json
// @Success 200 {object} domain.WalletStatsResponse
// @Router /api/v1/wallet/stats [get]
func (c *WalletController) GetStats(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	resp, err := c.walletService.GetStats(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrWalletNotFound {
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
		"data":    resp,
	})
}

// UpdateSettings updates wallet settings
// @Summary Update wallet settings
// @Description Updates wallet auto-withdraw settings
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body domain.WalletSettingsRequest true "Settings"
// @Success 200 {object} domain.WalletResponse
// @Router /api/v1/wallet/settings [patch]
func (c *WalletController) UpdateSettings(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	var req domain.WalletSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	resp, err := c.walletService.UpdateSettings(userID, &req)
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

// RequestWithdrawal creates a withdrawal request
// @Summary Request withdrawal
// @Description Requests a withdrawal from the wallet
// @Tags withdrawals
// @Accept json
// @Produce json
// @Param request body domain.WithdrawalRequest true "Withdrawal request"
// @Success 201 {object} domain.WithdrawalResponse
// @Router /api/v1/withdrawals [post]
func (c *WalletController) RequestWithdrawal(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	var req domain.WithdrawalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "datos inválidos: " + err.Error(),
		})
		return
	}

	resp, err := c.withdrawalService.RequestWithdrawal(userID, req.Amount)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrInsufficientBalance:
			status = http.StatusBadRequest
		case domain.ErrWithdrawalInProgress:
			status = http.StatusConflict
		case domain.ErrSellerNotActive:
			status = http.StatusPreconditionFailed
		}
		ctx.JSON(status, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetWithdrawals retrieves user's withdrawals
// @Summary Get withdrawals
// @Description Retrieves withdrawal history
// @Tags withdrawals
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} domain.WithdrawalsListResponse
// @Router /api/v1/withdrawals [get]
func (c *WalletController) GetWithdrawals(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	resp, err := c.withdrawalService.GetWithdrawals(userID, page, pageSize)
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

// GetWithdrawal retrieves a specific withdrawal
// @Summary Get withdrawal
// @Description Retrieves a specific withdrawal by UUID
// @Tags withdrawals
// @Produce json
// @Param id path string true "Withdrawal UUID"
// @Success 200 {object} domain.WithdrawalResponse
// @Router /api/v1/withdrawals/{id} [get]
func (c *WalletController) GetWithdrawal(ctx *gin.Context) {
	userID, ok := middleware.GetUserIdentifier(ctx)
	if !ok || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "usuario no autenticado",
		})
		return
	}

	withdrawalUUID := ctx.Param("id")

	resp, err := c.withdrawalService.GetWithdrawal(userID, withdrawalUUID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrWithdrawalNotFound:
			status = http.StatusNotFound
		case domain.ErrForbidden:
			status = http.StatusForbidden
		}
		ctx.JSON(status, gin.H{
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
