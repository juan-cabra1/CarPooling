/**
 * Chat Service for Trip Messages
 * Handles trip chat/messaging operations
 */

import { tripsApi as apiClient } from './api'
import type { ApiResponse } from '@/types'

/**
 * Message type definition
 */
export interface Message {
  id: string
  trip_id: string
  user_id: number
  user_name: string
  message: string
  created_at: string
}

/**
 * Send message request
 */
export interface SendMessageRequest {
  message: string
}

/**
 * Messages response from API
 */
export interface MessagesResponse {
  messages: Message[]
  count: number
}

/**
 * Send a chat message to a trip
 * Requires authentication
 * @param tripId - Trip ID
 * @param message - Message text
 * @returns Created message
 * @example
 * const message = await chatService.sendMessage('507f1f77bcf86cd799439011', 'Hello!')
 */
export async function sendMessage(tripId: string, message: string): Promise<Message> {
  const response = await apiClient.post<ApiResponse<Message>>(
    `/trips/${tripId}/messages`,
    { message }
  )
  return response.data.data!
}

/**
 * Get all messages for a trip
 * Requires authentication
 * @param tripId - Trip ID
 * @returns List of messages for the trip
 * @example
 * const { messages, count } = await chatService.getMessages('507f1f77bcf86cd799439011')
 */
export async function getMessages(tripId: string): Promise<MessagesResponse> {
  const response = await apiClient.get<ApiResponse<MessagesResponse>>(
    `/trips/${tripId}/messages`
  )

  // API returns { success: true, messages: [...], count: N }
  // We need to extract the data properly
  const data = response.data

  if (data.success) {
    return {
      messages: (data as any).messages || [],
      count: (data as any).count || 0,
    }
  }

  return { messages: [], count: 0 }
}

/**
 * Build the WebSocket URL for a trip's chat room.
 * Uses VITE_TRIPS_WS_URL env var in production, derives from window.location in dev.
 */
export function getChatWebSocketURL(tripId: string, token: string): string {
  const base = import.meta.env.VITE_TRIPS_WS_URL
  if (base) {
    return `${base}/${tripId}/ws?token=${encodeURIComponent(token)}`
  }
  // Fallback for local dev: derive from current host
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // Vite dev proxy is on the same host but trips-api is on port 8002
  return `${protocol}//localhost:8002/trips/${tripId}/ws?token=${encodeURIComponent(token)}`
}

const chatService = {
  sendMessage,
  getMessages,
  getChatWebSocketURL,
}

export default chatService
