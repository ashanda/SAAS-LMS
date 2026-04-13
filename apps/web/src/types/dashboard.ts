export interface AdminDashboardStats {
  students: number
  teachers: number
  subjects: number
  enrollments: number
}

export interface Subject {
  id: string
  name: string
  grade_name?: string
}

export interface Student {
  id: string
  name: string
  email?: string
}