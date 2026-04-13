import type { User } from './user'

export interface LoginPayload {
  email: string
  password: string
  tenant_domain: string
}

export interface LoginResponse {
  token: string
  user: User
}