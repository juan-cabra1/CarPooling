/**
 * Bottom Tab Navigator
 * Main navigation for authenticated users
 */

import React from 'react'
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs'
import { Ionicons } from '@expo/vector-icons'
import type { TabParamList } from './types'

// Placeholder screens (will be replaced with actual screens)
import HomeScreen from '@/screens/HomeScreen/HomeScreen'
import SearchScreen from '@/screens/SearchScreen/SearchScreen'
import MyTripsScreen from '@/screens/MyTripsScreen/MyTripsScreen'
import MyBookingsScreen from '@/screens/MyBookingsScreen/MyBookingsScreen'
import ProfileScreen from '@/screens/ProfileScreen/ProfileScreen'

const Tab = createBottomTabNavigator<TabParamList>()

export default function TabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={({ route }) => ({
        tabBarIcon: ({ focused, color, size }) => {
          let iconName: keyof typeof Ionicons.glyphMap = 'home'

          if (route.name === 'Home') {
            iconName = focused ? 'home' : 'home-outline'
          } else if (route.name === 'Search') {
            iconName = focused ? 'search' : 'search-outline'
          } else if (route.name === 'MyTrips') {
            iconName = focused ? 'car' : 'car-outline'
          } else if (route.name === 'MyBookings') {
            iconName = focused ? 'list' : 'list-outline'
          } else if (route.name === 'Profile') {
            iconName = focused ? 'person' : 'person-outline'
          }

          return <Ionicons name={iconName} size={size} color={color} />
        },
        tabBarActiveTintColor: '#0ea5e9', // primary-500
        tabBarInactiveTintColor: '#64748b', // muted-foreground
        tabBarStyle: {
          backgroundColor: '#ffffff',
          borderTopColor: '#cbd5e1',
          paddingBottom: 5,
          paddingTop: 5,
          height: 60,
        },
        headerStyle: {
          backgroundColor: '#0ea5e9', // primary-500
        },
        headerTintColor: '#ffffff',
        headerTitleStyle: {
          fontWeight: 'bold',
        },
      })}
    >
      <Tab.Screen
        name="Home"
        component={HomeScreen}
        options={{ title: 'Inicio' }}
      />
      <Tab.Screen
        name="Search"
        component={SearchScreen}
        options={{ title: 'Buscar' }}
      />
      <Tab.Screen
        name="MyTrips"
        component={MyTripsScreen}
        options={{ title: 'Mis Viajes' }}
      />
      <Tab.Screen
        name="MyBookings"
        component={MyBookingsScreen}
        options={{ title: 'Mis Reservas' }}
      />
      <Tab.Screen
        name="Profile"
        component={ProfileScreen}
        options={{ title: 'Perfil' }}
      />
    </Tab.Navigator>
  )
}
