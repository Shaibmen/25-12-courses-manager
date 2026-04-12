<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})

const attrs = useAttrs()

withDefaults(
  defineProps<{
    label: string
    modelValue: string | number
    type?: string
    placeholder?: string
    datePlaceholder?: string
    autocomplete?: string
    help?: string
    error?: string
    disabled?: boolean
  }>(),
  {
    type: 'text',
    placeholder: '',
    datePlaceholder: '',
    autocomplete: 'off',
    help: '',
    error: '',
    disabled: false
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<template>
  <label class="app-input">
    <span class="app-input__label">{{ label }}</span>

    <span class="app-input__field">
      <input
        v-bind="attrs"
        :type="type"
        :value="modelValue ?? ''"
        :placeholder="type === 'date' ? undefined : placeholder"
        :autocomplete="autocomplete"
        :disabled="disabled"
        :class="{ 'app-input__control--error': Boolean(error) }"
        class="app-input__control"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      >

      <span
        v-if="type === 'date' && datePlaceholder && !modelValue"
        class="app-input__date-placeholder"
      >
        {{ datePlaceholder }}
      </span>
    </span>

    <span v-if="error" class="app-input__message app-input__message--error">
      {{ error }}
    </span>
    <span v-else-if="help" class="app-input__message">
      {{ help }}
    </span>
  </label>
</template>

<style scoped>
.app-input {
  display: grid;
  gap: 0.45rem;
}

.app-input__label {
  font-size: 0.92rem;
  font-weight: 600;
  color: #0f172a;
}

.app-input__field {
  position: relative;
  display: block;
}

.app-input__control {
  width: 100%;
  min-height: 3rem;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.92);
  color: #0f172a;
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.app-input__control:focus {
  outline: none;
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}

.app-input__control--error {
  border-color: rgba(220, 38, 38, 0.45);
  box-shadow: 0 0 0 4px rgba(248, 113, 113, 0.08);
}

.app-input__control:disabled {
  cursor: not-allowed;
  opacity: 0.6;
  background: rgba(226, 232, 240, 0.65);
}

.app-input__date-placeholder {
  position: absolute;
  inset: 50% auto auto 1rem;
  color: #94a3b8;
  pointer-events: none;
  transform: translateY(-50%);
}

.app-input__message {
  font-size: 0.82rem;
  color: #64748b;
}

.app-input__message--error {
  color: #b91c1c;
}
</style>
