<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    to?: string
    type?: 'button' | 'submit' | 'reset'
    variant?: 'primary' | 'secondary' | 'ghost'
    disabled?: boolean
    block?: boolean
  }>(),
  {
    type: 'button',
    variant: 'primary',
    disabled: false,
    block: false
  }
)

const componentTag = computed(() => (props.to ? resolveComponent('NuxtLink') : 'button'))
</script>

<template>
  <component
    :is="componentTag"
    class="app-button"
    :class="[`app-button--${variant}`, { 'app-button--block': block }]"
    :to="to"
    :type="to ? undefined : type"
    :disabled="to ? undefined : disabled"
    :aria-disabled="disabled ? 'true' : undefined"
  >
    <span class="app-button__glow" />
    <span class="app-button__label">
      <slot />
    </span>
  </component>
</template>

<style scoped>
.app-button {
  position: relative;
  isolation: isolate;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  min-height: 2.9rem;
  padding: 0.78rem 1.05rem;
  border: none;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.96rem;
  text-decoration: none;
  cursor: pointer;
  overflow: hidden;
  transition:
    transform 180ms ease,
    box-shadow 220ms ease,
    background-color 220ms ease,
    color 220ms ease,
    border-color 220ms ease;
}

.app-button::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, transparent 20%, rgba(255, 255, 255, 0.35), transparent 80%);
  transform: translateX(-120%);
  transition: transform 480ms ease;
  z-index: 0;
}

.app-button:hover {
  transform: translateY(-2px) scale(1.01);
}

.app-button:hover::before {
  transform: translateX(120%);
}

.app-button:active {
  transform: translateY(1px) scale(0.99);
}

.app-button:focus-visible {
  outline: 3px solid rgba(59, 130, 246, 0.24);
  outline-offset: 3px;
}

.app-button[aria-disabled='true'] {
  cursor: not-allowed;
  opacity: 0.6;
  pointer-events: none;
}

.app-button--block {
  width: 100%;
}

.app-button__glow {
  position: absolute;
  inset: auto auto -1.1rem 50%;
  width: 60%;
  height: 1.1rem;
  transform: translateX(-50%);
  filter: blur(18px);
  opacity: 0.45;
  z-index: 0;
}

.app-button__label {
  position: relative;
  z-index: 1;
}

.app-button--primary {
  color: #eff6ff;
  background: linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%);
  box-shadow: 0 16px 34px rgba(29, 78, 216, 0.25);
}

.app-button--primary .app-button__glow {
  background: rgba(59, 130, 246, 0.75);
}

.app-button--secondary {
  color: #0f172a;
  background: linear-gradient(135deg, #f8fafc 0%, #dbeafe 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.1);
}

.app-button--secondary .app-button__glow {
  background: rgba(148, 163, 184, 0.55);
}

.app-button--ghost {
  color: #0f172a;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.08);
}

.app-button--ghost .app-button__glow {
  background: rgba(148, 163, 184, 0.45);
}

.app-button.router-link-active.app-button--ghost {
  color: #eff6ff;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}
</style>
