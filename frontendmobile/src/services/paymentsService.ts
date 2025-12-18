/**
 * Payments Service for CarPooling Mobile App
 * Handles all payment-related API calls to payments-api (port 8005)
 */

import { paymentsApi } from './api'
import type {
  CreatePreferenceRequest,
  CreatePreferenceResponse,
  Payment,
  PaymentResponse,
  Wallet,
  WalletResponse,
  WalletSettingsRequest,
  WalletTransaction,
  TransactionsListResponse,
  WalletStats,
  WalletStatsResponse,
  WithdrawalRequest,
  Withdrawal,
  WithdrawalResponse,
  WithdrawalsListResponse,
  SellerStatusResponse,
  OAuthURLResponse,
  OAuthCallbackRequest,
} from '@/types'

// ============================================================================
// Payment Endpoints
// ============================================================================

/**
 * Create a payment preference for a booking
 * Returns MercadoPago checkout URL (init_point)
 */
export async function createPaymentPreference(
  data: CreatePreferenceRequest
): Promise<CreatePreferenceResponse> {
  const response = await paymentsApi.post<{ success: boolean; data: CreatePreferenceResponse }>(
    '/api/v1/payments/preference',
    data
  )
  return response.data.data
}

/**
 * Get payment details by payment UUID
 */
export async function getPayment(paymentUuid: string): Promise<Payment> {
  const response = await paymentsApi.get<PaymentResponse>(`/api/v1/payments/${paymentUuid}`)
  return response.data.data
}

/**
 * Get payment details by booking ID
 */
export async function getPaymentByBooking(bookingId: string): Promise<Payment> {
  const response = await paymentsApi.get<PaymentResponse>(`/api/v1/payments/booking/${bookingId}`)
  return response.data.data
}

/**
 * Get all payments for current user (as passenger)
 */
export async function getMyPayments(params?: {
  page?: number
  page_size?: number
  status?: string
}): Promise<{ payments: Payment[]; total: number; page: number; total_pages: number }> {
  const response = await paymentsApi.get<{
    success: boolean
    data: { payments: Payment[]; total: number; page: number; page_size: number; total_pages: number }
  }>('/api/v1/payments/my', { params })
  return response.data.data
}

/**
 * Request refund for a payment (passenger action)
 */
export async function requestRefund(paymentUuid: string): Promise<Payment> {
  const response = await paymentsApi.post<PaymentResponse>(`/api/v1/payments/${paymentUuid}/refund`)
  return response.data.data
}

// ============================================================================
// Wallet Endpoints
// ============================================================================

/**
 * Get current user's wallet
 * Creates wallet automatically if it doesn't exist
 */
export async function getMyWallet(): Promise<Wallet> {
  const response = await paymentsApi.get<WalletResponse>('/api/v1/wallet')
  return response.data.data
}

/**
 * Get wallet transaction history
 */
export async function getWalletTransactions(params?: {
  page?: number
  page_size?: number
  type?: string
}): Promise<{ transactions: WalletTransaction[]; total: number; page: number; total_pages: number }> {
  const response = await paymentsApi.get<TransactionsListResponse>('/api/v1/wallet/transactions', { params })
  return response.data.data
}

/**
 * Get wallet statistics (earnings, trips count, etc.)
 */
export async function getWalletStats(): Promise<WalletStats> {
  const response = await paymentsApi.get<WalletStatsResponse>('/api/v1/wallet/stats')
  return response.data.data
}

/**
 * Update wallet settings (auto-withdraw, threshold)
 */
export async function updateWalletSettings(settings: WalletSettingsRequest): Promise<Wallet> {
  const response = await paymentsApi.patch<WalletResponse>('/api/v1/wallet/settings', settings)
  return response.data.data
}

// ============================================================================
// Withdrawal Endpoints
// ============================================================================

/**
 * Request a withdrawal from wallet to linked MercadoPago account
 * Amount is in centavos (e.g., 10000 = $100.00 ARS)
 */
export async function requestWithdrawal(data: WithdrawalRequest): Promise<Withdrawal> {
  const response = await paymentsApi.post<WithdrawalResponse>('/api/v1/withdrawals', data)
  return response.data.data
}

/**
 * Get withdrawal details by UUID
 */
export async function getWithdrawal(withdrawalUuid: string): Promise<Withdrawal> {
  const response = await paymentsApi.get<WithdrawalResponse>(`/api/v1/withdrawals/${withdrawalUuid}`)
  return response.data.data
}

/**
 * Get withdrawal history
 */
export async function getWithdrawals(params?: {
  page?: number
  page_size?: number
  status?: string
}): Promise<{ withdrawals: Withdrawal[]; total: number; page: number; total_pages: number }> {
  const response = await paymentsApi.get<WithdrawalsListResponse>('/api/v1/withdrawals', { params })
  return response.data.data
}

/**
 * Cancel a pending withdrawal
 */
export async function cancelWithdrawal(withdrawalUuid: string): Promise<Withdrawal> {
  const response = await paymentsApi.post<WithdrawalResponse>(`/api/v1/withdrawals/${withdrawalUuid}/cancel`)
  return response.data.data
}

// ============================================================================
// Seller (Driver MercadoPago Account) Endpoints
// ============================================================================

/**
 * Get seller (driver) MercadoPago account link status
 */
export async function getSellerStatus(): Promise<SellerStatusResponse['data']> {
  const response = await paymentsApi.get<SellerStatusResponse>('/api/v1/seller/status')
  return response.data.data
}

/**
 * Link MercadoPago account manually with access token
 * Alternative to OAuth flow - driver provides their MP access token directly
 */
export async function linkSellerManually(accessToken: string): Promise<void> {
  await paymentsApi.post('/api/v1/seller/link-manual', { access_token: accessToken })
}

/**
 * Get OAuth URL to link MercadoPago account
 * Opens this URL in browser for user to authorize
 */
export async function getSellerOAuthURL(redirectUri: string): Promise<string> {
  const response = await paymentsApi.get<OAuthURLResponse>('/api/v1/seller/oauth/url', {
    params: { redirect_uri: redirectUri },
  })
  return response.data.oauth_url
}

/**
 * Complete OAuth flow with authorization code
 * Called after user returns from MercadoPago authorization
 */
export async function completeSellerOAuth(data: OAuthCallbackRequest): Promise<void> {
  await paymentsApi.post('/api/v1/seller/oauth/callback', data)
}

/**
 * Unlink MercadoPago account (revoke authorization)
 */
export async function unlinkSellerAccount(): Promise<void> {
  await paymentsApi.delete('/api/v1/seller/disconnect')
}

// ============================================================================
// Driver Payment Endpoints
// ============================================================================

/**
 * Get payments received as driver
 */
export async function getDriverPayments(params?: {
  page?: number
  page_size?: number
  status?: string
  fund_status?: string
}): Promise<{ payments: Payment[]; total: number; page: number; total_pages: number }> {
  const response = await paymentsApi.get<{
    success: boolean
    data: { payments: Payment[]; total: number; page: number; page_size: number; total_pages: number }
  }>('/api/v1/payments/driver', { params })
  return response.data.data
}

/**
 * Confirm trip completion (releases held funds immediately)
 * Called by passenger after trip ends
 */
export async function confirmTripCompletion(paymentUuid: string): Promise<Payment> {
  const response = await paymentsApi.post<PaymentResponse>(`/api/v1/payments/${paymentUuid}/release`)
  return response.data.data
}

// ============================================================================
// Export all functions
// ============================================================================

export default {
  // Payments
  createPaymentPreference,
  getPayment,
  getPaymentByBooking,
  getMyPayments,
  requestRefund,
  // Wallet
  getMyWallet,
  getWalletTransactions,
  getWalletStats,
  updateWalletSettings,
  // Withdrawals
  requestWithdrawal,
  getWithdrawal,
  getWithdrawals,
  cancelWithdrawal,
  // Seller
  getSellerStatus,
  linkSellerManually,
  getSellerOAuthURL,
  completeSellerOAuth,
  unlinkSellerAccount,
  // Driver
  getDriverPayments,
  confirmTripCompletion,
}
