import { useState, useEffect } from 'react'
import { CheckCircle, XCircle, AlertCircle, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import adminService from '@/services/adminService'
import { getErrorMessage } from '@/services'
import type { User } from '@/types'

export default function AdminVerificationsPage() {
  const [users, setUsers] = useState<User[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [processingId, setProcessingId] = useState<number | null>(null)
  const [actionError, setActionError] = useState('')

  // Reject modal state
  const [rejectModal, setRejectModal] = useState<{
    userId: number
    verificationType: 'dni' | 'license'
  } | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  useEffect(() => {
    fetchPending()
  }, [])

  const fetchPending = async () => {
    try {
      setLoading(true)
      setError('')
      const data = await adminService.getPendingVerifications()
      setUsers(data.users)
      setTotal(data.total)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (userId: number, verificationType: 'dni' | 'license') => {
    setActionError('')
    setProcessingId(userId)
    try {
      await adminService.reviewVerification(userId, 'approve', verificationType)
      await fetchPending()
    } catch (err) {
      setActionError(getErrorMessage(err))
    } finally {
      setProcessingId(null)
    }
  }

  const handleReject = async () => {
    if (!rejectModal) return
    setActionError('')
    setProcessingId(rejectModal.userId)
    try {
      await adminService.reviewVerification(
        rejectModal.userId,
        'reject',
        rejectModal.verificationType,
        rejectReason
      )
      setRejectModal(null)
      setRejectReason('')
      await fetchPending()
    } catch (err) {
      setActionError(getErrorMessage(err))
    } finally {
      setProcessingId(null)
    }
  }

  return (
    <>
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <ShieldCheck className="w-7 h-7 text-primary" />
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Verificaciones de Identidad</h1>
            <p className="text-sm text-gray-500">{total} solicitud(es) pendiente(s)</p>
          </div>
        </div>
        <Button variant="outline" onClick={fetchPending} disabled={loading}>
          Actualizar
        </Button>
      </div>

      {actionError && (
        <div className="mb-4 p-3 bg-destructive/10 border border-destructive/20 rounded-lg flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-destructive" />
          <p className="text-sm text-destructive">{actionError}</p>
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center py-16">
          <div className="w-10 h-10 border-4 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      ) : error ? (
        <div className="p-4 bg-destructive/10 border border-destructive/20 rounded-lg">
          <p className="text-destructive">{error}</p>
        </div>
      ) : users.length === 0 ? (
        <Card className="text-center py-12">
          <CardContent>
            <ShieldCheck className="w-12 h-12 text-green-400 mx-auto mb-3" />
            <p className="text-lg font-medium text-gray-700">No hay verificaciones pendientes</p>
            <p className="text-sm text-gray-400">Todo está al día.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {users.map((user) => (
            <Card key={user.id} className="overflow-hidden">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-base">
                      {user.name} {user.lastname}
                    </CardTitle>
                    <p className="text-sm text-gray-500">{user.email}</p>
                    <p className="text-xs text-gray-400 mt-1">ID: {user.id}</p>
                  </div>
                  <Badge variant="secondary" className="flex items-center gap-1">
                    <AlertCircle className="w-3 h-3" />
                    Pendiente
                  </Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {/* DNI section */}
                  {((user as any).dni || (user as any).dni_photo_url) && (
                    <div className="p-3 bg-gray-50 rounded-lg space-y-2">
                      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">DNI</p>
                      {(user as any).dni && (
                        <p className="text-sm"><span className="font-medium">Número:</span> {(user as any).dni}</p>
                      )}
                      {(user as any).dni_photo_url && (
                        <div>
                          <p className="text-sm font-medium mb-1">Foto:</p>
                          <a
                            href={(user as any).dni_photo_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-sm text-primary underline break-all"
                          >
                            Ver foto del DNI
                          </a>
                        </div>
                      )}
                      {!(user as any).dni_verified && (
                        <div className="flex gap-2 pt-1">
                          <Button
                            size="sm"
                            onClick={() => handleApprove(user.id, 'dni')}
                            disabled={processingId === user.id}
                          >
                            <CheckCircle className="w-3 h-3 mr-1" />
                            Aprobar DNI
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => setRejectModal({ userId: user.id, verificationType: 'dni' })}
                            disabled={processingId === user.id}
                          >
                            <XCircle className="w-3 h-3 mr-1" />
                            Rechazar
                          </Button>
                        </div>
                      )}
                      {(user as any).dni_verified && (
                        <Badge variant="default" className="text-xs">
                          <CheckCircle className="w-3 h-3 mr-1" /> DNI aprobado
                        </Badge>
                      )}
                    </div>
                  )}

                  {/* License section */}
                  {((user as any).license_number || (user as any).license_photo_url) && (
                    <div className="p-3 bg-gray-50 rounded-lg space-y-2">
                      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wide">Licencia de Conducir</p>
                      {(user as any).license_number && (
                        <p className="text-sm"><span className="font-medium">Número:</span> {(user as any).license_number}</p>
                      )}
                      {(user as any).license_photo_url && (
                        <div>
                          <p className="text-sm font-medium mb-1">Foto:</p>
                          <a
                            href={(user as any).license_photo_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-sm text-primary underline break-all"
                          >
                            Ver foto de la licencia
                          </a>
                        </div>
                      )}
                      {!(user as any).license_verified && (
                        <div className="flex gap-2 pt-1">
                          <Button
                            size="sm"
                            onClick={() => handleApprove(user.id, 'license')}
                            disabled={processingId === user.id}
                          >
                            <CheckCircle className="w-3 h-3 mr-1" />
                            Aprobar Licencia
                          </Button>
                          <Button
                            size="sm"
                            variant="destructive"
                            onClick={() => setRejectModal({ userId: user.id, verificationType: 'license' })}
                            disabled={processingId === user.id}
                          >
                            <XCircle className="w-3 h-3 mr-1" />
                            Rechazar
                          </Button>
                        </div>
                      )}
                      {(user as any).license_verified && (
                        <Badge variant="default" className="text-xs">
                          <CheckCircle className="w-3 h-3 mr-1" /> Licencia aprobada
                        </Badge>
                      )}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>

    {/* Reject modal */}
    {rejectModal && (
      <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div className="bg-white rounded-xl shadow-xl p-6 max-w-sm w-full">
          <h3 className="text-lg font-semibold mb-2">Rechazar verificación</h3>
          <p className="text-sm text-gray-600 mb-4">
            Indicá el motivo del rechazo. El usuario podrá volver a enviar su documentación.
          </p>
          <textarea
            value={rejectReason}
            onChange={(e) => setRejectReason(e.target.value)}
            placeholder="Motivo del rechazo..."
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary mb-4"
          />
          {actionError && (
            <p className="text-sm text-destructive mb-3">{actionError}</p>
          )}
          <div className="flex gap-3">
            <Button
              variant="outline"
              onClick={() => { setRejectModal(null); setRejectReason('') }}
              disabled={processingId !== null}
              className="flex-1"
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleReject}
              disabled={processingId !== null || !rejectReason.trim()}
              className="flex-1"
            >
              {processingId !== null ? (
                <><div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />Procesando...</>
              ) : 'Rechazar'}
            </Button>
          </div>
        </div>
      </div>
    )}
    </>
  )
}
