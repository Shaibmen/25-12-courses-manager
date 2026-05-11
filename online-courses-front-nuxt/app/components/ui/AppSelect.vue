<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})

const attrs = useAttrs()

withDefaults(defineProps<{
  label: string
  modelValue: string | number
  placeholder?: string
  help?: string
  error?: string
  disabled?: boolean
}>(), {
  placeholder: '',
  help: '',
  error: '',
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<template>
  <label class="app-select">
    <span class="app-select__label">{{ label }}</span>

    <span class="app-select__field" :class="{ 'app-select__field--placeholder': !modelValue, 'app-select__field--error': Boolean(error) }">
      <select
        v-bind="attrs"
        :value="modelValue ?? ''"
        :disabled="disabled"
        class="app-select__control"
        @change="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
      >
        <option v-if="placeholder" value="">{{ placeholder }}</option>
        <slot />
      </select>

      <span class="app-select__icon" aria-hidden="true">
        <svg viewBox="0 0 20 20" fill="none">
          <path
            d="M5 7.5L10 12.5L15 7.5"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.8"
          />
        </svg>
      </span>
    </span>

    <span v-if="error" class="app-select__message app-select__message--error">{{ error }}</span>
    <span v-else-if="help" class="app-select__message">{{ help }}</span>
  </label>
</template>

<style scoped>
.app-select {
  display: grid;
  gap: 0.45rem;
}

.app-select__label {
  font-size: 0.92rem;
  font-weight: 600;
  color: #0f172a;
}

.app-select__field {
  position: relative;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.97) 0%, rgba(248, 250, 252, 0.94) 100%);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.7);
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.app-select__field:focus-within {
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}

.app-select__field--placeholder .app-select__control {
  color: #64748b;
}

.app-select__field--error {
  border-color: rgba(220, 38, 38, 0.35);
  box-shadow: 0 0 0 4px rgba(248, 113, 113, 0.08);
}

.app-select__control {
  appearance: none;
  -webkit-appearance: none;
  width: 100%;
  min-height: 3rem;
  padding: 0.85rem 3rem 0.85rem 1rem;
  border: none;
  border-radius: 1rem;
  background: transparent;
  color: #0f172a;
  font: inherit;
  outline: none;
  cursor: pointer;
}

.app-select__control:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.app-select__icon {
  position: absolute;
  top: 50%;
  right: 1rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1rem;
  height: 1rem;
  color: #475569;
  pointer-events: none;
  transform: translateY(-50%);
}

.app-select__icon svg {
  width: 100%;
  height: 100%;
}

.app-select__message {
  font-size: 0.82rem;
  color: #64748b;
}

.app-select__message--error {
  color: #b91c1c;
}
</style>
