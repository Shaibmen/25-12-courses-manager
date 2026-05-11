<script setup lang="ts">
const props = withDefaults(defineProps<{
  modelValue: boolean
  label: string
  disabled?: boolean
  type?: 'checkbox' | 'radio'
  name?: string
}>(), {
  disabled: false,
  type: 'checkbox',
  name: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()
</script>

<template>
  <label class="app-checkbox" :class="{ 'app-checkbox--disabled': disabled }">
    <input
      :checked="modelValue"
      :disabled="disabled"
      :type="type"
      :name="name"
      class="app-checkbox__input"
      @change="emit('update:modelValue', ($event.target as HTMLInputElement).checked)"
    >
    <span class="app-checkbox__control" aria-hidden="true">
      <svg viewBox="0 0 16 16" fill="none">
        <path d="M3.5 8.2L6.6 11.3L12.5 4.9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
      </svg>
    </span>
    <span class="app-checkbox__label">{{ label }}</span>
  </label>
</template>

<style scoped>
.app-checkbox {
  display: inline-flex;
  align-items: center;
  gap: 0.8rem;
  cursor: pointer;
  color: #0f172a;
}

.app-checkbox--disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.app-checkbox__input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.app-checkbox__control {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.35rem;
  height: 1.35rem;
  flex: 0 0 1.35rem;
  border-radius: 0.45rem;
  border: 1px solid rgba(15, 23, 42, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(241, 245, 249, 0.92) 100%);
  box-shadow: 0 10px 22px rgba(15, 23, 42, 0.08);
  color: transparent;
  transition:
    border-color 180ms ease,
    background-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.app-checkbox__control svg {
  width: 0.9rem;
  height: 0.9rem;
}

.app-checkbox__label {
  line-height: 1.4;
}

.app-checkbox__input:checked + .app-checkbox__control {
  border-color: rgba(37, 99, 235, 0.4);
  background: linear-gradient(135deg, #0f172a 0%, #2563eb 100%);
  color: #eff6ff;
}

.app-checkbox__input:focus-visible + .app-checkbox__control {
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}
</style>
