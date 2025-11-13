import React, { useState } from 'react';
import { Card, Button } from '@/components/common';
import axios from 'axios';

export const RegisterDebug: React.FC = () => {
  const [response, setResponse] = useState<any>(null);
  const [error, setError] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const testData = {
    email: 'test@example.com',
    password: 'password123',
    name: 'Test',
    lastname: 'User',
    phone: '+54 11 1234-5678',
    street: 'Test Street',
    number: 123,
    sex: 'otro',
    birthdate: '2000-01-01',
  };

  const testRegister = async () => {
    setLoading(true);
    setError(null);
    setResponse(null);

    try {
      console.log('🚀 Enviando request a: http://localhost:8001/users');
      console.log('📦 Datos:', testData);

      const result = await axios.post('http://localhost:8001/users', testData, {
        headers: {
          'Content-Type': 'application/json',
        },
      });

      console.log('✅ Respuesta exitosa:', result.data);
      setResponse(result.data);
    } catch (err: any) {
      console.error('❌ Error:', err);
      setError({
        message: err.message,
        status: err.response?.status,
        statusText: err.response?.statusText,
        data: err.response?.data,
        url: err.config?.url,
        method: err.config?.method,
        requestData: err.config?.data,
      });
    } finally {
      setLoading(false);
    }
  };

  const testHealth = async () => {
    setLoading(true);
    setError(null);
    setResponse(null);

    try {
      console.log('🏥 Probando health check...');
      const result = await axios.get('http://localhost:8001/health');
      console.log('✅ Health OK:', result.data);
      setResponse(result.data);
    } catch (err: any) {
      console.error('❌ Error en health:', err);
      setError({ message: err.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-8">🔍 Debug - Test de Registro</h1>

        <Card className="mb-6">
          <h2 className="text-xl font-bold mb-4">Datos de Prueba</h2>
          <pre className="bg-gray-100 p-4 rounded overflow-x-auto">
            {JSON.stringify(testData, null, 2)}
          </pre>
        </Card>

        <div className="flex gap-4 mb-6">
          <Button onClick={testHealth} isLoading={loading} variant="secondary">
            Test Health Check
          </Button>
          <Button onClick={testRegister} isLoading={loading} variant="primary">
            Test Register
          </Button>
        </div>

        {response && (
          <Card className="mb-6 bg-green-50">
            <h2 className="text-xl font-bold mb-4 text-green-700">✅ Respuesta Exitosa</h2>
            <pre className="bg-white p-4 rounded overflow-x-auto text-sm">
              {JSON.stringify(response, null, 2)}
            </pre>
          </Card>
        )}

        {error && (
          <Card className="mb-6 bg-red-50">
            <h2 className="text-xl font-bold mb-4 text-red-700">❌ Error</h2>
            <div className="space-y-4">
              <div>
                <h3 className="font-semibold">Mensaje:</h3>
                <p className="text-red-600">{error.message}</p>
              </div>
              {error.status && (
                <div>
                  <h3 className="font-semibold">Status:</h3>
                  <p className="text-red-600">{error.status} - {error.statusText}</p>
                </div>
              )}
              {error.data && (
                <div>
                  <h3 className="font-semibold">Respuesta del Backend:</h3>
                  <pre className="bg-white p-4 rounded overflow-x-auto text-sm">
                    {JSON.stringify(error.data, null, 2)}
                  </pre>
                </div>
              )}
              {error.requestData && (
                <div>
                  <h3 className="font-semibold">Datos Enviados:</h3>
                  <pre className="bg-white p-4 rounded overflow-x-auto text-sm">
                    {error.requestData}
                  </pre>
                </div>
              )}
            </div>
          </Card>
        )}

        <Card>
          <h2 className="text-xl font-bold mb-4">📋 Checklist</h2>
          <div className="space-y-2">
            <div>✅ Backend corriendo en http://localhost:8001</div>
            <div>✅ CORS configurado en backend</div>
            <div>✅ Base de datos MySQL conectada</div>
            <div>✅ Email service configurado (puede fallar pero no debería afectar)</div>
            <div>✅ Todos los campos requeridos en el request</div>
            <div>✅ Formato de fecha: YYYY-MM-DD</div>
            <div>✅ Sexo: 'hombre', 'mujer', o 'otro'</div>
          </div>
        </Card>

        <Card className="mt-6">
          <h2 className="text-xl font-bold mb-4">🛠️ Comandos Útiles</h2>
          <div className="space-y-2 text-sm">
            <div className="bg-gray-100 p-2 rounded">
              <code>cd backend/users-api && go run cmd/api/main.go</code>
            </div>
            <div className="bg-gray-100 p-2 rounded">
              <code>curl http://localhost:8001/health</code>
            </div>
            <div className="bg-gray-100 p-2 rounded">
              <code>docker ps | grep mysql</code>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};
