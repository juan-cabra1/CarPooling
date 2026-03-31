/**
 * EditProfileScreen
 * Allows users to edit their personal information
 */

import React, { useState } from 'react'
import {
  View,
  Text,
  ScrollView,
  TextInput,
  TouchableOpacity,
  ActivityIndicator,
  StyleSheet,
  Alert,
} from 'react-native'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { useAuth } from '@/context/AuthContext'
import { authService } from '@/services'
import type { RootStackParamList } from '@/navigation/types'
import type { UpdateUserData } from '@/types'

type EditProfileNavigationProp = StackNavigationProp<RootStackParamList>

export default function EditProfileScreen() {
  const navigation = useNavigation<EditProfileNavigationProp>()
  const { user, refreshUser } = useAuth()

  const [formData, setFormData] = useState<UpdateUserData>({
    name: user?.name ?? '',
    lastname: user?.lastname ?? '',
    phone: user?.phone ?? '',
    street: user?.street ?? '',
    number: user?.number ?? 0,
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSave = async () => {
    if (!user) return
    setError('')
    setLoading(true)
    try {
      await authService.updateProfile(user.id, formData)
      await refreshUser()
      Alert.alert('Éxito', 'Perfil actualizado correctamente', [
        { text: 'OK', onPress: () => navigation.goBack() },
      ])
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al guardar cambios'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <Text style={styles.title}>Editar Perfil</Text>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Nombre</Text>
        <TextInput
          value={formData.name ?? ''}
          onChangeText={v => setFormData(d => ({ ...d, name: v }))}
          style={styles.input}
          placeholder="Tu nombre"
          placeholderTextColor="#9ca3af"
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Apellido</Text>
        <TextInput
          value={formData.lastname ?? ''}
          onChangeText={v => setFormData(d => ({ ...d, lastname: v }))}
          style={styles.input}
          placeholder="Tu apellido"
          placeholderTextColor="#9ca3af"
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Teléfono</Text>
        <TextInput
          value={formData.phone ?? ''}
          onChangeText={v => setFormData(d => ({ ...d, phone: v }))}
          style={styles.input}
          placeholder="Tu teléfono"
          placeholderTextColor="#9ca3af"
          keyboardType="phone-pad"
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Calle</Text>
        <TextInput
          value={formData.street ?? ''}
          onChangeText={v => setFormData(d => ({ ...d, street: v }))}
          style={styles.input}
          placeholder="Nombre de la calle"
          placeholderTextColor="#9ca3af"
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Número</Text>
        <TextInput
          value={formData.number ? formData.number.toString() : ''}
          onChangeText={v => setFormData(d => ({ ...d, number: parseInt(v) || 0 }))}
          style={styles.input}
          placeholder="Número de puerta"
          placeholderTextColor="#9ca3af"
          keyboardType="numeric"
        />
      </View>

      {error !== '' && (
        <View style={styles.errorBox}>
          <Text style={styles.errorText}>{error}</Text>
        </View>
      )}

      <TouchableOpacity
        style={[styles.saveButton, loading && styles.saveButtonDisabled]}
        onPress={handleSave}
        disabled={loading}
      >
        {loading ? (
          <ActivityIndicator color="#ffffff" />
        ) : (
          <Text style={styles.saveButtonText}>Guardar cambios</Text>
        )}
      </TouchableOpacity>

      <TouchableOpacity
        style={styles.cancelButton}
        onPress={() => navigation.goBack()}
        disabled={loading}
      >
        <Text style={styles.cancelButtonText}>Cancelar</Text>
      </TouchableOpacity>
    </ScrollView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#ffffff',
  },
  content: {
    padding: 24,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#111827',
    marginBottom: 24,
  },
  fieldGroup: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '500',
    color: '#374151',
    marginBottom: 6,
  },
  input: {
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 10,
    fontSize: 14,
    color: '#111827',
    backgroundColor: '#f9fafb',
  },
  errorBox: {
    backgroundColor: '#fef2f2',
    borderWidth: 1,
    borderColor: '#fecaca',
    borderRadius: 8,
    padding: 12,
    marginBottom: 16,
  },
  errorText: {
    color: '#dc2626',
    fontSize: 14,
  },
  saveButton: {
    backgroundColor: '#0ea5e9',
    borderRadius: 8,
    paddingVertical: 14,
    alignItems: 'center',
    marginBottom: 12,
  },
  saveButtonDisabled: {
    opacity: 0.6,
  },
  saveButtonText: {
    color: '#ffffff',
    fontWeight: '600',
    fontSize: 16,
  },
  cancelButton: {
    borderWidth: 1,
    borderColor: '#d1d5db',
    borderRadius: 8,
    paddingVertical: 14,
    alignItems: 'center',
  },
  cancelButtonText: {
    color: '#374151',
    fontWeight: '500',
    fontSize: 16,
  },
})
