<script setup lang="ts">
const notifications = useNotifications()

const titleByVariant = {
  success: 'Успешно',
  error: 'Нужно внимание',
  info: 'Информация'
} as const
</script>

<template>
  <Teleport to="body">
    <div v-if="notifications.notifications.value.length" class="notifications" aria-live="polite" aria-atomic="true">
      <article
        v-for="item in notifications.notifications.value"
        :key="item.id"
        class="notifications__item"
        :class="`notifications__item--${item.variant}`"
      >
        <div class="notifications__content">
          <strong class="notifications__title">
            {{ item.title || titleByVariant[item.variant] }}
          </strong>
          <p class="notifications__message">{{ item.message }}</p>
        </div>

        <button
          type="button"
          class="notifications__close"
          aria-label="Закрыть уведомление"
          @click="notifications.remove(item.id)"
        >
          ×
        </button>
      </article>
    </div>
  </Teleport>
</template>

<style scoped>
.notifications {
  position: fixed;
  top: 1.25rem;
  right: 1.25rem;
  z-index: 70;
  display: grid;
  gap: 0.85rem;
  width: min(26rem, calc(100vw - 2rem));
}

.notifications__item {
  display: flex;
  justify-content: space-between;
  gap: 0.9rem;
  padding: 1rem 1rem 1rem 1.1rem;
  border-radius: 1.15rem;
  border: 1px solid rgba(15, 23, 42, 0.1);
  box-shadow:
    0 24px 54px rgba(15, 23, 42, 0.18),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(20px);
  background: rgba(255, 255, 255, 0.94);
  animation: notification-enter 220ms ease;
}

.notifications__item--success {
  border-color: rgba(22, 163, 74, 0.2);
  background:
    linear-gradient(180deg, rgba(240, 253, 244, 0.96), rgba(255, 255, 255, 0.96));
}

.notifications__item--error {
  border-color: rgba(220, 38, 38, 0.2);
  background:
    linear-gradient(180deg, rgba(254, 242, 242, 0.96), rgba(255, 255, 255, 0.96));
}

.notifications__item--info {
  border-color: rgba(14, 116, 144, 0.18);
  background:
    linear-gradient(180deg, rgba(240, 249, 255, 0.96), rgba(255, 255, 255, 0.96));
}

.notifications__content {
  display: grid;
  gap: 0.3rem;
  min-width: 0;
}

.notifications__title {
  color: #0f172a;
  font-size: 0.95rem;
}

.notifications__message {
  margin: 0;
  color: #334155;
  line-height: 1.4;
}

.notifications__close {
  border: none;
  width: 2rem;
  height: 2rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  color: #64748b;
  font-size: 1.1rem;
  line-height: 1;
  cursor: pointer;
  padding: 0;
}

.notifications__close:hover {
  color: #0f172a;
  background: rgba(255, 255, 255, 0.96);
}

@keyframes notification-enter {
  from {
    opacity: 0;
    transform: translateY(-0.45rem) scale(0.98);
  }

  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (max-width: 720px) {
  .notifications {
    top: auto;
    right: 0.75rem;
    bottom: 0.75rem;
    left: 0.75rem;
    width: auto;
  }
}
</style>
