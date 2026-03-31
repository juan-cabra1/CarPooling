/**
 * EditTripScreen
 * Edit an existing trip
 */

import React, { useState, useEffect } from 'react'
import { View, Text, ScrollView, TouchableOpacity, Alert, Platform, ActivityIndicator } from 'react-native'
import { useNavigation, useRoute } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import DateTimePicker from '@react-native-community/datetimepicker'
import { Ionicons } from '@expo/vector-icons'
import { Button, Input, Label, Card } from '@/components/ui'
import { Switch } from 'react-native'
import TextArea from '@/components/ui/TextArea'
import { LocationInput } from '@/components/search'
import { tripsService, getErrorMessage } from '@/services'
import { styles } from './EditTripScreen.styles'
import type { Trip, LocationInput as LocationInputType } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type EditTripScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function EditTripScreen() {
  const navigation = useNavigation<EditTripScreenNavigationProp>()
  const route = useRoute()
  const { id } = route.params as { id: string }

  const [trip, setTrip] = useState<Trip | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [origin, setOrigin] = useState<LocationInputType | null>(null)
  const [destination, setDestination] = useState<LocationInputType | null>(null)
  const [showDeparturePicker, setShowDeparturePicker] = useState(false)
  const [showArrivalPicker, setShowArrivalPicker] = useState(false)

  const [formData, setFormData] = useState({
    departure_datetime: new Date(),
    estimated_arrival_datetime: new Date(),
    price_per_seat: '',
    total_seats: '1',
    car: { brand: '', model: '', year: new Date().getFullYear().toString(), color: '', plate: '' },
    preferences: { pets_allowed: false, smoking_allowed: false, music_allowed: false },
    description: '',
  })

  useEffect(() => {
    if (id) fetchTrip()
  }, [id])

  const fetchTrip = async () => {
    try {
      setLoading(true)
      const data = await tripsService.getTripById(id)
      setTrip(data)

      setOrigin({
        city: data.origin.city,
        province: data.origin.province,
        address: data.origin.address,
        coordinates: data.origin.coordinates,
      })

      setDestination({
        city: data.destination.city,
        province: data.destination.province,
        address: data.destination.address,
        coordinates: data.destination.coordinates,
      })

      setFormData({
        departure_datetime: new Date(data.departure_datetime),
        estimated_arrival_datetime: new Date(data.estimated_arrival_datetime),
        price_per_seat: data.price_per_seat.toString(),
        total_seats: data.total_seats.toString(),
        car: {
          brand: data.car.brand,
          model: data.car.model,
          year: data.car.year.toString(),
          color: data.car.color,
          plate: data.car.plate,
        },
        preferences: data.preferences,
        description: data.description || '',
      })
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async () => {
    try {
      setSaving(true)
      setError('')

      if (!origin?.coordinates || !destination?.coordinates) {
        setError('Debe seleccionar las ubicaciones del menú de autocompletado')
        setSaving(false)
        return
      }

      const updateData = {
        origin: {
          city: origin.city,
          province: origin.province,
          address: origin.address || `${origin.city}, ${origin.province}`,
          coordinates: origin.coordinates,
        },
        destination: {
          city: destination.city,
          province: destination.province,
          address: destination.address || `${destination.city}, ${destination.province}`,
          coordinates: destination.coordinates,
        },
        departure_datetime: formData.departure_datetime.toISOString(),
        estimated_arrival_datetime: formData.estimated_arrival_datetime.toISOString(),
        price_per_seat: parseFloat(formData.price_per_seat),
        total_seats: parseInt(formData.total_seats),
        car: {
          brand: formData.car.brand,
          model: formData.car.model,
          year: parseInt(formData.car.year),
          color: formData.car.color,
          plate: formData.car.plate.toUpperCase(),
        },
        preferences: formData.preferences,
        description: formData.description,
      }

      await tripsService.updateTrip(id, updateData)
      Alert.alert('Éxito', 'Viaje actualizado exitosamente', [
        { text: 'Ver viaje', onPress: () => navigation.navigate('TripDetail', { id }) },
      ])
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setSaving(false)
    }
  }

  const handleChange = (field: string, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  const handleNestedChange = (parent: string, field: string, value: any) => {
    setFormData((prev) => {
      const parentData = prev[parent as keyof typeof prev]
      if (typeof parentData === 'object' && parentData !== null) {
        return { ...prev, [parent]: { ...parentData, [field]: value } }
      }
      return prev
    })
  }

  const formatDateTime = (date: Date) => {
    return new Intl.DateTimeFormat('es-AR', {
      weekday: 'short',
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando viaje...</Text>
      </View>
    )
  }

  if (error && !trip) {
    return (
      <View style={styles.errorContainer}>
        <Ionicons name="alert-circle-outline" size={64} color="#dc2626" />
        <Text style={styles.errorTitle}>Error al cargar viaje</Text>
        <Text style={styles.errorText}>{error}</Text>
        <Button title="Volver" onPress={() => navigation.goBack()} />
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <ScrollView>
        <View style={styles.header}>
          <TouchableOpacity style={styles.backButton} onPress={() => navigation.goBack()}>
            <Ionicons name="arrow-back" size={24} color="#0ea5e9" />
            <Text style={styles.backButtonText}>Volver</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Editar Viaje</Text>
          <Text style={styles.headerSubtitle}>Actualiza los detalles de tu viaje</Text>
        </View>

        <View style={styles.content}>
          {error && (
            <View style={styles.errorBox}>
              <Ionicons name="alert-circle" size={20} color="#dc2626" />
              <Text style={styles.errorBoxText}>{error}</Text>
            </View>
          )}

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="location" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Ubicaciones</Text>
            </View>
            <LocationInput label="Origen" value={origin} onChange={setOrigin} required />
            <LocationInput label="Destino" value={destination} onChange={setDestination} required />
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="calendar" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Fechas</Text>
            </View>
            <View style={styles.inputGroup}>
              <TouchableOpacity onPress={() => setShowDeparturePicker(true)}>
                <Text style={styles.dateButtonText}>{formatDateTime(formData.departure_datetime)}</Text>
              </TouchableOpacity>
            </View>
            {showDeparturePicker && (
              <DateTimePicker value={formData.departure_datetime} mode="datetime" display={Platform.OS === 'ios' ? 'spinner' : 'default'} onChange={(e, date) => { setShowDeparturePicker(Platform.OS === 'ios'); if (date) handleChange('departure_datetime', date) }} minimumDate={new Date()} />
            )}
            <View style={styles.inputGroup}>
              <TouchableOpacity onPress={() => setShowArrivalPicker(true)}>
                <Text style={styles.dateButtonText}>{formatDateTime(formData.estimated_arrival_datetime)}</Text>
              </TouchableOpacity>
            </View>
            {showArrivalPicker && (
              <DateTimePicker value={formData.estimated_arrival_datetime} mode="datetime" display={Platform.OS === 'ios' ? 'spinner' : 'default'} onChange={(e, date) => { setShowArrivalPicker(Platform.OS === 'ios'); if (date) handleChange('estimated_arrival_datetime', date) }} minimumDate={formData.departure_datetime} />
            )}
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="cash" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Precio y Capacidad</Text>
            </View>
            <Input label="Precio por asiento (ARS) *" value={formData.price_per_seat} onChangeText={(v) => handleChange('price_per_seat', v)} keyboardType="numeric" />
            <Input label="Asientos disponibles *" value={formData.total_seats} onChangeText={(v) => handleChange('total_seats', v)} keyboardType="numeric" />
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="car" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Vehículo</Text>
            </View>
            <Input label="Marca *" value={formData.car.brand} onChangeText={(v) => handleNestedChange('car', 'brand', v)} />
            <Input label="Modelo *" value={formData.car.model} onChangeText={(v) => handleNestedChange('car', 'model', v)} />
            <Input label="Año *" value={formData.car.year} onChangeText={(v) => handleNestedChange('car', 'year', v)} keyboardType="numeric" />
            <Input label="Color *" value={formData.car.color} onChangeText={(v) => handleNestedChange('car', 'color', v)} />
            <Input label="Patente *" value={formData.car.plate} onChangeText={(v) => handleNestedChange('car', 'plate', v.toUpperCase())} maxLength={7} autoCapitalize="characters" />
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="settings" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Preferencias</Text>
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🐕 Permitir mascotas</Text>
              <Switch value={formData.preferences.pets_allowed} onValueChange={(v: boolean) => handleNestedChange('preferences', 'pets_allowed', v)} />
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🚬 Permitir fumar</Text>
              <Switch value={formData.preferences.smoking_allowed} onValueChange={(v: boolean) => handleNestedChange('preferences', 'smoking_allowed', v)} />
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🎵 Permitir música</Text>
              <Switch value={formData.preferences.music_allowed} onValueChange={(v: boolean) => handleNestedChange('preferences', 'music_allowed', v)} />
            </View>
          </Card>

          <Card style={styles.card}>
            <TextArea label="Descripción (opcional)" value={formData.description} onChangeText={(v: string) => handleChange('description', v)} numberOfLines={4} />
          </Card>

          <View style={styles.actionButtons}>
            <Button title="Cancelar" variant="outline" onPress={() => navigation.goBack()} disabled={saving} style={styles.cancelButton} />
            <Button title={saving ? 'Guardando...' : 'Guardar Cambios'} onPress={handleSubmit} loading={saving} disabled={saving} style={styles.submitButton} />
          </View>
        </View>
      </ScrollView>
    </View>
  )
}
