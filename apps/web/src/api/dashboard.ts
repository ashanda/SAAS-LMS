import api from './axios'
import type { AdminDashboardStats, Subject, Student } from '@/types/dashboard'

export async function getAdminDashboard(): Promise<AdminDashboardStats> {
  const { data } = await api.get<AdminDashboardStats>('/lms/dashboard/admin')
  return data
}

export async function getTeacherSubjects(): Promise<Subject[]> {
  const { data } = await api.get<Subject[]>('/lms/teacher/subjects')
  return data
}

export async function getTeacherStudents(subjectId: string): Promise<Student[]> {
  const { data } = await api.get<Student[]>(`/lms/teacher/subjects/${subjectId}/students`)
  return data
}

export async function getStudentSubjects(): Promise<Subject[]> {
  const { data } = await api.get<Subject[]>('/lms/dashboard/student/subjects')
  return data
}