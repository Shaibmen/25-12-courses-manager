<script setup lang="ts">
import GraphicsDashboardView from '../components/features/graphics/GraphicsDashboardView.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'

const { data, status, error, refresh } = await useAsyncData(
  'graphics-dashboard',
  async () => await getGraphicsDashboard()
)
</script>

<template>
  <section class="stack">
    <AppCard v-if="status === 'pending'" title="Загрузка аналитики">
      <p>Загрузка данных для аналитического модуля и построения графиков.</p>
    </AppCard>

    <AppCard v-else-if="error" title="Не удалось загрузить аналитику">
      <p>{{ error.message }}</p>
      <AppButton class="graphics-page__retry" variant="secondary" @click="refresh">
        Повторить запрос
      </AppButton>
    </AppCard>

    <GraphicsDashboardView
      v-else-if="data"
      :data="data"
    />
  </section>
</template>

<style scoped>
.graphics-page__retry {
  margin-top: 1rem;
}
</style>
