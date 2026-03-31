/**
 * Payment-related types for the CarPooling Payments API (Port 8005)
 * Database: MySQL
 * Integration: MercadoPago
 */

/**
 * Request to create a MercadoPago payment preference
 */
export interface CreatePreferenceRequest {
  booking_id: string
  trip_id: string
  driver_id: string
  seats_count: number
  price_per_seat: number   // In centavos (ARS cents)
  origin: string           // "City, Province"
  destination: string      // "City, Province"
  departure_at: string     // ISO 8601
  payer_email?: string
  payer_name?: string
  return_base_url?: string // Web: window.location.origin, Mobile: "carpooling:/"
}

/**
 * Response after creating a payment preference
 */
export interface CreatePreferenceResponse {
  payment_uuid: string
  preference_id: string
  init_point: string        // Production checkout URL (MercadoPago)
  sandbox_init_point: string
  amount: number            // Total in centavos
  platform_fee: number      // Platform fee in centavos
  driver_amount: number     // Amount driver will receive in centavos
  expires_at: string        // ISO 8601
}

/**
 * Payment details
 */
export interface PaymentResponse {
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
  status: 'pending' | 'approved' | 'rejected' | 'refunded' | 'released'
  fund_status: string
  mp_payment_id?: string
  approved_at?: string
  released_at?: string
  created_at: string
}

/**
 * Wallet details response
 */
export interface WalletResponse {
  wallet_uuid: string
  user_id: string
  available_balance: number    // In centavos
  pending_balance: number
  held_balance: number
  total_balance: number
  currency: string
  total_earned: number
  total_withdrawn: number
  total_trips: number
  auto_withdraw_enabled: boolean
  auto_withdraw_threshold: number
  // Pre-formatted strings from backend
  available_balance_formatted: string
  total_earned_formatted: string
}

/**
 * Wallet settings update
 */
export interface WalletSettingsRequest {
  auto_withdraw_enabled?: boolean
  auto_withdraw_threshold?: number  // In centavos
}

/**
 * Single wallet transaction
 */
export interface TransactionResponse {
  transaction_uuid: string
  type: 'credit' | 'debit' | 'withdrawal' | 'refund'
  amount: number         // In centavos
  currency: string
  balance_before: number
  balance_after: number
  description: string
  payment_id?: string
  withdrawal_id?: string
  created_at: string
  amount_formatted: string
  signed_amount: number  // Positive=credit, negative=debit
}

/**
 * Paginated transactions list
 */
export interface TransactionsListResponse {
  transactions: TransactionResponse[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

/**
 * Wallet statistics
 */
export interface WalletStatsResponse {
  current_month_earnings: number
  current_month_trips: number
  last_month_earnings: number
  last_month_trips: number
  earnings_growth_percent: number
  total_earnings: number
  total_trips: number
  total_withdrawals: number
  average_per_trip: number
}

/**
 * Withdrawal request
 */
export interface WithdrawalRequest {
  amount: number  // In centavos
}

/**
 * Withdrawal details response
 */
export interface WithdrawalResponse {
  withdrawal_uuid: string
  amount: number
  fee: number
  net_amount: number
  currency: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  destination_account?: string
  error_message?: string
  processed_at?: string
  completed_at?: string
  created_at: string
  amount_formatted: string
}

/**
 * Paginated withdrawals list
 */
export interface WithdrawalsListResponse {
  withdrawals: WithdrawalResponse[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

/**
 * Seller (driver MP account) status
 */
export interface SellerStatusResponse {
  linked: boolean
  mp_user_id?: string
  mp_email?: string
  created_at?: string
}
