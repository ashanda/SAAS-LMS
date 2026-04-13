<template>
  <MainLayout title="Admin Dashboard">
    <div v-if="loading">Loading...</div>

    <div v-else class="stats-grid">
      <div class="stat-card">
        <h3>Students</h3>
        <p>{{ stats.students }}</p>
      </div>
      <div class="stat-card">
        <h3>Teachers</h3>
        <p>{{ stats.teachers }}</p>
      </div>
      <div class="stat-card">
        <h3>Subjects</h3>
        <p>{{ stats.subjects }}</p>
      </div>
      <div class="stat-card">
        <h3>Enrollments</h3>
        <p>{{ stats.enrollments }}</p>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { getAdminDashboard } from '@/api/dashboard'
import type { AdminDashboardStats } from '@/types/dashboard'

const loading = ref(true)
const stats = ref<AdminDashboardStats>({
  students: 0,
  teachers: 0,
  subjects: 0,
  enrollments: 0,
})

onMounted(async () => {
  try {
    stats.value = await getAdminDashboard()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}
.stat-card {
  background: white;
  padding: 24px;
  border-radius: 16px;
}
.stat-card p {
  font-size: 28px;
  font-weight: bold;
}
</style>