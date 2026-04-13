<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const form = reactive({
  tenant_domain: '',
  email: '',
  password: '',
})

async function submitLogin() {
  try {
    await auth.login(form)

    const role = auth.role

    if (['owner', 'admin', 'staff'].includes(role || '')) {
      router.push('/admin')
    } else if (role === 'teacher') {
      router.push('/teacher')
    } else if (role === 'student') {
      router.push('/student')
    } else {
      router.push('/login')
    }
  } catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-shell">
      <div class="login-brand">
        <div class="brand-badge">LMS SaaS</div>
        <h1>Run your academy from one dashboard</h1>
        <p>
          Manage teachers, students, subjects, and analytics with a modern multi-tenant LMS.
        </p>

        <div class="feature-list">
          <div class="feature-card">
            <span>Multi-tenant</span>
            <strong>Domain based access</strong>
          </div>
          <div class="feature-card">
            <span>Secure Auth</span>
            <strong>JWT gateway protection</strong>
          </div>
          <div class="feature-card">
            <span>Role-based UI</span>
            <strong>Admin, teacher, student dashboards</strong>
          </div>
        </div>
      </div>

      <div class="login-card">
        <div class="login-card-header">
          <h2>Welcome back</h2>
          <p>Sign in to continue to your dashboard</p>
        </div>

        <form class="login-form" @submit.prevent="submitLogin">
          <div class="form-group">
            <label>Tenant Domain</label>
            <input v-model="form.tenant_domain" type="text" placeholder="guru.localhost" />
          </div>

          <div class="form-group">
            <label>Email</label>
            <input v-model="form.email" type="email" placeholder="teacher@example1.com" />
          </div>

          <div class="form-group">
            <label>Password</label>
            <input v-model="form.password" type="password" placeholder="••••••••" />
          </div>

          <p v-if="auth.error" class="error-message">{{ auth.error }}</p>

          <button type="submit" class="login-btn" :disabled="auth.loading">
            {{ auth.loading ? 'Signing in...' : 'Login' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(99, 102, 241, 0.14), transparent 30%),
    radial-gradient(circle at bottom right, rgba(14, 165, 233, 0.14), transparent 30%),
    #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-shell {
  width: 100%;
  max-width: 1180px;
  display: grid;
  grid-template-columns: 1.15fr 0.85fr;
  gap: 28px;
  align-items: stretch;
}

.login-brand {
  color: white;
  padding: 32px;
  border-radius: 28px;
  background: linear-gradient(145deg, rgba(255,255,255,0.08), rgba(255,255,255,0.03));
  border: 1px solid rgba(255,255,255,0.08);
  backdrop-filter: blur(10px);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.brand-badge {
  display: inline-flex;
  width: fit-content;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(255,255,255,0.12);
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 20px;
}

.login-brand h1 {
  font-size: 46px;
  line-height: 1.08;
  margin: 0 0 16px;
  max-width: 580px;
}

.login-brand p {
  margin: 0;
  color: rgba(255,255,255,0.78);
  font-size: 17px;
  line-height: 1.7;
  max-width: 560px;
}

.feature-list {
  display: grid;
  gap: 14px;
  margin-top: 34px;
}

.feature-card {
  padding: 18px 20px;
  border-radius: 20px;
  background: rgba(255,255,255,0.08);
  border: 1px solid rgba(255,255,255,0.06);
}

.feature-card span {
  display: block;
  font-size: 13px;
  color: rgba(255,255,255,0.68);
  margin-bottom: 6px;
}

.feature-card strong {
  font-size: 16px;
  font-weight: 600;
}

.login-card {
  background: #ffffff;
  border-radius: 28px;
  padding: 34px;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.28);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login-card-header h2 {
  font-size: 30px;
  margin: 0 0 8px;
  color: #0f172a;
}

.login-card-header p {
  margin: 0 0 28px;
  color: #64748b;
}

.login-form {
  display: grid;
  gap: 18px;
}

.form-group {
  display: grid;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: #334155;
}

.form-group input {
  height: 52px;
  border-radius: 14px;
  border: 1px solid #dbe3ef;
  padding: 0 16px;
  font-size: 15px;
  outline: none;
  transition: 0.2s ease;
  background: #f8fafc;
}

.form-group input:focus {
  border-color: #6366f1;
  background: white;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.12);
}

.error-message {
  margin: 0;
  color: #dc2626;
  font-size: 14px;
}

.login-btn {
  height: 54px;
  border: 0;
  border-radius: 14px;
  background: linear-gradient(135deg, #4f46e5, #0ea5e9);
  color: white;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.16s ease, opacity 0.16s ease;
}

.login-btn:hover {
  transform: translateY(-1px);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 960px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .login-brand {
    order: 2;
  }

  .login-card {
    order: 1;
  }

  .login-brand h1 {
    font-size: 34px;
  }
}

@media (max-width: 640px) {
  .login-page {
    padding: 14px;
  }

  .login-card,
  .login-brand {
    padding: 22px;
    border-radius: 22px;
  }

  .login-brand h1 {
    font-size: 28px;
  }
}
</style>