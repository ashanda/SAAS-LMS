<template>
  <MainLayout title="Teacher Dashboard">
    <div class="dashboard-grid">
      <section class="hero-card">
        <div>
          <p class="hero-label">Overview</p>
          <h2>Welcome back, teacher</h2>
          <p class="hero-text">
            View your assigned subjects and quickly access the students for each class.
          </p>
        </div>

        <div class="hero-stats">
          <div class="mini-stat">
            <span>Total Subjects</span>
            <strong>{{ subjects.length }}</strong>
          </div>
        </div>
      </section>

      <section class="content-card full-width">
        <div class="section-head">
          <div>
            <h3>Your Subjects</h3>
            <p>All subjects assigned to your account</p>
          </div>
        </div>

        <div v-if="loading" class="empty-state">Loading subjects...</div>

        <div v-else-if="subjects.length === 0" class="empty-state">
          No subjects assigned yet.
        </div>

        <div v-else class="subject-grid">
          <div v-for="subject in subjects" :key="subject.id" class="subject-card">
            <div class="subject-top">
              <div>
                <h4>{{ subject.name }}</h4>
                <p>{{ subject.grade_name || 'No grade assigned' }}</p>
              </div>

              <button class="view-btn" @click="loadStudents(subject.id)">
                View Students
              </button>
            </div>

            <transition name="fade">
              <div
                v-if="selectedSubjectId === subject.id"
                class="students-panel"
              >
                <div class="students-head">
                  <h5>Students</h5>
                  <span>{{ students.length }}</span>
                </div>

                <div v-if="studentsLoading" class="empty-state small">
                  Loading students...
                </div>

                <div v-else-if="students.length === 0" class="empty-state small">
                  No students found for this subject.
                </div>

                <ul v-else class="student-list">
                  <li v-for="student in students" :key="student.id" class="student-item">
                    <div class="avatar">
                      {{ student.name?.charAt(0)?.toUpperCase() || 'S' }}
                    </div>
                    <div>
                      <strong>{{ student.name }}</strong>
                      <p>{{ student.email || 'No email available' }}</p>
                    </div>
                  </li>
                </ul>
              </div>
            </transition>
          </div>
        </div>
      </section>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import { getTeacherStudents, getTeacherSubjects } from '@/api/dashboard'
import type { Student, Subject } from '@/types/dashboard'

const loading = ref(true)
const studentsLoading = ref(false)
const subjects = ref<Subject[]>([])
const students = ref<Student[]>([])
const selectedSubjectId = ref<string | null>(null)

onMounted(async () => {
  try {
    subjects.value = await getTeacherSubjects()
  } catch (error) {
    console.error('Teacher subjects error:', error)
  } finally {
    loading.value = false
  }
})

async function loadStudents(subjectId: string) {
  try {
    studentsLoading.value = true
    selectedSubjectId.value = subjectId
    students.value = await getTeacherStudents(subjectId)
  } catch (error) {
    console.error('Teacher students error:', error)
  } finally {
    studentsLoading.value = false
  }
}
</script>

<style scoped>
.dashboard-grid {
  display: grid;
  gap: 22px;
}

.full-width {
  grid-column: 1 / -1;
}

.hero-card {
  background: linear-gradient(135deg, #4f46e5, #0ea5e9);
  color: white;
  border-radius: 24px;
  padding: 28px;
  display: flex;
  justify-content: space-between;
  gap: 18px;
  align-items: center;
  flex-wrap: wrap;
}

.hero-label {
  margin: 0 0 10px;
  font-size: 13px;
  color: rgba(255,255,255,0.74);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.hero-card h2 {
  margin: 0 0 10px;
  font-size: 30px;
}

.hero-text {
  margin: 0;
  max-width: 640px;
  color: rgba(255,255,255,0.85);
  line-height: 1.7;
}

.hero-stats {
  display: flex;
  gap: 14px;
}

.mini-stat {
  background: rgba(255,255,255,0.14);
  border: 1px solid rgba(255,255,255,0.14);
  padding: 16px 18px;
  border-radius: 18px;
  min-width: 150px;
}

.mini-stat span {
  display: block;
  font-size: 13px;
  color: rgba(255,255,255,0.76);
  margin-bottom: 6px;
}

.mini-stat strong {
  font-size: 28px;
}

.content-card {
  background: white;
  border-radius: 24px;
  padding: 24px;
  border: 1px solid #e2e8f0;
}

.section-head {
  margin-bottom: 20px;
}

.section-head h3 {
  margin: 0;
  font-size: 22px;
  color: #0f172a;
}

.section-head p {
  margin: 6px 0 0;
  color: #64748b;
}

.subject-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.subject-card {
  border: 1px solid #e2e8f0;
  border-radius: 20px;
  padding: 18px;
  background: #fff;
}

.subject-top {
  display: flex;
  justify-content: space-between;
  align-items: start;
  gap: 14px;
}

.subject-card h4 {
  margin: 0;
  font-size: 18px;
  color: #0f172a;
}

.subject-card p {
  margin: 6px 0 0;
  color: #64748b;
}

.view-btn {
  border: 0;
  border-radius: 12px;
  background: #0f172a;
  color: white;
  padding: 11px 14px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
}

.students-panel {
  margin-top: 18px;
  padding-top: 18px;
  border-top: 1px solid #e2e8f0;
}

.students-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.students-head h5 {
  margin: 0;
  font-size: 16px;
  color: #0f172a;
}

.students-head span {
  background: #eef2ff;
  color: #4338ca;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 700;
}

.student-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 12px;
}

.student-item {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px;
  background: #f8fafc;
  border-radius: 16px;
}

.avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #0ea5e9);
  color: white;
  display: grid;
  place-items: center;
  font-weight: 700;
  flex-shrink: 0;
}

.student-item strong {
  display: block;
  color: #0f172a;
}

.student-item p {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 14px;
}

.empty-state {
  min-height: 160px;
  display: grid;
  place-items: center;
  color: #64748b;
  border: 1px dashed #cbd5e1;
  border-radius: 18px;
  background: #f8fafc;
  text-align: center;
  padding: 20px;
}

.empty-state.small {
  min-height: 90px;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

@media (max-width: 960px) {
  .subject-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .hero-card,
  .content-card {
    padding: 18px;
    border-radius: 20px;
  }

  .hero-card h2 {
    font-size: 24px;
  }

  .subject-top {
    flex-direction: column;
    align-items: stretch;
  }

  .view-btn {
    width: 100%;
  }
}
</style>