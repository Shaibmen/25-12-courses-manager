<script setup lang="ts">
import type { GroupItem } from '../../../types/group'

const props = withDefaults(defineProps<{
  modelValue: string
  options: GroupItem[]
  label?: string
  placeholder?: string
  help?: string
  error?: string
  loading?: boolean
  disabled?: boolean
}>(), {
  label: 'Группа',
  placeholder: 'Начните вводить группу',
  help: '',
  error: '',
  loading: false,
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  select: [group: GroupItem]
}>()

const rootRef = ref<HTMLElement | null>(null)
const isOpen = ref(false)

const hasStateMessage = computed(() => Boolean(props.loading || props.modelValue.trim()))

const handleDocumentClick = (event: MouseEvent) => {
  if (!rootRef.value) {
    return
  }

  if (!rootRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

const handleInput = (event: Event) => {
  isOpen.value = true
  emit('update:modelValue', (event.target as HTMLInputElement).value)
}

const handleSelect = (group: GroupItem) => {
  emit('update:modelValue', group.name_group)
  emit('select', group)
  isOpen.value = false
}

onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick)
})
</script>

<template>
  <label ref="rootRef" class="group-search">
    <span class="group-search__label">{{ label }}</span>

    <div class="group-search__field">
      <input
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="{ 'group-search__control--error': Boolean(error) }"
        class="group-search__control"
        autocomplete="off"
        @focus="isOpen = true"
        @input="handleInput"
      >

      <div v-if="isOpen" class="group-search__dropdown">
        <div v-if="loading" class="group-search__state">
          Ищу подходящие группы...
        </div>

        <template v-else-if="options.length">
          <button
            v-for="group in options"
            :key="group.group"
            type="button"
            class="group-search__option"
            @click="handleSelect(group)"
          >
            <span class="group-search__option-name">{{ group.name_group }}</span>
          </button>
        </template>

        <div v-else-if="hasStateMessage" class="group-search__state">
          Похожие группы не найдены.
        </div>
      </div>
    </div>

    <span v-if="error" class="group-search__message group-search__message--error">{{ error }}</span>
    <span v-else-if="help" class="group-search__message">{{ help }}</span>
  </label>
</template>

<style scoped>
.group-search {
  display: grid;
  gap: 0.45rem;
}

.group-search__label {
  font-size: 0.92rem;
  font-weight: 600;
  color: #0f172a;
}

.group-search__field {
  position: relative;
}

.group-search__control {
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

.group-search__control:focus {
  outline: none;
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}

.group-search__control--error {
  border-color: rgba(220, 38, 38, 0.45);
  box-shadow: 0 0 0 4px rgba(248, 113, 113, 0.08);
}

.group-search__control:disabled {
  cursor: not-allowed;
  opacity: 0.6;
  background: rgba(226, 232, 240, 0.65);
}

.group-search__dropdown {
  position: absolute;
  top: calc(100% + 0.45rem);
  left: 0;
  right: 0;
  z-index: 20;
  display: grid;
  gap: 0.35rem;
  max-height: 16rem;
  padding: 0.45rem;
  overflow-y: auto;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.14);
}

.group-search__option {
  display: grid;
  gap: 0.2rem;
  padding: 0.8rem 0.9rem;
  border: none;
  border-radius: 0.85rem;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: background 160ms ease;
}

.group-search__option:hover {
  background: rgba(219, 234, 254, 0.55);
}

.group-search__option-name {
  font-weight: 600;
  color: #0f172a;
}

.group-search__state,
.group-search__message {
  font-size: 0.82rem;
  color: #64748b;
}

.group-search__message--error {
  color: #b91c1c;
}
</style>
