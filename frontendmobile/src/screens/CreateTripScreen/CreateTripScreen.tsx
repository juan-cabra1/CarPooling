/**
 * CreateTripScreen
 * Create and publish a new trip
 */

import React, { useState } from 'react'
import { View, Text, ScrollView, TouchableOpacity, Alert, Platform, Switch } from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import DateTimePicker from '@react-native-community/datetimepicker'
import { Ionicons } from '@expo/vector-icons'
import { Button, Input, Label, Card } from '@/components/ui'
import TextArea from '@/components/ui/TextArea'
import { LocationInput } from '@/components/search'
import { tripsService, getErrorMessage } from '@/services'
import { styles } from './CreateTripScreen.styles'
import type { CreateTripRequest, LocationInput as LocationInputType } from '@/types'
import type { RootStackParamList } from '@/navigation/types'

type CreateTripScreenNavigationProp = StackNavigationProp<RootStackParamList>

export default function CreateTripScreen() {
  const navigation = useNavigation<CreateTripScreenNavigationProp>()
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [origin, setOrigin] = useState<LocationInputType | null>(null)
  const [destination, setDestination] = useState<LocationInputType | null>(null)
  const [showDeparturePicker, setShowDeparturePicker] = useState(false)
  const [showArrivalPicker, setShowArrivalPicker] = useState(false)

  const [formData, setFormData] = useState({
    departure_datetime: new Date(Date.now() + 24 * 60 * 60 * 1000),
    estimated_arrival_datetime: new Date(Date.now() + 30 * 60 * 60 * 1000),
    price_per_seat: '',
    total_seats: '1',
    car: { brand: '', model: '', year: new Date().getFullYear().toString(), color: '', plate: '' },
    preferences: { pets_allowed: false, smoking_allowed: false, music_allowed: false },
    description: '',
  })

  const handleSubmit = async () => {
    try {
      setSaving(true)
      setError('')

      if (!origin?.coordinates || !destination?.coordinates) {
        setError('Debe seleccionar las ubicaciones del menú de autocompletado')
        setSaving(false)
        return
      }

      const now = new Date()
      if (formData.departure_datetime < now) {
        setError('La fecha de salida debe ser en el futuro')
        setSaving(false)
        return
      }

      if (formData.estimated_arrival_datetime <= formData.departure_datetime) {
        setError('La fecha de llegada debe ser posterior a la fecha de salida')
        setSaving(false)
        return
      }

      const createData: CreateTripRequest = {
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

      const newTrip = await tripsService.createTrip(createData)
      Alert.alert('Éxito', 'Viaje publicado exitosamente', [
        { text: 'Ver viaje', onPress: () => navigation.navigate('TripDetail', { id: newTrip.id.toString() }) },
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
    setFormData((prev) => ({ ...prev, [parent]: { ...(prev[parent as keyof typeof prev] as any), [field]: value } }))
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

  return (
    <View style={styles.container}>
      <ScrollView>
        <View style={styles.header}>
          <TouchableOpacity style={styles.backButton} onPress={() => navigation.goBack()}>
            <Ionicons name="arrow-back" size={24} color="#0ea5e9" />
            <Text style={styles.backButtonText}>Volver</Text>
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Publicar Nuevo Viaje</Text>
          <Text style={styles.headerSubtitle}>Completa los detalles de tu viaje</Text>
        </View>

        <View style={styles.content}>
          {error && (
            <View style={styles.errorContainer}>
              <Ionicons name="alert-circle" size={20} color="#dc2626" />
              <Text style={styles.errorText}>{error}</Text>
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
              <Text style={styles.sectionTitle}>Fechas y Horarios</Text>
            </View>

            <View style={styles.inputGroup}>
              <Label>Fecha de salida <Text style={styles.required}>*</Text></Label>
              <TouchableOpacity style={styles.dateButton} onPress={() => setShowDeparturePicker(true)}>
                <Ionicons name="calendar" size={20} color="#6b7280" />
                <Text style={styles.dateButtonText}>{formatDateTime(formData.departure_datetime)}</Text>
              </TouchableOpacity>
            </View>

            {showDeparturePicker && (
              <DateTimePicker
                value={formData.departure_datetime}
                mode="datetime"
                display={Platform.OS === 'ios' ? 'spinner' : 'default'}
                onChange={(event, date) => {
                  setShowDeparturePicker(Platform.OS === 'ios')
                  if (date) handleChange('departure_datetime', date)
                }}
                minimumDate={new Date()}
              />
            )}

            <View style={styles.inputGroup}>
              <Label>Fecha de llegada <Text style={styles.required}>*</Text></Label>
              <TouchableOpacity style={styles.dateButton} onPress={() => setShowArrivalPicker(true)}>
                <Ionicons name="calendar" size={20} color="#6b7280" />
                <Text style={styles.dateButtonText}>{formatDateTime(formData.estimated_arrival_datetime)}</Text>
              </TouchableOpacity>
            </View>

            {showArrivalPicker && (
              <DateTimePicker
                value={formData.estimated_arrival_datetime}
                mode="datetime"
                display={Platform.OS === 'ios' ? 'spinner' : 'default'}
                onChange={(event, date) => {
                  setShowArrivalPicker(Platform.OS === 'ios')
                  if (date) handleChange('estimated_arrival_datetime', date)
                }}
                minimumDate={formData.departure_datetime}
              />
            )}
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="cash" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Precio y Capacidad</Text>
            </View>
            <Input label="Precio por asiento (ARS) *" value={formData.price_per_seat} onChangeText={(v) => handleChange('price_per_seat', v)} keyboardType="numeric" placeholder="5000" />
            <Input label="Asientos disponibles *" value={formData.total_seats} onChangeText={(v) => handleChange('total_seats', v)} keyboardType="numeric" placeholder="1-8" />
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="car" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Vehículo</Text>
            </View>
            <Input label="Marca *" value={formData.car.brand} onChangeText={(v) => handleNestedChange('car', 'brand', v)} placeholder="Toyota" />
            <Input label="Modelo *" value={formData.car.model} onChangeText={(v) => handleNestedChange('car', 'model', v)} placeholder="Corolla" />
            <Input label="Año *" value={formData.car.year} onChangeText={(v) => handleNestedChange('car', 'year', v)} keyboardType="numeric" />
            <Input label="Color *" value={formData.car.color} onChangeText={(v) => handleNestedChange('car', 'color', v)} placeholder="Blanco" />
            <Input label="Patente *" value={formData.car.plate} onChangeText={(v) => handleNestedChange('car', 'plate', v.toUpperCase())} placeholder="ABC123" maxLength={7} autoCapitalize="characters" />
          </Card>

          <Card style={styles.card}>
            <View style={styles.sectionHeader}>
              <Ionicons name="settings" size={20} color="#0ea5e9" />
              <Text style={styles.sectionTitle}>Preferencias</Text>
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🐕 Permitir mascotas</Text>
              <Switch value={formData.preferences.pets_allowed} onValueChange={(v) => handleNestedChange('preferences', 'pets_allowed', v)} />
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🚬 Permitir fumar</Text>
              <Switch value={formData.preferences.smoking_allowed} onValueChange={(v) => handleNestedChange('preferences', 'smoking_allowed', v)} />
            </View>
            <View style={styles.preferenceItem}>
              <Text style={styles.preferenceLabel}>🎵 Permitir música</Text>
              <Switch value={formData.preferences.music_allowed} onValueChange={(v) => handleNestedChange('preferences', 'music_allowed', v)} />
            </View>
          </Card>

          <Card style={styles.card}>
            <TextArea label="Descripción (opcional)" value={formData.description} onChangeText={(v) => handleChange('description', v)} placeholder="Detalles del viaje..." numberOfLines={4} />
          </Card>

          <View style={styles.actionButtons}>
            <Button title="Cancelar" variant="outline" onPress={() => navigation.goBack()} disabled={saving} style={styles.cancelButton} />
            <Button title={saving ? 'Publicando...' : 'Publicar'} onPress={handleSubmit} loading={saving} disabled={saving} style={styles.submitButton} />
          </View>
        </View>
      </ScrollView>
    </View>
  )
}
