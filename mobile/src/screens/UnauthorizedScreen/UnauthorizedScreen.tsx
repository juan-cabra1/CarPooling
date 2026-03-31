import React from 'react'
import { View, Text } from 'react-native'

export default function UnauthorizedScreen() {
  return (
    <View className="flex-1 items-center justify-center bg-background">
      <Text className="text-2xl font-bold">UnauthorizedScreen</Text>
      <Text className="text-muted-foreground mt-2">No tienes permisos</Text>
    </View>
  )
}
