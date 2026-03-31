/**
 * Payment and Wallet types for CarPooling payments-api
 * Handles payment processing, wallet management, and withdrawals
 */

// ============================================================================
// Payment Types
// ============================================================================

export type PaymentStatus =
  | 'pending'
  | 'in_process'
  | 'approved'
  | 'rejected'
  | 'cancelled'
  | 'refunded'
  | 'partial_refund'

export type FundStatus = 'held' | 'released' | 'refunded'

export interface CreatePreferenceRequest {
  booking_id: string
  trip_id: string
  driver_id: string
  seats_count: number
  price_per_seat: number // in centavos
  origin: string
  destination: string
  departure_at: string // ISO date
  payer_email?: string
  payer_name?: string
  return_base_url?: string // Mobile: "carpooling:/" (deep link), Web: window.location.origin
}

export interface CreatePreferenceResponse {
  payment_uuid: string
  preference_id: string
  init_point: string
  sandbox_init_point: string
  amount: number
  platform_fee: number
  driver_amount: number
  expires_at: string
}

export interface Payment {
  payment_uuid: string
  booking_id: string
  trip_id: string
  passenger_id: string
  driver_id: string
  amount: number
  currency: string
  platform_fee: number
  driver_amount: number
  seats_count: number
  status: PaymentStatus
  fund_status: FundStatus
  mp_payment_id?: string
  approved_at?: string
  released_at?: string
  created_at: string
}

export interface PaymentResponse {
  success: boolean
  data: Payment
}

// ============================================================================
// Wallet Types
// ============================================================================

export interface Wallet {
  wallet_uuid: string
  user_id: string
  available_balance: number
  pending_balance: number
  held_balance: number
  total_balance: number
  currency: string
  total_earned: number
  total_withdrawn: number
  total_trips: number
  auto_withdraw_enabled: boolean
  auto_withdraw_threshold: number
  available_balance_formatted: string
  total_earned_formatted: string
}

export interface WalletResponse {
  success: boolean
  data: Wallet
}

export interface WalletSettingsRequest {
  auto_withdraw_enabled?: boolean
  auto_withdraw_threshold?: number
}

export type TransactionType =
  | 'credit'
  | 'debit'
  | 'hold'
  | 'release'
  | 'refund'
  | 'withdrawal'
  | 'fee'

export interface WalletTransaction {
  transaction_uuid: string
  type: TransactionType
  amount: number
  currency: string
  balance_before: number
  balance_after: number
  description: string
  payment_id?: string
  withdrawal_id?: string
  created_at: string
  amount_formatted: string
  signed_amount: number
}

export interface TransactionsListResponse {
  success: boolean
  data: {
    transactions: WalletTransaction[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }
}

export interface WalletStats {
  total_earnings: number
  total_trips: number
  total_withdrawals: number
  average_per_trip: number
  this_month_earnings?: number
  last_month_earnings?: number
}

export interface WalletStatsResponse {
  success: boolean
  data: WalletStats
}

// ============================================================================
// Withdrawal Types
// ============================================================================

export type WithdrawalStatus =
  | 'pending'
  | 'processing'
  | 'completed'
  | 'failed'
  | 'cancelled'

export interface WithdrawalRequest {
  amount: number // in centavos
}

export interface Withdrawal {
  withdrawal_uuid: string
  amount: number
  fee: number
  net_amount: number
  currency: string
  status: WithdrawalStatus
  destination_account: string
  error_message?: string
  processed_at?: string
  completed_at?: string
  created_at: string
  amount_formatted: string
}

export interface WithdrawalResponse {
  success: boolean
  data: Withdrawal
}

export interface WithdrawalsListResponse {
  success: boolean
  data: {
    withdrawals: Withdrawal[]
    total: number
    page: number
    page_size: number
    total_pages: number
  }
}

// ============================================================================
// Seller (Driver MP Account) Types
// ============================================================================

export type SellerStatus =
  | 'not_linked'
  | 'pending'
  | 'active'
  | 'revoked'
  | 'error'

export interface SellerStatusResponse {
  success: boolean
  data: {
    is_linked: boolean
    status: SellerStatus
    mp_user_id?: string
    linked_at?: string
    error_message?: string
  }
}

export interface OAuthURLResponse {
  success: boolean
  oauth_url: string
}

export interface OAuthCallbackRequest {
  code: string
  redirect_uri: string
}

// ============================================================================
// Helper Types
// ============================================================================

export interface PaymentErrorResponse {
  success: false
  error: string
}

export function formatMoney(centavos: number): string {
  const pesos = centavos / 100
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
  }).format(pesos)
}

export function getStatusLabel(status: PaymentStatus): string {
  const labels: Record<PaymentStatus, string> = {
    pending: 'Pendiente',
    in_process: 'En proceso',
    approved: 'Aprobado',
    rejected: 'Rechazado',
    cancelled: 'Cancelado',
    refunded: 'Reembolsado',
    partial_refund: 'Reembolso parcial',
  }
  return labels[status] || status
}

export function getWithdrawalStatusLabel(status: WithdrawalStatus): string {
  const labels: Record<WithdrawalStatus, string> = {
    pending: 'Pendiente',
    processing: 'Procesando',
    completed: 'Completado',
    failed: 'Fallido',
    cancelled: 'Cancelado',
  }
  return labels[status] || status
}

export function getTransactionTypeLabel(type: TransactionType): string {
  const labels: Record<TransactionType, string> = {
    credit: 'Ingreso',
    debit: 'Egreso',
    hold: 'Retenido',
    release: 'Liberado',
    refund: 'Reembolso',
    withdrawal: 'Retiro',
    fee: 'Comision',
  }
  return labels[type] || type
}
