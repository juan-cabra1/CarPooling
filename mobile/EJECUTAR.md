# Guía Rápida de Ejecución - CarPooling Mobile

## Pasos para Ejecutar la Aplicación

### 1. Instalar Dependencias

```bash
cd frontendmobile
npm install
```

### 2. Asegurarse de que los Backends Estén Corriendo

Antes de ejecutar la app móvil, verifica que todos los servicios backend estén activos:

```bash
# En terminales separadas, desde el directorio raíz del proyecto:

# Terminal 1 - Users API (Puerto 8080)
cd backend/users-api
go run cmd/api/main.go

# Terminal 2 - Trips API (Puerto 8081)
cd backend/trips-api
go run cmd/api/main.go

# Terminal 3 - Bookings API (Puerto 8082)
cd backend/bookings-api
go run cmd/api/main.go

# Terminal 4 - Search API (Puerto 8083)
cd backend/search-api
go run cmd/api/main.go
```

### 3. Ejecutar la Aplicación Móvil

```bash
cd frontendmobile
npm start
```

Esto abrirá Expo DevTools en tu navegador con un código QR.

### 4. Opciones para Ver la Aplicación

#### Opción A: En tu Teléfono (Recomendado)

1. **Instala Expo Go** en tu dispositivo:
   - iOS: [Expo Go en App Store](https://apps.apple.com/app/expo-go/id982107779)
   - Android: [Expo Go en Google Play](https://play.google.com/store/apps/details?id=host.exp.exponent)

2. **Escanea el código QR** que aparece en la terminal o en el navegador:
   - iOS: Usa la cámara del iPhone
   - Android: Usa la app Expo Go

3. Espera a que se cargue la aplicación (puede tardar un poco la primera vez)

#### Opción B: En Emulador de Android

1. Asegúrate de tener Android Studio instalado con un emulador configurado
2. Inicia el emulador desde Android Studio
3. En la terminal de Expo, presiona `a` o ejecuta:
   ```bash
   npm run android
   ```

#### Opción C: En Simulador de iOS (Solo macOS)

1. Asegúrate de tener Xcode instalado
2. En la terminal de Expo, presiona `i` o ejecuta:
   ```bash
   npm run ios
   ```

#### Opción D: En el Navegador Web (Funcionalidad Limitada)

```bash
npm run web
```

**Nota**: La versión web tiene limitaciones ya que algunos componentes nativos no funcionan en el navegador.

### 5. Comandos Útiles

```bash
# Limpiar caché de Expo (si hay problemas)
npx expo start -c

# Ejecutar en modo desarrollo con reload automático
npm start

# Ver logs en tiempo real
# Los logs aparecerán en la terminal de Expo

# Detener el servidor
# Presiona Ctrl+C en la terminal donde ejecutaste npm start
```

## Solución de Problemas Comunes

### "No se puede conectar al servidor de desarrollo"

**En dispositivo físico:**
1. Asegúrate de que tu teléfono y computadora estén en la misma red WiFi
2. Verifica que el firewall no esté bloqueando el puerto 8081
3. Intenta reiniciar Expo con `npx expo start -c`

### "Error al conectar con las APIs"

1. Verifica que los 4 servicios backend estén corriendo
2. En dispositivo físico, es posible que necesites cambiar `localhost` por la IP local de tu computadora
3. Puedes encontrar tu IP local ejecutando:
   - Windows: `ipconfig`
   - macOS/Linux: `ifconfig` o `ip addr`

### "Metro bundler no inicia"

```bash
# Matar procesos existentes de Metro
killall node
# o en Windows
taskkill /F /IM node.exe

# Limpiar caché e iniciar de nuevo
npx expo start -c
```

### "Errores de dependencias"

```bash
# Borrar node_modules y reinstalar
rm -rf node_modules package-lock.json
npm install
```

## Características Disponibles

Una vez que la aplicación esté corriendo, puedes probar:

✅ **Autenticación**
- Registro de nuevo usuario
- Inicio de sesión
- Recuperación de contraseña
- Cierre de sesión

✅ **Búsqueda de Viajes**
- Buscar por origen y destino
- Filtrar por fecha
- Filtrar por número de pasajeros
- Ver detalles del viaje

✅ **Gestión de Viajes (Como Conductor)**
- Crear nuevo viaje
- Editar viaje existente
- Ver mis viajes publicados
- Eliminar viaje

✅ **Reservas (Como Pasajero)**
- Reservar asientos en un viaje
- Ver mis reservas activas
- Cancelar reserva

✅ **Perfil de Usuario**
- Ver información personal
- Editar perfil
- Cambiar foto de perfil

## Flujo de Prueba Recomendado

1. **Registro**: Crea una cuenta nueva
2. **Verificación**: (Opcional, según configuración del backend)
3. **Login**: Inicia sesión con tu cuenta
4. **Búsqueda**: Busca viajes disponibles
5. **Crear Viaje**: Publica un nuevo viaje como conductor
6. **Ver Viaje**: Revisa los detalles del viaje creado
7. **Reserva**: Con otra cuenta, reserva asientos en el viaje
8. **Perfil**: Actualiza tu información de perfil

## Configuración de APIs

Las URLs de las APIs están configuradas en `/src/services/api.ts`:

```typescript
Users API:    http://localhost:8080/api
Trips API:    http://localhost:8081/api
Bookings API: http://localhost:8082/api
Search API:   http://localhost:8083/api
```

Si necesitas cambiar estas URLs (por ejemplo, para usar tu IP local en dispositivo físico), edita el archivo `src/services/api.ts`.

## Más Ayuda

- Documentación completa: Ver [README.md](README.md)
- Documentación de Expo: https://docs.expo.dev/
- Documentación de React Native: https://reactnative.dev/

---

**¡Listo!** Ahora deberías poder ejecutar y probar la aplicación móvil de CarPooling.
