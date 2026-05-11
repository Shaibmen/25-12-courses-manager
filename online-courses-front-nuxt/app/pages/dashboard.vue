<script setup lang="ts">
import type { AdminDashboardResponse, UserDashboardResponse } from '../types/dashboard'
import AdminDashboardOverview from '../components/features/dashboard/AdminDashboardOverview.vue'
import WorkerDashboardOverview from '../components/features/dashboard/WorkerDashboardOverview.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'

const auth = useAuth()

const { data, status, error, refresh } = await useAsyncData(
  () => `dashboard:${auth.state.value.role || 'guest'}`,
  async () => {
    const role = auth.state.value.role

    if (!role) {
      throw new Error('Не удалось определить роль пользователя')
    }

    return await getDashboardByRole(role)
  }
)

const userDashboard = computed<UserDashboardResponse | null>(() =>
  data.value?.kind === 'user' ? data.value.data : null
)

const adminDashboard = computed<AdminDashboardResponse | null>(() =>
  data.value?.kind === 'admin' ? data.value.data : null
)
</script>

<template>
  <section class="stack">
    <AppCard v-if="status === 'pending'" title="Загрузка данных">
      <p>Подтягиваю актуальные данные для вашей роли.</p>
    </AppCard>

    <AppCard v-else-if="error" title="Не удалось загрузить данные">
      <p>{{ error.message }}</p>
      <AppButton class="dashboard__retry" variant="secondary" @click="refresh">
        Повторить запрос
      </AppButton>
    </AppCard>

    <AdminDashboardOverview
      v-else-if="adminDashboard"
      :data="adminDashboard"
    />

    <template v-else-if="userDashboard">
      <WorkerDashboardOverview
        :data="userDashboard"
      />
    </template>
  </section>
</template>

<style scoped>
.dashboard__retry {
  margin-top: 1rem;
}
</style>
