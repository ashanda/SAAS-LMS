<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  title: string
}>()

const router = useRouter()
const auth = useAuthStore()

const roleLabel = computed(() => auth.role || 'user')

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="dashboard-shell">
    <aside class="sidebar">
      <div class="sidebar-top">
        <div class="logo">
          <div class="logo-icon">L</div>
          <div>
            <h2>LMS SaaS</h2>
            <p>{{ roleLabel }}</p>
          </div>
        </div>
      </div>

      <nav class="menu">
        <RouterLink to="/teacher" class="menu-link" v-if="auth.role === 'teacher'">
          Teacher Dashboard
        </RouterLink>

        <RouterLink to="/student" class="menu-link" v-if="auth.role === 'student'">
          Student Dashboard
        </RouterLink>

        <RouterLink
          to="/admin"
          class="menu-link"
          v-if="['owner', 'admin', 'staff'].includes(auth.role || '')"
        >
          Admin Dashboard
        </RouterLink>
      </nav>

      <div class="sidebar-bottom">
        <button class="logout-btn" @click="handleLogout">Logout</button>
      </div>
    </aside>

    <div class="dashboard-main">
      <header class="topbar">
        <div>
          <h1>{{ props.title }}</h1>
          <p>{{ auth.user?.email }}</p>
        </div>
      </header>

      <main class="page-content">
        <slot />
      </main>
    </div>
  </div>
</template>

<style scoped>
.dashboard-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 280px 1fr;
  background: #f8fafc;
}

.sidebar {
  background: #0f172a;
  color: white;
  display: flex;
  flex-direction: column;
  padding: 22px;
  border-right: 1px solid rgba(255,255,255,0.06);
}

.sidebar-top {
  margin-bottom: 30px;
}

.logo {
  display: flex;
  gap: 14px;
  align-items: center;
}

.logo-icon {
  width: 46px;
  height: 46px;
  border-radius: 14px;
  background: linear-gradient(135deg, #4f46e5, #0ea5e9);
  display: grid;
  place-items: center;
  font-weight: 800;
  font-size: 20px;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.logo p {
  margin: 4px 0 0;
  color: rgba(255,255,255,0.65);
  font-size: 13px;
  text-transform: capitalize;
}

.menu {
  display: grid;
  gap: 10px;
}

.menu-link {
  color: rgba(255,255,255,0.84);
  text-decoration: none;
  padding: 13px 14px;
  border-radius: 14px;
  transition: 0.2s ease;
  background: rgba(255,255,255,0.03);
}

.menu-link:hover,
.router-link-active.menu-link {
  background: rgba(79, 70, 229, 0.24);
  color: white;
}

.sidebar-bottom {
  margin-top: auto;
}

.logout-btn {
  width: 100%;
  height: 48px;
  border: 0;
  border-radius: 14px;
  background: #ef4444;
  color: white;
  font-weight: 700;
  cursor: pointer;
}

.dashboard-main {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.topbar {
  background: white;
  border-bottom: 1px solid #e2e8f0;
  padding: 24px 30px;
}

.topbar h1 {
  margin: 0;
  font-size: 28px;
  color: #0f172a;
}

.topbar p {
  margin: 6px 0 0;
  color: #64748b;
}

.page-content {
  padding: 28px;
}

@media (max-width: 900px) {
  .dashboard-shell {
    grid-template-columns: 1fr;
  }

  .sidebar {
    padding: 16px;
  }

  .menu {
    grid-template-columns: 1fr;
  }

  .page-content {
    padding: 18px;
  }

  .topbar {
    padding: 18px;
  }
}
</style>