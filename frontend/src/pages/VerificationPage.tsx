import { useState, useEffect } from 'react'
import { CheckCircle, Clock, XCircle, AlertCircle, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { usersApi } from '@/services/api'
import { getErrorMessage } from '@/services'
import { useAuth } from '@/context/AuthContext'
import type { ApiResponse } from '@/types'

interface VerificationStatus {
  verification_status: 'none' | 'pending' | 'verified' | 'rejected'
  dni_verified: boolean
  dni_verified_at: string | null
  license_verified: boolean
  license_verified_at: string | null
  rejection_reason: string
}

interface SubmitVerificationData {
  dni?: string
  dni_photo_url?: string
  license_number?: string
  license_photo_url?: string
}

export default function VerificationPage() {
  const { user } = useAuth()

  const [status, setStatus] = useState<VerificationStatus | null>(null)
  const [statusLoading, setStatusLoading] = useState(true)
  const [statusError, setStatusError] = useState('')

  const [dni, setDni] = useState('')
  const [dniPhotoUrl, setDniPhotoUrl] = useState('')
  const [licenseNumber, setLicenseNumber] = useState('')
  const [licensePhotoUrl, setLicensePhotoUrl] = useState('')

  const [submitting, setSubmitting] = useState(false)
  const [submitSuccess, setSubmitSuccess] = useState('')
  const [submitError, setSubmitError] = useState('')

  useEffect(() => {
    if (user?.id) {
      fetchStatus()
    }
  }, [user])

  const fetchStatus = async () => {
    try {
      setStatusLoading(true)
      setStatusError('')
      const response = await usersApi.get<ApiResponse<VerificationStatus>>(
        `/users/${user!.id}/verification`
      )
      setStatus(response.data.data!)
    } catch (err) {
      setStatusError(getErrorMessage(err))
    } finally {
      setStatusLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitError('')
    setSubmitSuccess('')

    const data: SubmitVerificationData = {}
    if (dni.trim()) data.dni = dni.trim()
    if (dniPhotoUrl.trim()) data.dni_photo_url = dniPhotoUrl.trim()
    if (licenseNumber.trim()) data.license_number = licenseNumber.trim()
    if (licensePhotoUrl.trim()) data.license_photo_url = licensePhotoUrl.trim()

    if (Object.keys(data).length === 0) {
      setSubmitError('Completá al menos un campo de verificación')
      return
    }

    try {
      setSubmitting(true)
      await usersApi.post(`/users/${user!.id}/verification`, data)
      setSubmitSuccess('Información enviada correctamente. Será revisada por un administrador.')
      await fetchStatus()
      setDni('')
      setDniPhotoUrl('')
      setLicenseNumber('')
      setLicensePhotoUrl('')
    } catch (err) {
      setSubmitError(getErrorMessage(err))
    } finally {
      setSubmitting(false)
    }
  }

  const getStatusBadge = (verStatus: string) => {
    const config: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; icon: React.ReactNode }> = {
      none: { label: 'Sin verificar', variant: 'outline', icon: <AlertCircle className="w-4 h-4" /> },
      pending: { label: 'En revisión', variant: 'secondary', icon: <Clock className="w-4 h-4" /> },
      verified: { label: 'Verificado', variant: 'default', icon: <CheckCircle className="w-4 h-4" /> },
      rejected: { label: 'Rechazado', variant: 'destructive', icon: <XCircle className="w-4 h-4" /> },
    }
    const c = config[verStatus] || config.none
    return (
      <Badge variant={c.variant} className="flex items-center gap-1 px-3 py-1 text-sm">
        {c.icon}
        {c.label}
      </Badge>
    )
  }

  const formatDate = (dateString: string | null) => {
    if (!dateString) return '-'
    return new Intl.DateTimeFormat('es-AR', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    }).format(new Date(dateString))
  }

  const canSubmit = !status || status.verification_status === 'none' || status.verification_status === 'rejected'

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-2xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center gap-3 mb-2">
              <ShieldCheck className="w-8 h-8 text-primary" />
              <h1 className="text-3xl font-bold text-gray-900">Verificación de Identidad</h1>
            </div>
            <p className="text-gray-600">
              Verificá tu DNI y/o licencia de conducir para generar más confianza en la comunidad.
            </p>
          </div>

          {/* Current Status */}
          <Card className="mb-6">
            <CardHeader>
              <CardTitle className="text-lg">Estado de Verificación</CardTitle>
            </CardHeader>
            <CardContent>
              {statusLoading ? (
                <div className="flex items-center gap-2 text-gray-500">
                  <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                  <span className="text-sm">Cargando estado...</span>
                </div>
              ) : statusError ? (
                <p className="text-sm text-destructive">{statusError}</p>
              ) : status ? (
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-gray-700">Estado general</span>
                    {getStatusBadge(status.verification_status)}
                  </div>

                  {status.rejection_reason && (
                    <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-lg">
                      <p className="text-sm text-destructive">
                        <span className="font-medium">Motivo de rechazo:</span> {status.rejection_reason}
                      </p>
                    </div>
                  )}

                  <div className="grid grid-cols-2 gap-4 pt-2 border-t">
                    <div>
                      <p className="text-xs text-gray-500 mb-1">DNI</p>
                      {status.dni_verified ? (
                        <div className="flex items-center gap-1 text-green-600">
                          <CheckCircle className="w-4 h-4" />
                          <span className="text-sm font-medium">Verificado</span>
                        </div>
                      ) : (
                        <div className="flex items-center gap-1 text-gray-400">
                          <XCircle className="w-4 h-4" />
                          <span className="text-sm">No verificado</span>
                        </div>
                      )}
                      {status.dni_verified_at && (
                        <p className="text-xs text-gray-400 mt-1">{formatDate(status.dni_verified_at)}</p>
                      )}
                    </div>
                    <div>
                      <p className="text-xs text-gray-500 mb-1">Licencia de conducir</p>
                      {status.license_verified ? (
                        <div className="flex items-center gap-1 text-green-600">
                          <CheckCircle className="w-4 h-4" />
                          <span className="text-sm font-medium">Verificada</span>
                        </div>
                      ) : (
                        <div className="flex items-center gap-1 text-gray-400">
                          <XCircle className="w-4 h-4" />
                          <span className="text-sm">No verificada</span>
                        </div>
                      )}
                      {status.license_verified_at && (
                        <p className="text-xs text-gray-400 mt-1">{formatDate(status.license_verified_at)}</p>
                      )}
                    </div>
                  </div>
                </div>
              ) : null}
            </CardContent>
          </Card>

          {/* Pending notice */}
          {status?.verification_status === 'pending' && (
            <div className="mb-6 p-4 rounded-lg bg-yellow-50 border border-yellow-200 flex items-start gap-3">
              <Clock className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
              <div>
                <p className="font-semibold text-yellow-800">Verificación en proceso</p>
                <p className="text-sm text-yellow-700">
                  Tu solicitud está siendo revisada por nuestro equipo. Te notificaremos cuando tengamos novedades.
                </p>
              </div>
            </div>
          )}

          {/* Submit form */}
          {canSubmit && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">
                  {status?.verification_status === 'rejected' ? 'Volver a enviar documentación' : 'Enviar Documentación'}
                </CardTitle>
                <CardDescription>
                  Completá los campos con tu información. Podés enviar solo el DNI, solo la licencia, o ambos a la vez.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                  {/* DNI section */}
                  <div className="space-y-3">
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide">DNI</h3>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Número de DNI
                      </label>
                      <input
                        type="text"
                        value={dni}
                        onChange={(e) => setDni(e.target.value)}
                        placeholder="Ej: 12345678"
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        URL de foto del DNI
                      </label>
                      <input
                        type="url"
                        value={dniPhotoUrl}
                        onChange={(e) => setDniPhotoUrl(e.target.value)}
                        placeholder="https://..."
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                      <p className="text-xs text-gray-400 mt-1">
                        Subí la foto a un servicio externo (Imgur, Google Drive, etc.) y pegá el enlace.
                      </p>
                    </div>
                  </div>

                  <div className="border-t" />

                  {/* License section */}
                  <div className="space-y-3">
                    <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wide">Licencia de Conducir</h3>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Número de licencia
                      </label>
                      <input
                        type="text"
                        value={licenseNumber}
                        onChange={(e) => setLicenseNumber(e.target.value)}
                        placeholder="Ej: 12345678"
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        URL de foto de la licencia
                      </label>
                      <input
                        type="url"
                        value={licensePhotoUrl}
                        onChange={(e) => setLicensePhotoUrl(e.target.value)}
                        placeholder="https://..."
                        className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                      <p className="text-xs text-gray-400 mt-1">
                        Subí la foto a un servicio externo (Imgur, Google Drive, etc.) y pegá el enlace.
                      </p>
                    </div>
                  </div>

                  {submitError && (
                    <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-lg flex items-start gap-2">
                      <AlertCircle className="w-4 h-4 text-destructive flex-shrink-0 mt-0.5" />
                      <p className="text-sm text-destructive">{submitError}</p>
                    </div>
                  )}

                  {submitSuccess && (
                    <div className="p-3 bg-green-50 border border-green-200 rounded-lg flex items-start gap-2">
                      <CheckCircle className="w-4 h-4 text-green-600 flex-shrink-0 mt-0.5" />
                      <p className="text-sm text-green-700">{submitSuccess}</p>
                    </div>
                  )}

                  <Button type="submit" disabled={submitting} className="w-full">
                    {submitting ? (
                      <><div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />Enviando...</>
                    ) : 'Enviar para revisión'}
                  </Button>
                </form>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
