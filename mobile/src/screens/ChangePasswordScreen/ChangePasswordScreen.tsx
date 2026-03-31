/**
 * ChangePasswordScreen
 * Allows users to change their password
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
import { authService } from '@/services'
import type { RootStackParamList } from '@/navigation/types'

type ChangePasswordNavigationProp = StackNavigationProp<RootStackParamList>

export default function ChangePasswordScreen() {
  const navigation = useNavigation<ChangePasswordNavigationProp>()

  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSave = async () => {
    setError('')

    if (newPassword !== confirmPassword) {
      setError('Las contraseñas nuevas no coinciden')
      return
    }
    if (newPassword.length < 8) {
      setError('La nueva contraseña debe tener al menos 8 caracteres')
      return
    }

    setLoading(true)
    try {
      await authService.changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      })
      Alert.alert('Éxito', 'Contraseña cambiada correctamente', [
        { text: 'OK', onPress: () => navigation.goBack() },
      ])
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al cambiar contraseña'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      <Text style={styles.title}>Cambiar Contraseña</Text>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Contraseña actual</Text>
        <TextInput
          value={currentPassword}
          onChangeText={setCurrentPassword}
          style={styles.input}
          placeholder="Ingresa tu contraseña actual"
          placeholderTextColor="#9ca3af"
          secureTextEntry
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Nueva contraseña</Text>
        <TextInput
          value={newPassword}
          onChangeText={setNewPassword}
          style={styles.input}
          placeholder="Mínimo 8 caracteres"
          placeholderTextColor="#9ca3af"
          secureTextEntry
        />
      </View>

      <View style={styles.fieldGroup}>
        <Text style={styles.label}>Confirmar nueva contraseña</Text>
        <TextInput
          value={confirmPassword}
          onChangeText={setConfirmPassword}
          style={styles.input}
          placeholder="Repite la nueva contraseña"
          placeholderTextColor="#9ca3af"
          secureTextEntry
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
          <Text style={styles.saveButtonText}>Cambiar contraseña</Text>
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
