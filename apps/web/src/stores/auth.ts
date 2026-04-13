import { defineStore } from 'pinia'
import { loginRequest } from '@/api/auth'
import { storage } from '@/utils/storage'
import type { User } from '@/types/user'

interface AuthState {
  token: string | null
  user: User | null
  tenant: string | null
  loading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: storage.getToken(),
    user: storage.getUser(),
    tenant: storage.getTenant(),
    loading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    role: (state) => state.user?.roles?.[0] || null,
  },

  actions: {
    hydrate() {
      this.token = storage.getToken()
      this.user = storage.getUser()
      this.tenant = storage.getTenant()
    },

    async login(payload: {
      email: string
      password: string
      tenant_domain: string
    }) {
      this.loading = true
      this.error = null

      try {
        // 🔥 IMPORTANT: set tenant BEFORE request
        storage.setTenant(payload.tenant_domain)

        const data = await loginRequest(payload)

        this.token = data.token
        this.user = data.user
        this.tenant = payload.tenant_domain

        storage.setToken(data.token)
        storage.setUser(data.user)

        return data
      } catch (err: any) {
        this.error = err?.response?.data?.error || 'Login failed'
        throw err
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.token = null
      this.user = null
      this.tenant = null
      storage.clear()
    },
  },
})