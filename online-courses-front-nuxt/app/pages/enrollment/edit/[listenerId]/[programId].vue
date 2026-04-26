<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const listenerId = computed(() => String(route.params.listenerId || ''))

onMounted(async () => {
  notifications.info('Изменение записей на курс отключено.', 'Запись на курс')
  await router.replace(`/enrollment/details/${listenerId.value}`)
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Переадресация">
      <p>Функция изменения записи на курс отключена. Возвращаю на страницу записей.</p>
    </AppCard>
  </section>
</template>
