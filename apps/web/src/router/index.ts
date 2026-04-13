import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// ✅ extend meta typing (IMPORTANT)
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    guestOnly?: boolean
    roles?: string[]
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/admin',
    name: 'admin-dashboard',
    component: () => import('@/views/admin/AdminDashboardView.vue'),
    meta: { requiresAuth: true, roles: ['owner', 'admin', 'staff'] },
  },
  {
    path: '/teacher',
    name: 'teacher-dashboard',
    component: () => import('@/views/teacher/TeacherDashboardView.vue'),
    meta: { requiresAuth: true, roles: ['teacher'] },
  },
  {
    path: '/student',
    name: 'student-dashboard',
    component: () => import('@/views/student/StudentDashboardView.vue'),
    meta: { requiresAuth: true, roles: ['student'] },
  },
  {
    path: '/unauthorized',
    name: 'unauthorized',
    component: () => import('@/views/UnauthorizedView.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  // ✅ SAFE hydrate (only if exists)
  if (typeof auth.hydrate === 'function') {
    auth.hydrate()
  }

  // =========================
  // 🔐 AUTH CHECK
  // =========================
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' }
  }

  // =========================
  // 🚫 GUEST CHECK
  // =========================
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return redirectByRole(auth.role)
  }

  // =========================
  // 🎭 ROLE CHECK
  // =========================
  if (to.meta.roles && to.meta.roles.length > 0) {
    const role = auth.role

    if (!role || !to.meta.roles.includes(role)) {
      return { name: 'unauthorized' }
    }
  }

  return true
})

// =========================
// 🔁 REDIRECT BY ROLE
// =========================
function redirectByRole(role: string | null) {
  if (!role) return { name: 'login' }

  if (['owner', 'admin', 'staff'].includes(role)) {
    return { name: 'admin-dashboard' }
  }

  if (role === 'teacher') {
    return { name: 'teacher-dashboard' }
  }

  if (role === 'student') {
    return { name: 'student-dashboard' }
  }

  return { name: 'login' }
}

export default router