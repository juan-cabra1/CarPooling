import { useState } from 'react'
import { useAuth } from '@/context/AuthContext'
import { authService } from '@/services'
import type { UpdateUserData, ChangePasswordRequest, UserSex } from '@/types'

export default function ProfilePage() {
  const { user, refreshUser, logout } = useAuth()

  // Personal info form state
  const [profileData, setProfileData] = useState<UpdateUserData>({
    name: user?.name ?? '',
    lastname: user?.lastname ?? '',
    phone: user?.phone ?? '',
    street: user?.street ?? '',
    number: user?.number ?? 0,
  })
  const [profileSuccess, setProfileSuccess] = useState('')
  const [profileError, setProfileError] = useState('')
  const [profileLoading, setProfileLoading] = useState(false)

  // Password form state
  const [pwData, setPwData] = useState<ChangePasswordRequest>({
    current_password: '',
    new_password: '',
  })
  const [confirmPassword, setConfirmPassword] = useState('')
  const [pwSuccess, setPwSuccess] = useState('')
  const [pwError, setPwError] = useState('')
  const [pwLoading, setPwLoading] = useState(false)

  // Delete account modal
  const [showDeleteModal, setShowDeleteModal] = useState(false)
  const [deleteLoading, setDeleteLoading] = useState(false)
  const [deleteError, setDeleteError] = useState('')

  if (!user) return null

  const handleProfileSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setProfileError('')
    setProfileSuccess('')
    setProfileLoading(true)
    try {
      await authService.updateProfile(user.id, profileData)
      await refreshUser()
      setProfileSuccess('Cambios guardados correctamente')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al guardar cambios'
      setProfileError(msg)
    } finally {
      setProfileLoading(false)
    }
  }

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setPwError('')
    setPwSuccess('')
    if (pwData.new_password !== confirmPassword) {
      setPwError('Las contraseñas no coinciden')
      return
    }
    if (pwData.new_password.length < 8) {
      setPwError('La nueva contraseña debe tener al menos 8 caracteres')
      return
    }
    setPwLoading(true)
    try {
      await authService.changePassword(pwData)
      setPwSuccess('Contraseña cambiada correctamente')
      setPwData({ current_password: '', new_password: '' })
      setConfirmPassword('')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al cambiar contraseña'
      setPwError(msg)
    } finally {
      setPwLoading(false)
    }
  }

  const handleDeleteAccount = async () => {
    setDeleteError('')
    setDeleteLoading(true)
    try {
      await authService.deleteAccount(user.id)
      logout()
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar cuenta'
      setDeleteError(msg)
      setDeleteLoading(false)
    }
  }

  const sexOptions: { value: UserSex; label: string }[] = [
    { value: 'hombre', label: 'Hombre' },
    { value: 'mujer', label: 'Mujer' },
    { value: 'otro', label: 'Otro' },
  ]

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-8">Mi Perfil</h1>

      {/* Personal Info Section */}
      <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">Información Personal</h2>
        <form onSubmit={handleProfileSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
              <input
                type="text"
                value={profileData.name ?? ''}
                onChange={e => setProfileData(d => ({ ...d, name: e.target.value }))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Apellido</label>
              <input
                type="text"
                value={profileData.lastname ?? ''}
                onChange={e => setProfileData(d => ({ ...d, lastname: e.target.value }))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Teléfono</label>
            <input
              type="tel"
              value={profileData.phone ?? ''}
              onChange={e => setProfileData(d => ({ ...d, phone: e.target.value }))}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700 mb-1">Calle</label>
              <input
                type="text"
                value={profileData.street ?? ''}
                onChange={e => setProfileData(d => ({ ...d, street: e.target.value }))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Número</label>
              <input
                type="number"
                value={profileData.number ?? 0}
                onChange={e => setProfileData(d => ({ ...d, number: parseInt(e.target.value) || 0 }))}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Sexo</label>
              <select
                value={user.sex}
                disabled
                className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm bg-gray-50 text-gray-500 cursor-not-allowed"
              >
                {sexOptions.map(o => (
                  <option key={o.value} value={o.value}>{o.label}</option>
                ))}
              </select>
              <p className="text-xs text-gray-400 mt-1">No editable</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Fecha de nacimiento</label>
              <input
                type="date"
                value={user.birthdate?.slice(0, 10) ?? ''}
                disabled
                className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm bg-gray-50 text-gray-500 cursor-not-allowed"
              />
              <p className="text-xs text-gray-400 mt-1">No editable</p>
            </div>
          </div>

          {profileSuccess && (
            <p className="text-sm text-green-600 bg-green-50 border border-green-200 rounded-lg px-3 py-2">{profileSuccess}</p>
          )}
          {profileError && (
            <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">{profileError}</p>
          )}

          <button
            type="submit"
            disabled={profileLoading}
            className="w-full bg-blue-600 text-white font-medium py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {profileLoading ? 'Guardando...' : 'Guardar cambios'}
          </button>
        </form>
      </section>

      {/* Password Section */}
      <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">Seguridad</h2>
        <form onSubmit={handlePasswordSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Contraseña actual</label>
            <input
              type="password"
              value={pwData.current_password}
              onChange={e => setPwData(d => ({ ...d, current_password: e.target.value }))}
              required
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Nueva contraseña</label>
            <input
              type="password"
              value={pwData.new_password}
              onChange={e => setPwData(d => ({ ...d, new_password: e.target.value }))}
              required
              minLength={8}
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Confirmar nueva contraseña</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
              required
              className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {pwSuccess && (
            <p className="text-sm text-green-600 bg-green-50 border border-green-200 rounded-lg px-3 py-2">{pwSuccess}</p>
          )}
          {pwError && (
            <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">{pwError}</p>
          )}

          <button
            type="submit"
            disabled={pwLoading}
            className="w-full bg-blue-600 text-white font-medium py-2 px-4 rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {pwLoading ? 'Cambiando...' : 'Cambiar contraseña'}
          </button>
        </form>
      </section>

      {/* Account Section */}
      <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h2 className="text-xl font-semibold mb-4">Cuenta</h2>

        <div className="space-y-3 mb-6">
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Email verificado</span>
            <span className={`text-xs font-medium px-2 py-1 rounded-full ${
              user.email_verified
                ? 'bg-green-100 text-green-700'
                : 'bg-yellow-100 text-yellow-700'
            }`}>
              {user.email_verified ? 'Verificado' : 'Sin verificar'}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Rating como conductor</span>
            <span className="text-sm font-medium">
              {user.avg_driver_rating > 0
                ? `${user.avg_driver_rating.toFixed(1)} ⭐ (${user.total_trips_driver} viajes)`
                : 'Sin viajes como conductor'}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Rating como pasajero</span>
            <span className="text-sm font-medium">
              {user.avg_passenger_rating > 0
                ? `${user.avg_passenger_rating.toFixed(1)} ⭐ (${user.total_trips_passenger} viajes)`
                : 'Sin viajes como pasajero'}
            </span>
          </div>
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Rol</span>
            <span className={`text-xs font-medium px-2 py-1 rounded-full ${
              user.role === 'admin' ? 'bg-purple-100 text-purple-700' : 'bg-gray-100 text-gray-700'
            }`}>
              {user.role === 'admin' ? 'Administrador' : 'Usuario'}
            </span>
          </div>
        </div>

        <button
          onClick={() => setShowDeleteModal(true)}
          className="w-full border border-red-300 text-red-600 font-medium py-2 px-4 rounded-lg hover:bg-red-50 transition-colors"
        >
          Eliminar cuenta
        </button>
      </section>

      {/* Delete Confirmation Modal */}
      {showDeleteModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl p-6 max-w-sm w-full">
            <h3 className="text-lg font-semibold mb-2">Eliminar cuenta</h3>
            <p className="text-sm text-gray-600 mb-4">
              Esta acción es permanente y no se puede deshacer. Todos tus datos serán eliminados.
            </p>
            {deleteError && (
              <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2 mb-4">{deleteError}</p>
            )}
            <div className="flex gap-3">
              <button
                onClick={() => { setShowDeleteModal(false); setDeleteError('') }}
                disabled={deleteLoading}
                className="flex-1 border border-gray-300 text-gray-700 font-medium py-2 px-4 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancelar
              </button>
              <button
                onClick={handleDeleteAccount}
                disabled={deleteLoading}
                className="flex-1 bg-red-600 text-white font-medium py-2 px-4 rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors"
              >
                {deleteLoading ? 'Eliminando...' : 'Eliminar'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
