import api from './axios'

export async function loginRequest(payload: {
  email: string
  password: string
  tenant_domain: string
}) {
  const { data } = await api.post('/auth/login', payload)
  return data
}