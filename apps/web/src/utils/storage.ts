import type { User } from '@/types/user'

const TOKEN_KEY = 'lms_token'
const USER_KEY = 'lms_user'
const TENANT_KEY = 'lms_tenant'

export const storage = {
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  },

  getUser(): User | null {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? JSON.parse(raw) : null
  },

  setUser(user: User): void {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },

  getTenant(): string | null {
    return localStorage.getItem(TENANT_KEY)
  },

  setTenant(tenant: string): void {
    localStorage.setItem(TENANT_KEY, tenant)
  },

  clear(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(TENANT_KEY)
  },
}