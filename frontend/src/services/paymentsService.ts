/**
 * Payments Service for CarPooling Payments API (Port 8005)
 * Handles MercadoPago preferences, wallet, withdrawals, and seller accounts
 */

import { paymentsApi, tripsApi } from './api'
import type {
  CreatePreferenceRequest,
  CreatePreferenceResponse,
  PaymentResponse,
  WalletResponse,
  WalletSettingsRequest,
  TransactionsListResponse,
  WalletStatsResponse,
  WithdrawalRequest,
  WithdrawalResponse,
  WithdrawalsListResponse,
  SellerStatusResponse,
} from '@/types/payment'
import type { ApiResponse } from '@/types'

// ─── Payments ───────────────────────────────────────────────────────────────

export async function createPreference(
  data: CreatePreferenceRequest
): Promise<CreatePreferenceResponse> {
  const response = await paymentsApi.post<ApiResponse<CreatePreferenceResponse>>(
    '/payments/preference',
    data
  )
  return response.data.data!
}

export async function getPaymentByBooking(bookingId: string): Promise<PaymentResponse> {
  const response = await paymentsApi.get<ApiResponse<PaymentResponse>>(
    `/payments/by-booking?booking_id=${encodeURIComponent(bookingId)}`
  )
  return response.data.data!
}

export async function getPayment(paymentUuid: string): Promise<PaymentResponse> {
  const response = await paymentsApi.get<ApiResponse<PaymentResponse>>(
    `/payments/${paymentUuid}`
  )
  return response.data.data!
}

export async function releasePayment(paymentUuid: string): Promise<void> {
  await paymentsApi.post(`/payments/${paymentUuid}/release`, {
    payment_uuid: paymentUuid,
  })
}

// ─── Wallet ──────────────────────────────────────────────────────────────────

export async function getWallet(): Promise<WalletResponse> {
  const response = await paymentsApi.get<ApiResponse<WalletResponse>>('/wallet')
  return response.data.data!
}

export async function getTransactions(page = 1, pageSize = 20): Promise<TransactionsListResponse> {
  const response = await paymentsApi.get<ApiResponse<TransactionsListResponse>>(
    `/wallet/transactions?page=${page}&page_size=${pageSize}`
  )
  return response.data.data!
}

export async function getWalletStats(): Promise<WalletStatsResponse> {
  const response = await paymentsApi.get<ApiResponse<WalletStatsResponse>>('/wallet/stats')
  return response.data.data!
}

export async function updateWalletSettings(data: WalletSettingsRequest): Promise<WalletResponse> {
  const response = await paymentsApi.patch<ApiResponse<WalletResponse>>('/wallet/settings', data)
  return response.data.data!
}

// ─── Withdrawals ─────────────────────────────────────────────────────────────

export async function requestWithdrawal(amountCentavos: number): Promise<WithdrawalResponse> {
  const data: WithdrawalRequest = { amount: amountCentavos }
  const response = await paymentsApi.post<ApiResponse<WithdrawalResponse>>('/withdrawals', data)
  return response.data.data!
}

export async function getWithdrawals(page = 1): Promise<WithdrawalsListResponse> {
  const response = await paymentsApi.get<ApiResponse<WithdrawalsListResponse>>(
    `/withdrawals?page=${page}`
  )
  return response.data.data!
}

export async function getWithdrawal(withdrawalUuid: string): Promise<WithdrawalResponse> {
  const response = await paymentsApi.get<ApiResponse<WithdrawalResponse>>(
    `/withdrawals/${withdrawalUuid}`
  )
  return response.data.data!
}

// ─── Seller Account ──────────────────────────────────────────────────────────

export async function getSellerStatus(): Promise<SellerStatusResponse> {
  const response = await paymentsApi.get<ApiResponse<SellerStatusResponse>>('/seller/status')
  return response.data.data!
}

export async function getOAuthURL(): Promise<string> {
  const response = await paymentsApi.get<ApiResponse<{ url: string }>>('/seller/oauth/url')
  return response.data.data!.url
}

export async function linkManualToken(accessToken: string): Promise<void> {
  await paymentsApi.post('/seller/link-manual', { access_token: accessToken })
}

export async function disconnectSeller(): Promise<void> {
  await paymentsApi.delete('/seller/disconnect')
}

// ─── Route embed (Google Maps, server-side key) — served by trips-api ────────

export async function getRouteEmbedURL(
  originLat: number,
  originLng: number,
  destLat: number,
  destLng: number
): Promise<string> {
  const response = await tripsApi.get<{ success: boolean; embed_url: string }>(
    `/route?origin_lat=${originLat}&origin_lng=${originLng}&dest_lat=${destLat}&dest_lng=${destLng}`
  )
  return response.data.embed_url
}

export default {
  createPreference,
  getPaymentByBooking,
  getPayment,
  releasePayment,
  getWallet,
  getTransactions,
  getWalletStats,
  updateWalletSettings,
  requestWithdrawal,
  getWithdrawals,
  getWithdrawal,
  getSellerStatus,
  getOAuthURL,
  linkManualToken,
  disconnectSeller,
  getRouteEmbedURL,
}
