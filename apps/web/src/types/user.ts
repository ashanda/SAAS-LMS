export type UserRole = 'owner' | 'admin' | 'staff' | 'teacher' | 'student'

export interface User {
  id: string
  name: string
  email: string
  roles: UserRole
}