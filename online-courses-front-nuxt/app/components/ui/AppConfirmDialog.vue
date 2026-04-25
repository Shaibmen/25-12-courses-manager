<script setup lang="ts">
import AppButton from './AppButton.vue'

defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel?: string
  cancelLabel?: string
  loading?: boolean
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="dialog-backdrop" @click.self="$emit('cancel')">
      <section class="dialog-card" role="dialog" aria-modal="true" :aria-label="title">
        <header class="dialog-card__header">
          <h3>{{ title }}</h3>
        </header>

        <div class="dialog-card__body">
          <p>{{ message }}</p>
        </div>

        <footer class="dialog-card__actions">
          <AppButton variant="ghost" :disabled="loading" @click="$emit('cancel')">
            {{ cancelLabel || 'Отмена' }}
          </AppButton>
          <AppButton :disabled="loading" @click="$emit('confirm')">
            {{ loading ? 'Выполняем...' : (confirmLabel || 'Подтвердить') }}
          </AppButton>
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<style scoped>
.dialog-backdrop {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  background: rgba(15, 23, 42, 0.42);
  backdrop-filter: blur(10px);
}

.dialog-card {
  width: min(32rem, 100%);
  padding: 1.35rem;
  border-radius: 1.35rem;
  border: 1px solid rgba(148, 163, 184, 0.2);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(248, 250, 252, 0.96));
  box-shadow: 0 30px 90px rgba(15, 23, 42, 0.24);
  animation: dialog-enter 180ms ease;
}

.dialog-card__header h3,
.dialog-card__body p {
  margin: 0;
}

.dialog-card__body {
  margin-top: 0.85rem;
  color: #334155;
  line-height: 1.5;
}

.dialog-card__actions {
  margin-top: 1.25rem;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  flex-wrap: wrap;
}

@keyframes dialog-enter {
  from {
    opacity: 0;
    transform: translateY(0.5rem) scale(0.98);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (max-width: 720px) {
  .dialog-card__actions {
    flex-direction: column-reverse;
  }
}
</style>
