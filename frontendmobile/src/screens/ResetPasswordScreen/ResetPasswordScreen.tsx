/**
 * ResetPasswordScreen - Placeholder
 * TODO: Implement full functionality
 */

import React from 'react'
import { View, Text } from 'react-native'
import { Button } from '@/components/ui'
import { useNavigation } from '@react-navigation/native'

export default function ResetPasswordScreen() {
  const navigation = useNavigation()

  return (
    <View className="flex-1 items-center justify-center bg-background p-6">
      <Text className="text-2xl font-bold mb-4">ResetPasswordScreen</Text>
      <Text className="text-muted-foreground mb-6 text-center">Esta pantalla está pendiente de implementación completa</Text>
      <Button
        title="Volver a Login"
        onPress={() => navigation.navigate('Login' as never)}
      />
    </View>
  )
}
