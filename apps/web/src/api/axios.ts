import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// 🔥 FINAL FIX
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('lms_token')
  const tenant = localStorage.getItem('lms_tenant')

  // ✅ AUTH HEADER (CRITICAL)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  // ✅ TENANT HEADER
  if (tenant) {
    config.headers['X-Tenant-Domain'] = tenant
  }

  return config
})

export default api