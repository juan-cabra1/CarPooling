/**
 * PasswordInput Component
 * Input with toggle visibility for passwords
 */

import React, { useState } from 'react'
import { View, TouchableOpacity } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import Input, { type InputProps } from '../Input/Input'

export interface PasswordInputProps extends Omit<InputProps, 'secureTextEntry'> {}

export default function PasswordInput(props: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false)

  return (
    <View className="relative">
      <Input {...props} secureTextEntry={!showPassword} />
      <TouchableOpacity
        onPress={() => setShowPassword(!showPassword)}
        className="absolute right-4 top-10"
        activeOpacity={0.7}
      >
        <Ionicons
          name={showPassword ? 'eye-off-outline' : 'eye-outline'}
          size={24}
          color="#64748b"
        />
      </TouchableOpacity>
    </View>
  )
}
