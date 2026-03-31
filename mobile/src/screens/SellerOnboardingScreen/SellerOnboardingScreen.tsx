/**
 * SellerOnboardingScreen
 * Guide drivers to link their MercadoPago account for receiving payments
 */

import React, { useState, useEffect } from 'react'
import { View, Text, ScrollView, Alert, ActivityIndicator, TextInput, TouchableOpacity } from 'react-native'
import * as WebBrowser from 'expo-web-browser'
import { useNavigation } from '@react-navigation/native'
import { type StackNavigationProp } from '@react-navigation/stack'
import { Ionicons } from '@expo/vector-icons'
import { Button, Card } from '@/components/ui'
import * as paymentsService from '@/services/paymentsService'
import { getErrorMessage } from '@/types'
import type { SellerStatusResponse } from '@/types'
import { styles } from './SellerOnboardingScreen.styles'
import type { RootStackParamList } from '@/navigation/types'

type SellerOnboardingNavigationProp = StackNavigationProp<RootStackParamList>

// Backend HTTPS redirect URI — MP requires HTTPS, backend bounces back to app deep link
const REDIRECT_URI = 'https://payments-api-production-e0e9.up.railway.app/api/v1/seller/oauth/web-callback'

export default function SellerOnboardingScreen() {
  const navigation = useNavigation<SellerOnboardingNavigationProp>()

  const [sellerStatus, setSellerStatus] = useState<SellerStatusResponse['data'] | null>(null)
  const [loading, setLoading] = useState(true)
  const [linking, setLinking] = useState(false)
  const [error, setError] = useState('')
  const [showManualInput, setShowManualInput] = useState(false)
  const [accessToken, setAccessToken] = useState('')

  useEffect(() => {
    loadSellerStatus()
  }, [])

  const loadSellerStatus = async () => {
    try {
      const status = await paymentsService.getSellerStatus()
      setSellerStatus(status)

      if (status.is_linked) {
        Alert.alert('Cuenta vinculada', 'Tu cuenta de MercadoPago ya esta vinculada.', [
          { text: 'OK', onPress: () => navigation.goBack() },
        ])
      }
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleLinkAccount = async () => {
    try {
      setLinking(true)
      setError('')

      const oauthUrl = await paymentsService.getSellerOAuthURL(REDIRECT_URI)

      // openAuthSessionAsync usa ASWebAuthenticationSession en iOS — maneja mejor
      // los formularios de login y captura automáticamente el redirect al deep link
      const result = await WebBrowser.openAuthSessionAsync(oauthUrl, 'carpooling://seller/callback')

      if (result.type === 'success' && result.url) {
        const urlParams = new URL(result.url.replace('carpooling://', 'https://'))
        const code = urlParams.searchParams.get('code')

        if (code) {
          await paymentsService.completeSellerOAuth({ code, redirect_uri: REDIRECT_URI })
          Alert.alert(
            'Cuenta vinculada',
            'Tu cuenta de MercadoPago fue vinculada exitosamente. Ya puedes recibir pagos y retirar fondos.',
            [{ text: 'Continuar', onPress: () => navigation.goBack() }]
          )
        } else {
          Alert.alert('Error', 'No se recibió el código de autorización')
        }
      }
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLinking(false)
    }
  }

  const handleLinkManually = async () => {
    if (!accessToken.trim()) {
      setError('Por favor ingresa tu Access Token de MercadoPago')
      return
    }

    try {
      setLinking(true)
      setError('')

      await paymentsService.linkSellerManually(accessToken.trim())

      Alert.alert(
        'Cuenta vinculada',
        'Tu cuenta de MercadoPago fue vinculada exitosamente. Ya puedes recibir pagos y retirar fondos.',
        [{ text: 'Continuar', onPress: () => navigation.goBack() }]
      )
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLinking(false)
    }
  }

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text style={styles.loadingText}>Cargando...</Text>
      </View>
    )
  }

  return (
    <View style={styles.container}>
      <ScrollView>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.iconContainer}>
            <Ionicons name="card" size={48} color="#0ea5e9" />
          </View>
          <Text style={styles.title}>Vincula tu cuenta de MercadoPago</Text>
          <Text style={styles.subtitle}>
            Para recibir pagos de tus pasajeros y retirar tus ganancias, necesitas vincular tu
            cuenta de MercadoPago.
          </Text>
        </View>

        {/* Benefits */}
        <Card style={styles.benefitsCard}>
          <Text style={styles.cardTitle}>Beneficios</Text>

          <View style={styles.benefitItem}>
            <View style={styles.benefitIcon}>
              <Ionicons name="flash" size={20} color="#22c55e" />
            </View>
            <View style={styles.benefitText}>
              <Text style={styles.benefitTitle}>Pagos automaticos</Text>
              <Text style={styles.benefitDescription}>
                Recibe el dinero de tus viajes automaticamente en tu billetera
              </Text>
            </View>
          </View>

          <View style={styles.benefitItem}>
            <View style={styles.benefitIcon}>
              <Ionicons name="shield-checkmark" size={20} color="#22c55e" />
            </View>
            <View style={styles.benefitText}>
              <Text style={styles.benefitTitle}>Seguridad garantizada</Text>
              <Text style={styles.benefitDescription}>
                El dinero se retiene hasta que el pasajero confirme el viaje
              </Text>
            </View>
          </View>

          <View style={styles.benefitItem}>
            <View style={styles.benefitIcon}>
              <Ionicons name="wallet" size={20} color="#22c55e" />
            </View>
            <View style={styles.benefitText}>
              <Text style={styles.benefitTitle}>Retiros faciles</Text>
              <Text style={styles.benefitDescription}>
                Retira tu dinero a tu cuenta de MercadoPago cuando quieras
              </Text>
            </View>
          </View>
        </Card>

        {/* How it works */}
        <Card style={styles.stepsCard}>
          <Text style={styles.cardTitle}>Como funciona</Text>

          <View style={styles.stepItem}>
            <View style={styles.stepNumber}>
              <Text style={styles.stepNumberText}>1</Text>
            </View>
            <Text style={styles.stepText}>Toca "Vincular cuenta" abajo</Text>
          </View>

          <View style={styles.stepItem}>
            <View style={styles.stepNumber}>
              <Text style={styles.stepNumberText}>2</Text>
            </View>
            <Text style={styles.stepText}>Inicia sesion en MercadoPago</Text>
          </View>

          <View style={styles.stepItem}>
            <View style={styles.stepNumber}>
              <Text style={styles.stepNumberText}>3</Text>
            </View>
            <Text style={styles.stepText}>Autoriza a Carpooling a recibir pagos</Text>
          </View>

          <View style={styles.stepItem}>
            <View style={styles.stepNumber}>
              <Text style={styles.stepNumberText}>4</Text>
            </View>
            <Text style={styles.stepText}>Listo! Ya puedes recibir pagos</Text>
          </View>
        </Card>

        {/* Info */}
        <Card style={styles.infoCard}>
          <Ionicons name="information-circle" size={24} color="#0369a1" />
          <Text style={styles.infoText}>
            Solo vinculamos tu cuenta para transferirte dinero. No podemos hacer compras ni retirar
            dinero sin tu autorizacion.
          </Text>
        </Card>

        {/* Error */}
        {error && (
          <View style={styles.errorContainer}>
            <Ionicons name="alert-circle" size={20} color="#dc2626" />
            <Text style={styles.errorText}>{error}</Text>
          </View>
        )}

        {/* Manual Token Input */}
        {showManualInput && (
          <Card style={styles.manualInputCard}>
            <Text style={styles.cardTitle}>Vincular manualmente</Text>
            <Text style={styles.manualInputDescription}>
              Ingresa tu Access Token de MercadoPago. Lo puedes obtener en:{'\n'}
              https://www.mercadopago.com.ar/developers/panel/credentials
            </Text>

            <View style={styles.tokenInputContainer}>
              <TextInput
                style={styles.tokenInput}
                value={accessToken}
                onChangeText={setAccessToken}
                placeholder="Toca aquí y mantén presionado para pegar&#10;APP_USR-8687385123270840-121023-..."
                placeholderTextColor="#9ca3af"
                autoCapitalize="none"
                autoCorrect={false}
                autoComplete="off"
                textContentType="none"
                selectTextOnFocus={true}
                contextMenuHidden={false}
                multiline={true}
                numberOfLines={3}
                editable={true}
              />
              {accessToken.length > 0 && (
                <TouchableOpacity
                  style={styles.clearButton}
                  onPress={() => setAccessToken('')}
                >
                  <Ionicons name="close-circle" size={20} color="#6b7280" />
                </TouchableOpacity>
              )}
            </View>

            <Text style={styles.tokenHint}>
              Toca y mantén presionado dentro del cuadro para pegar
            </Text>

            <Button
              title={linking ? 'Vinculando...' : 'Vincular con token'}
              onPress={handleLinkManually}
              loading={linking}
              disabled={linking || !accessToken.trim()}
              style={styles.manualButton}
            />
          </Card>
        )}

        {/* CTA */}
        <View style={styles.ctaContainer}>
          <Button
            title={linking ? 'Conectando...' : 'Vincular cuenta de MercadoPago'}
            onPress={handleLinkAccount}
            loading={linking}
            disabled={linking}
            style={styles.ctaButton}
          />

          <Button
            title={showManualInput ? 'Usar OAuth' : 'Vincular manualmente'}
            variant="outline"
            onPress={() => setShowManualInput(!showManualInput)}
            disabled={linking}
            style={styles.toggleButton}
          />

          <Button
            title="Cancelar"
            variant="outline"
            onPress={() => navigation.goBack()}
            disabled={linking}
            style={styles.cancelButton}
          />
        </View>
      </ScrollView>
    </View>
  )
}
