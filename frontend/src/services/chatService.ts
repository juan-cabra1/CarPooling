/**
 * Chat Service for Trip Messages
 * Handles trip chat/messaging operations
 */

import apiClient from './api'
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

const chatService = {
  sendMessage,
  getMessages,
}

export default chatService
