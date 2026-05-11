<script setup lang="ts">
import type { GroupItem } from '../../../types/group'
import type { EnrollmentProgramDetails, ListenerDetailsResponse } from '../../../types/listener'
import { formatApiDate, normalizeApiDate } from '../../../utils/date'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../../components/ui/AppConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const listenerId = computed(() => String(route.params.listenerId || ''))

const listener = ref<ListenerDetailsResponse | null>(null)
const enrollments = ref<EnrollmentProgramDetails[]>([])
const groups = ref<GroupItem[]>([])
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const enrollmentToDelete = ref<EnrollmentProgramDetails | null>(null)

const groupNameById = computed(() =>
  new Map(groups.value.map((group) => [group.group, group.name_group]))
)

const resolveGroupName = (id: string) => groupNameById.value.get(id) || id || '—'

const load = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const [listenerData, enrollmentData, groupData] = await Promise.all([
      getListenerDetails(listenerId.value),
      getListenerEnrollments(listenerId.value),
      getGroups(1, '')
    ])

    listener.value = listenerData
    enrollments.value = enrollmentData
    groups.value = groupData
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить записи слушателя'
  } finally {
    loading.value = false
  }
}

const removeEnrollmentRecord = async () => {
  if (!enrollmentToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteEnrollment(listenerId.value, enrollmentToDelete.value.id_program_education)
    enrollmentToDelete.value = null
    await load()
    notifications.success('Запись удалена.', 'Запись на курс')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить запись',
      'Запись на курс'
    )
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Записи слушателя">
      <div class="summary-header">
        <div v-if="listener" class="summary-text">
          <strong>{{ listener.listener.second_name }} {{ listener.listener.first_name }} {{ listener.listener.middle_name || '' }}</strong>
          <span>{{ listener.listener.contact_phone || '—' }} · {{ listener.listener.email || '—' }}</span>
        </div>

        <div class="summary-actions">
          <AppButton v-if="listener" variant="secondary" @click="router.push(`/listeners/${listenerId}`)">
            Профиль слушателя
          </AppButton>
          <AppButton @click="router.push(`/enrollment/create/${listenerId}`)">Добавить запись</AppButton>
          <AppButton variant="ghost" to="/enrollments">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю записи слушателя.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="load">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Курсы слушателя">
      <div v-if="!enrollments.length" class="empty-state">
        Курсы у слушателя пока не найдены.
      </div>

      <div v-else class="course-list">
        <article
          v-for="course in enrollments"
          :key="course.id_program_education"
          class="course-card"
        >
          <div class="course-card__main">
            <strong>{{ course.name_prof_education }}</strong>
            <span>{{ formatApiDate(course.start_date) }} - {{ formatApiDate(course.end_date) }}</span>
            <span>Цена: {{ new Intl.NumberFormat('ru-RU').format(course.price) }} ₽</span>
            <span>Группа: {{ resolveGroupName(course.id_group) }}</span>
            <span>Тип: {{ course.type_of_retraining || '—' }}</span>
          </div>

          <div class="course-card__actions">
            <AppButton variant="ghost" @click="enrollmentToDelete = course">
              Удалить
            </AppButton>
          </div>
        </article>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(enrollmentToDelete)"
      title="Удаление записи"
      :message="enrollmentToDelete ? `Удалить запись на программу «${enrollmentToDelete.name_prof_education}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="enrollmentToDelete = null"
      @confirm="removeEnrollmentRecord"
    />
  </section>
</template>

<style scoped>
.summary-header,
.summary-actions,
.course-card__actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.summary-header {
  justify-content: space-between;
  align-items: center;
}

.summary-text {
  display: grid;
  gap: 0.25rem;
  color: #334155;
}

.course-list {
  display: grid;
  gap: 0.9rem;
}

.course-card {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  padding: 1rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.course-card__main {
  display: grid;
  gap: 0.3rem;
}

.empty-state {
  color: #64748b;
}

@media (max-width: 900px) {
  .summary-header,
  .summary-actions,
  .course-card,
  .course-card__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
