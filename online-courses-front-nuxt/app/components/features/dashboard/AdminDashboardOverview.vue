<script setup lang="ts">
import type { AdminDashboardResponse } from '../../../types/dashboard'
import AdminUserRegistrationCard from './AdminUserRegistrationCard.vue'
import AppCard from '../../ui/AppCard.vue'
import AppStatCard from '../../ui/AppStatCard.vue'

defineProps<{
  data: AdminDashboardResponse
}>()
</script>

<template>
  <section class="dashboard-stack">
    <div class="stats-grid">
      <AppStatCard title="База данных" :value="data.admin_stat.database_name || '-'" />
      <AppStatCard title="Размер базы" :value="data.admin_stat.total_size || '-'" />
      <AppStatCard title="Активные соединения" :value="data.admin_stat.active_connection ?? '-'" />
      <AppStatCard title="Коммит транзакций" :value="data.admin_stat.committed_tx ?? '-'" />
      <AppStatCard title="Откат транзакций" :value="data.admin_stat.rolledback_tx ?? '-'" />
      <AppStatCard title="Чтение блоков диска" :value="data.admin_stat.disk_block_read ?? '-'" />
      <AppStatCard title="Буферные попадания" :value="data.admin_stat.buffer_hits ?? '-'" />
    </div>

    <AppCard title="Процессы PostgreSQL">
      <div v-if="!data.pg_stat?.length" class="empty-state">
        Активных процессов сейчас нет.
      </div>

      <div v-else class="process-grid">
        <article v-for="process in data.pg_stat" :key="process.pid" class="process-card">
          <strong>PID: {{ process.pid }}</strong>
          <span>Состояние: {{ process.state || '—' }}</span>
          <span>
            Клиент: {{ process.client_addr || '—' }}<template v-if="process.client_port">:{{ process.client_port }}</template>
          </span>
          <span>
            Начало запроса:
            {{ process.query_start ? new Date(process.query_start).toLocaleString('ru-RU') : '—' }}
          </span>
        </article>
      </div>
    </AppCard>

    <AdminUserRegistrationCard :roles="data.role || []" />
  </section>
</template>

<style scoped>
.dashboard-stack {
  display: grid;
  gap: 1rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.process-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 0.9rem;
}

.process-card {
  display: grid;
  gap: 0.45rem;
  padding: 1rem;
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid rgba(15, 23, 42, 0.07);
  color: #334155;
}

.empty-state {
  color: #64748b;
}

@media (max-width: 900px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
