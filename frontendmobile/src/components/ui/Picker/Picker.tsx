/**
 * Picker Component
 * Simple picker/select for React Native
 */

import React, { useState } from 'react'
import { View, Text, TouchableOpacity, Modal, FlatList } from 'react-native'
import { Ionicons } from '@expo/vector-icons'

export interface PickerOption {
  label: string
  value: string
}

export interface PickerProps {
  label?: string
  value: string
  onValueChange: (value: string) => void
  options: PickerOption[]
  placeholder?: string
  error?: string
  disabled?: boolean
}

export default function Picker({
  label,
  value,
  onValueChange,
  options,
  placeholder = 'Seleccionar...',
  error,
  disabled = false,
}: PickerProps) {
  const [modalVisible, setModalVisible] = useState(false)

  const selectedOption = options.find((opt) => opt.value === value)

  return (
    <View className="mb-4">
      {label && <Text className="text-sm font-medium text-foreground mb-2">{label}</Text>}

      <TouchableOpacity
        onPress={() => !disabled && setModalVisible(true)}
        disabled={disabled}
        className={`bg-white border rounded-lg px-4 py-3 flex-row justify-between items-center ${
          error ? 'border-destructive border-2' : 'border-input'
        } ${disabled ? 'bg-muted' : ''}`}
      >
        <Text className={selectedOption ? 'text-foreground' : 'text-muted-foreground'}>
          {selectedOption ? selectedOption.label : placeholder}
        </Text>
        <Ionicons name="chevron-down" size={20} color="#64748b" />
      </TouchableOpacity>

      {error && <Text className="text-destructive text-sm mt-1">{error}</Text>}

      <Modal visible={modalVisible} transparent animationType="slide">
        <View className="flex-1 justify-end bg-black/50">
          <View className="bg-white rounded-t-3xl max-h-96">
            <View className="flex-row justify-between items-center p-4 border-b border-border">
              <Text className="text-lg font-bold">{label || 'Seleccionar'}</Text>
              <TouchableOpacity onPress={() => setModalVisible(false)}>
                <Ionicons name="close" size={24} color="#1e293b" />
              </TouchableOpacity>
            </View>
            <FlatList
              data={options}
              keyExtractor={(item) => item.value}
              renderItem={({ item }) => (
                <TouchableOpacity
                  onPress={() => {
                    onValueChange(item.value)
                    setModalVisible(false)
                  }}
                  className="p-4 border-b border-border flex-row justify-between items-center"
                >
                  <Text className="text-base">{item.label}</Text>
                  {value === item.value && (
                    <Ionicons name="checkmark" size={24} color="#0ea5e9" />
                  )}
                </TouchableOpacity>
              )}
            />
          </View>
        </View>
      </Modal>
    </View>
  )
}
