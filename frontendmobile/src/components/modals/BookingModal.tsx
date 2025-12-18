/**
 * BookingModal Component (Mobile)
 * Modal for creating seat reservations with MercadoPago payment
 */

import React, { useState, useEffect } from 'react'
import { View, Text, Modal, TouchableOpacity, Linking, Alert } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { Button, Label } from '@/components/ui'
import { bookingsService, getErrorMessage } from '@/services'
import * as paymentsService from '@/services/paymentsService'
import { useAuth } from '@/context/AuthContext'
import { DevConfig } from '@/config/dev'
import { styles } from './BookingModal.styles'
import type { Trip, SearchTrip } from '@/types'

type BookingStep = 'form' | 'processing' | 'payment' | 'success' | 'error'

interface BookingModalProps {
  trip: Trip | SearchTrip
  visible: boolean
  onClose: () => void
  onSuccess: () => void
}

export function BookingModal({ trip, visible, onClose, onSuccess }: BookingModalProps) {
  const { user } = useAuth()
  const [seatsRequested, setSeatsRequested] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [step, setStep] = useState<BookingStep>('form')
  const [processingMessage, setProcessingMessage] = useState('')
  const [bookingId, setBookingId] = useState<string | null>(null)
  const [pollingInterval, setPollingInterval] = useState<NodeJS.Timeout | null>(null)
  const [isDevMode, setIsDevMode] = useState(false)

  // Handle both Trip and SearchTrip types
  const tripId = 'trip_id' in trip ? trip.trip_id : trip.id
  const driverId = trip.driver_id
  const totalPrice = seatsRequested * trip.price_per_seat

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 0,
    }).format(price)
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }

  // Check dev mode on mount
  useEffect(() => {
    const checkDevMode = async () => {
      const devMode = await DevConfig.isDevModeEnabled()
      setIsDevMode(devMode)
    }
    checkDevMode()
  }, [])

  // Cleanup polling on unmount
  useEffect(() => {
    return () => {
      if (pollingInterval) {
        clearInterval(pollingInterval)
      }
    }
  }, [pollingInterval])

  // Poll payment status
  const pollPaymentStatus = async (bookingIdToCheck: string) => {
    try {
      const payment = await paymentsService.getPaymentByBooking(bookingIdToCheck)

      console.log('Payment status:', payment.status)

      if (payment.status === 'approved') {
        // Payment successful!
        if (pollingInterval) {
          clearInterval(pollingInterval)
          setPollingInterval(null)
        }

        setStep('success')
        setProcessingMessage('Pago confirmado!')

        setTimeout(() => {
          onSuccess()
          handleClose()
        }, 2000)

        return true
      } else if (payment.status === 'rejected' || payment.status === 'cancelled') {
        // Payment failed
        if (pollingInterval) {
          clearInterval(pollingInterval)
          setPollingInterval(null)
        }

        setError('El pago fue rechazado o cancelado. Por favor intenta nuevamente.')
        setStep('error')

        return true
      }

      // Still pending, keep polling
      return false
    } catch (err) {
      console.error('Error polling payment:', err)
      // Continue polling on error (might be temporary network issue)
      return false
    }
  }

  // Start polling after opening MercadoPago
  const startPaymentPolling = (bookingIdToCheck: string) => {
    setProcessingMessage('Esperando confirmación de pago...')

    // Poll every 3 seconds
    const interval = setInterval(async () => {
      await pollPaymentStatus(bookingIdToCheck)
    }, 3000)

    setPollingInterval(interval)

    // Stop polling after 10 minutes (timeout)
    setTimeout(() => {
      if (pollingInterval) {
        clearInterval(interval)
        setPollingInterval(null)
        setError('Tiempo de espera agotado. Verifica el estado de tu reserva en "Mis Viajes".')
        setStep('error')
      }
    }, 600000) // 10 minutes
  }

  const handleSubmit = async () => {
    if (!user) {
      setError('Debes iniciar sesion para reservar')
      return
    }

    if (seatsRequested < 1) {
      setError('Debes reservar al menos 1 asiento')
      return
    }

    if (seatsRequested > trip.available_seats) {
      setError(`Solo hay ${trip.available_seats} asientos disponibles`)
      return
    }

    try {
      setLoading(true)
      setError('')
      setStep('processing')
      setProcessingMessage('Creando reserva...')

      // Step 1: Create booking
      // If dev mode: status will be "confirmed", otherwise "pending_payment"
      const booking = await bookingsService.createBooking({
        trip_id: tripId,
        passenger_id: user.id,
        seats_reserved: seatsRequested,
        dev_mode: isDevMode, // Backend will create as "confirmed" if true
      })

      console.log('Booking created:', booking)

      // DEV MODE: Skip payment flow
      if (isDevMode) {
        console.log('🔧 DEV MODE: Skipping payment flow')
        setStep('success')
        setProcessingMessage('Reserva creada en modo desarrollo!')

        setTimeout(() => {
          onSuccess()
          handleClose()
        }, 2000)

        return
      }

      // PRODUCTION MODE: Continue with payment flow
      setProcessingMessage('Generando link de pago...')

      // Step 2: Create payment preference
      const paymentPreference = await paymentsService.createPaymentPreference({
        booking_id: booking.id,
        trip_id: tripId,
        driver_id: String(driverId),
        seats_count: seatsRequested,
        price_per_seat: Math.round(trip.price_per_seat * 100), // Convert to centavos
        origin: trip.origin.city,
        destination: trip.destination.city,
        departure_at: trip.departure_datetime,
        payer_email: user.email,
        payer_name: `${user.name} ${user.lastname}`,
      })

      console.log('Payment preference created:', paymentPreference)
      setBookingId(booking.id)

      // Step 3: Open MercadoPago checkout
      // Use sandbox_init_point for test mode, init_point for production
      const checkoutUrl = paymentPreference.sandbox_init_point || paymentPreference.init_point
      const canOpen = await Linking.canOpenURL(checkoutUrl)

      if (canOpen) {
        setStep('payment')
        await Linking.openURL(checkoutUrl)

        // Start polling for payment confirmation
        setLoading(false) // Allow user to interact while polling
        startPaymentPolling(booking.id)

        // Show info message
        Alert.alert(
          'Completa tu pago',
          'Se abrio MercadoPago. Completa el pago y volveremos a verificar el estado automaticamente.',
          [{ text: 'OK' }]
        )
      } else {
        throw new Error('No se pudo abrir el link de pago')
      }
    } catch (err) {
      console.error('Error in booking/payment flow:', err)
      setError(getErrorMessage(err))
      setStep('error')
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    if (!loading) {
      // Clear polling if active
      if (pollingInterval) {
        clearInterval(pollingInterval)
        setPollingInterval(null)
      }

      setSeatsRequested(1)
      setError('')
      setStep('form')
      setProcessingMessage('')
      setBookingId(null)
      onClose()
    }
  }

  const handleRetry = () => {
    setError('')
    setStep('form')
  }

  const incrementSeats = () => {
    if (seatsRequested < trip.available_seats) {
      setSeatsRequested(seatsRequested + 1)
    }
  }

  const decrementSeats = () => {
    if (seatsRequested > 1) {
      setSeatsRequested(seatsRequested - 1)
    }
  }

  const renderContent = () => {
    switch (step) {
      case 'processing':
        return (
          <View style={styles.successContainer}>
            <Ionicons name="hourglass" size={64} color="#0ea5e9" />
            <Text style={styles.successTitle}>{processingMessage}</Text>
            <Text style={styles.successText}>Por favor espera...</Text>
          </View>
        )

      case 'payment':
        return (
          <View style={styles.successContainer}>
            <Ionicons name="hourglass" size={64} color="#0ea5e9" />
            <Text style={styles.successTitle}>{processingMessage}</Text>
            <Text style={styles.successText}>
              Completa el pago en MercadoPago. Verificaremos el estado automaticamente.
            </Text>
            <Button
              title="Verificar ahora"
              variant="outline"
              onPress={() => bookingId && pollPaymentStatus(bookingId)}
              style={{ marginTop: 16 }}
            />
            <Button
              title="Cancelar"
              variant="outline"
              onPress={handleClose}
              style={{ marginTop: 8 }}
            />
          </View>
        )

      case 'success':
        return (
          <View style={styles.successContainer}>
            <Ionicons name="checkmark-circle" size={64} color="#22c55e" />
            <Text style={styles.successTitle}>Reserva iniciada!</Text>
            <Text style={styles.successText}>
              Completa el pago en MercadoPago para confirmar tu reserva.
            </Text>
          </View>
        )

      case 'error':
        return (
          <View style={styles.form}>
            <View style={styles.errorContainer}>
              <Ionicons name="alert-circle" size={20} color="#dc2626" />
              <Text style={styles.errorText}>{error}</Text>
            </View>
            <View style={styles.actionButtons}>
              <Button
                title="Cancelar"
                variant="outline"
                onPress={handleClose}
                style={styles.cancelButton}
              />
              <Button title="Reintentar" onPress={handleRetry} style={styles.confirmButton} />
            </View>
          </View>
        )

      default:
        return (
          <View style={styles.form}>
            {/* Trip Info */}
            <View style={styles.tripInfoCard}>
              <Text style={styles.tripRoute}>
                {trip.origin.city} → {trip.destination.city}
              </Text>
              <Text style={styles.tripDate}>{formatDate(trip.departure_datetime)}</Text>
              <View style={styles.tripSeatsInfo}>
                <Ionicons name="people" size={16} color="#6b7280" />
                <Text style={styles.tripSeatsText}>
                  {trip.available_seats} asientos disponibles
                </Text>
              </View>
            </View>

            {/* Seats Selection */}
            <View style={styles.seatsSection}>
              <Label style={styles.seatsLabel}>Cantidad de asientos</Label>
              <View style={styles.seatsControl}>
                <TouchableOpacity
                  style={styles.seatsButton}
                  onPress={decrementSeats}
                  disabled={seatsRequested <= 1 || loading}
                >
                  <Text style={styles.seatsButtonText}>-</Text>
                </TouchableOpacity>
                <Text style={styles.seatsValue}>{seatsRequested}</Text>
                <TouchableOpacity
                  style={styles.seatsButton}
                  onPress={incrementSeats}
                  disabled={seatsRequested >= trip.available_seats || loading}
                >
                  <Text style={styles.seatsButtonText}>+</Text>
                </TouchableOpacity>
              </View>
              <Text style={styles.seatsHint}>Maximo: {trip.available_seats} asientos</Text>
            </View>

            {/* Price Summary */}
            <View style={styles.priceSummary}>
              <View style={styles.priceRow}>
                <Text style={styles.priceLabel}>Precio por asiento</Text>
                <Text style={styles.priceValue}>{formatPrice(trip.price_per_seat)}</Text>
              </View>
              <View style={styles.priceRow}>
                <Text style={styles.priceLabel}>Asientos</Text>
                <Text style={styles.priceValue}>{seatsRequested}</Text>
              </View>
              <View style={styles.priceDivider} />
              <View style={styles.priceTotalRow}>
                <Text style={styles.priceTotalLabel}>Total</Text>
                <View style={styles.priceTotalValue}>
                  <Ionicons name="card" size={20} color="#0ea5e9" />
                  <Text style={styles.priceTotalAmount}>{formatPrice(totalPrice)}</Text>
                </View>
              </View>
            </View>

            {/* Dev Mode Warning */}
            {isDevMode && (
              <View style={styles.devModeCard}>
                <Ionicons name="construct" size={20} color="#f59e0b" />
                <Text style={styles.devModeText}>
                  MODO DESARROLLO: La reserva se creará sin procesar pagos
                </Text>
              </View>
            )}

            {/* Payment Info */}
            {!isDevMode && (
              <View style={styles.paymentInfoCard}>
                <Ionicons name="shield-checkmark" size={20} color="#22c55e" />
                <Text style={styles.paymentInfoText}>
                  Pago seguro con MercadoPago. El dinero se retiene hasta que confirmes el viaje.
                </Text>
              </View>
            )}

            {/* Error Message */}
            {error && (
              <View style={styles.errorContainer}>
                <Ionicons name="alert-circle" size={20} color="#dc2626" />
                <Text style={styles.errorText}>{error}</Text>
              </View>
            )}

            {/* Action Buttons */}
            <View style={styles.actionButtons}>
              <Button
                title="Cancelar"
                variant="outline"
                onPress={handleClose}
                disabled={loading}
                style={styles.cancelButton}
              />
              <Button
                title={loading ? 'Procesando...' : isDevMode ? 'Reservar (Sin Pago)' : 'Pagar con MercadoPago'}
                onPress={handleSubmit}
                loading={loading}
                disabled={loading}
                style={styles.confirmButton}
              />
            </View>
          </View>
        )
    }
  }

  return (
    <Modal visible={visible} animationType="slide" transparent onRequestClose={handleClose}>
      <View style={styles.overlay}>
        <View style={styles.modalContainer}>
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.headerTitle}>Reservar Asientos</Text>
            <TouchableOpacity onPress={handleClose} disabled={loading}>
              <Ionicons name="close" size={28} color="#6b7280" />
            </TouchableOpacity>
          </View>

          {/* Content */}
          <View style={styles.content}>{renderContent()}</View>
        </View>
      </View>
    </Modal>
  )
}
