package domain

import "errors"

// Payment errors
var (
	ErrPaymentNotFound       = errors.New("pago no encontrado")
	ErrPaymentAlreadyExists  = errors.New("ya existe un pago para esta reserva")
	ErrPaymentNotPending     = errors.New("el pago no está pendiente")
	ErrPaymentNotApproved    = errors.New("el pago no está aprobado")
	ErrPaymentAlreadyReleased = errors.New("los fondos ya fueron liberados")
	ErrPaymentAlreadyRefunded = errors.New("el pago ya fue reembolsado")
	ErrCannotReleasePayment  = errors.New("no se pueden liberar los fondos de este pago")
	ErrCannotRefundPayment   = errors.New("no se puede reembolsar este pago")
	ErrInvalidRefundAmount   = errors.New("monto de reembolso inválido")
)

// Wallet errors
var (
	ErrWalletNotFound        = errors.New("billetera no encontrada")
	ErrInsufficientBalance   = errors.New("saldo insuficiente")
	ErrInvalidAmount         = errors.New("monto inválido")
	ErrWalletAlreadyExists   = errors.New("la billetera ya existe")
)

// Withdrawal errors
var (
	ErrWithdrawalNotFound    = errors.New("retiro no encontrado")
	ErrWithdrawalInProgress  = errors.New("ya hay un retiro en proceso")
	ErrWithdrawalFailed      = errors.New("el retiro falló")
	ErrMinimumWithdrawal     = errors.New("el monto mínimo de retiro no se cumple")
)

// Seller errors
var (
	ErrSellerNotFound        = errors.New("cuenta de vendedor no encontrada")
	ErrSellerNotActive       = errors.New("la cuenta de vendedor no está activa")
	ErrSellerAlreadyLinked   = errors.New("ya tienes una cuenta de Mercado Pago vinculada")
	ErrOAuthFailed           = errors.New("error en la autorización de Mercado Pago")
	ErrTokenRefreshFailed    = errors.New("error al refrescar el token de Mercado Pago")
	ErrInvalidToken          = errors.New("token de Mercado Pago inválido")
)

// Mercado Pago errors
var (
	ErrMPPreferenceCreation  = errors.New("error al crear preferencia de pago")
	ErrMPPaymentNotFound     = errors.New("pago no encontrado en Mercado Pago")
	ErrMPRefundFailed        = errors.New("error al procesar reembolso en Mercado Pago")
	ErrMPTransferFailed      = errors.New("error al procesar transferencia en Mercado Pago")
	ErrMPWebhookInvalid      = errors.New("webhook de Mercado Pago inválido")
)

// User errors
var (
	ErrUserNotFound   = errors.New("usuario no encontrado")
	ErrDriverNotFound = errors.New("conductor no encontrado")
)

// General errors
var (
	ErrUnauthorized          = errors.New("no autorizado")
	ErrForbidden             = errors.New("acceso denegado")
	ErrInvalidRequest        = errors.New("solicitud inválida")
	ErrInternalServer        = errors.New("error interno del servidor")
	ErrEventAlreadyProcessed = errors.New("evento ya procesado")
)
