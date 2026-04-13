<template>
  <MainLayout title="Student Dashboard">
    <div v-if="loading">Loading...</div>

    <div v-else class="subject-grid">
      <div v-for="subject in subjects" :key="subject.id" class="subject-card">
        <h3>{{ subject.name }}</h3>
        <p>{{ subject.grade_name }}</p>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { getStudentSubjects } from '@/api/dashboard'
import type { Subject } from '@/types/dashboard'

const loading = ref(true)
const subjects = ref<Subject[]>([])

onMounted(async () => {
  try {
    subjects.value = await getStudentSubjects()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.subject-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}
.subject-card {
  background: white;
  padding: 20px;
  border-radius: 16px;
}
</style>