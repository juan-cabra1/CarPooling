/**
 * BookingModal Component (Mobile)
 * Modal for creating seat reservations + MercadoPago payment
 */

import React, { useState } from 'react'
import { View, Text, Modal, TouchableOpacity, Linking } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { Button, Label } from '@/components/ui'
import { bookingsService, getErrorMessage } from '@/services'
import { createPaymentPreference } from '@/services/paymentsService'
import { useAuth } from '@/context/AuthContext'
import { styles } from './BookingModal.styles'
import type { Trip, SearchTrip } from '@/types'
import type { Booking } from '@/types/booking'

interface BookingModalProps {
  trip: Trip | SearchTrip
  visible: boolean
  onClose: () => void
  onSuccess: () => void
}

type Step = 'booking' | 'payment'

export function BookingModal({ trip, visible, onClose, onSuccess }: BookingModalProps) {
  const { user } = useAuth()
  const [step, setStep] = useState<Step>('booking')
  const [seatsRequested, setSeatsRequested] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [booking, setBooking] = useState<Booking | null>(null)
  const [checkoutURL, setCheckoutURL] = useState('')
  const [prefAmount, setPrefAmount] = useState(0)

  const tripId = 'trip_id' in trip ? trip.trip_id : trip.id
  const totalPrice = seatsRequested * trip.price_per_seat

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 0,
    }).format(price)

  const formatDate = (dateString: string) =>
    new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(dateString))

  // Step 1: create booking + preference
  const handleSubmit = async () => {
    if (!user) { setError('Debés iniciar sesión para reservar'); return }
    if (seatsRequested < 1) { setError('Debés reservar al menos 1 asiento'); return }
    if (seatsRequested > trip.available_seats) {
      setError(`Solo hay ${trip.available_seats} asientos disponibles`)
      return
    }

    try {
      setLoading(true)
      setError('')

      const newBooking = await bookingsService.createBooking({
        trip_id: tripId,
        passenger_id: user.id,
        seats_reserved: seatsRequested,
      })
      setBooking(newBooking)

      const pref = await createPaymentPreference({
        booking_id: newBooking.id,
        trip_id: tripId,
        driver_id: String(trip.driver_id),
        seats_count: seatsRequested,
        price_per_seat: Math.round(trip.price_per_seat * 100), // ARS → centavos
        origin: `${trip.origin.city}, ${trip.origin.province}`,
        destination: `${trip.destination.city}, ${trip.destination.province}`,
        departure_at: trip.departure_datetime,
        payer_email: user.email,
        return_base_url: 'carpooling:/',  // deep link para volver a la app
      })

      setCheckoutURL(pref.init_point)
      setPrefAmount(pref.amount)
      setStep('payment')
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  // Step 2: open MercadoPago in browser
  const handleGoToPayment = async () => {
    if (!checkoutURL) return
    try {
      await Linking.openURL(checkoutURL)
      // After opening browser, close modal and call onSuccess (booking already created)
      onSuccess()
      handleClose()
    } catch {
      setError('No se pudo abrir MercadoPago. Intentá de nuevo.')
    }
  }

  const handleClose = () => {
    if (!loading) {
      setSeatsRequested(1)
      setError('')
      setStep('booking')
      setBooking(null)
      setCheckoutURL('')
      setPrefAmount(0)
      onClose()
    }
  }

  return (
    <Modal visible={visible} animationType="slide" transparent onRequestClose={handleClose}>
      <View style={styles.overlay}>
        <View style={styles.modalContainer}>
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.headerTitle}>
              {step === 'booking' ? 'Reservar Asientos' : 'Confirmar Pago'}
            </Text>
            <TouchableOpacity onPress={handleClose} disabled={loading}>
              <Ionicons name="close" size={28} color="#6b7280" />
            </TouchableOpacity>
          </View>

          <View style={styles.content}>
            {/* ── Step 1: Booking form ── */}
            {step === 'booking' && (
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
                      onPress={() => seatsRequested > 1 && setSeatsRequested(seatsRequested - 1)}
                      disabled={seatsRequested <= 1 || loading}
                    >
                      <Text style={styles.seatsButtonText}>-</Text>
                    </TouchableOpacity>
                    <Text style={styles.seatsValue}>{seatsRequested}</Text>
                    <TouchableOpacity
                      style={styles.seatsButton}
                      onPress={() =>
                        seatsRequested < trip.available_seats &&
                        setSeatsRequested(seatsRequested + 1)
                      }
                      disabled={seatsRequested >= trip.available_seats || loading}
                    >
                      <Text style={styles.seatsButtonText}>+</Text>
                    </TouchableOpacity>
                  </View>
                  <Text style={styles.seatsHint}>Máximo: {trip.available_seats} asientos</Text>
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
                      <Ionicons name="cash" size={20} color="#0ea5e9" />
                      <Text style={styles.priceTotalAmount}>{formatPrice(totalPrice)}</Text>
                    </View>
                  </View>
                </View>

                {error !== '' && (
                  <View style={styles.errorContainer}>
                    <Ionicons name="alert-circle" size={20} color="#dc2626" />
                    <Text style={styles.errorText}>{error}</Text>
                  </View>
                )}

                <View style={styles.actionButtons}>
                  <Button
                    title="Cancelar"
                    variant="outline"
                    onPress={handleClose}
                    disabled={loading}
                    style={styles.cancelButton}
                  />
                  <Button
                    title={loading ? 'Procesando...' : 'Continuar al pago'}
                    onPress={handleSubmit}
                    loading={loading}
                    disabled={loading}
                    style={styles.confirmButton}
                  />
                </View>
              </View>
            )}

            {/* ── Step 2: Payment confirmation ── */}
            {step === 'payment' && (
              <View style={styles.form}>
                <View style={{ alignItems: 'center', paddingVertical: 8, marginBottom: 16 }}>
                  <Ionicons name="checkmark-circle" size={56} color="#22c55e" />
                  <Text style={[styles.successTitle, { marginTop: 8 }]}>¡Reserva creada!</Text>
                  <Text style={[styles.successText, { textAlign: 'center' }]}>
                    Completá el pago para confirmar tu lugar.
                  </Text>
                </View>

                {/* Payment summary */}
                <View style={styles.tripInfoCard}>
                  <Text style={styles.tripRoute}>
                    {trip.origin.city} → {trip.destination.city}
                  </Text>
                  <View style={styles.priceRow}>
                    <Text style={styles.priceLabel}>Total a pagar</Text>
                    <Text style={[styles.priceTotalAmount, { color: '#0ea5e9' }]}>
                      {formatPrice(prefAmount / 100)}
                    </Text>
                  </View>
                </View>

                <Text style={{
                  fontSize: 12,
                  color: '#9ca3af',
                  textAlign: 'center',
                  marginBottom: 16,
                }}>
                  Se abrirá MercadoPago en tu navegador para completar el pago.
                </Text>

                {error !== '' && (
                  <View style={styles.errorContainer}>
                    <Ionicons name="alert-circle" size={20} color="#dc2626" />
                    <Text style={styles.errorText}>{error}</Text>
                  </View>
                )}

                <View style={styles.actionButtons}>
                  <Button
                    title="Pagar después"
                    variant="outline"
                    onPress={() => { onSuccess(); handleClose() }}
                    style={styles.cancelButton}
                  />
                  <Button
                    title="Pagar con MP"
                    onPress={handleGoToPayment}
                    style={styles.confirmButton}
                  />
                </View>
              </View>
            )}
          </View>
        </View>
      </View>
    </Modal>
  )
}
