/**
 * Development Mode Configuration
 *
 * Permite hacer testing sin implementación completa de pagos.
 * En modo dev, las reservas se crean directamente sin pasar por MercadoPago.
 */

import AsyncStorage from '@react-native-async-storage/async-storage'

const DEV_MODE_STORAGE_KEY = '@carpooling:dev_mode'

/**
 * Configuración de modo desarrollo
 */
export const DevConfig = {
  /**
   * Verifica si el modo desarrollo está activado
   */
  async isDevModeEnabled(): Promise<boolean> {
    try {
      const value = await AsyncStorage.getItem(DEV_MODE_STORAGE_KEY)
      return value === 'true'
    } catch (error) {
      console.error('Error reading dev mode:', error)
      return false
    }
  },

  /**
   * Activa el modo desarrollo
   */
  async enableDevMode(): Promise<void> {
    try {
      await AsyncStorage.setItem(DEV_MODE_STORAGE_KEY, 'true')
      console.log('✅ Dev mode enabled - payments will be skipped')
    } catch (error) {
      console.error('Error enabling dev mode:', error)
      throw error
    }
  },

  /**
   * Desactiva el modo desarrollo
   */
  async disableDevMode(): Promise<void> {
    try {
      await AsyncStorage.setItem(DEV_MODE_STORAGE_KEY, 'false')
      console.log('✅ Dev mode disabled - payments will be processed normally')
    } catch (error) {
      console.error('Error disabling dev mode:', error)
      throw error
    }
  },

  /**
   * Toggle dev mode
   */
  async toggleDevMode(): Promise<boolean> {
    const isEnabled = await this.isDevModeEnabled()
    if (isEnabled) {
      await this.disableDevMode()
      return false
    } else {
      await this.enableDevMode()
      return true
    }
  },
}

export default DevConfig
