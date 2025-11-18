<template>
  <div class="modal fade" tabindex="-1" ref="modalRef">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">

        <div class="modal-header bg-danger text-white">
          <h5 class="modal-title">{{ title }}</h5>
          <button type="button" class="btn-close btn-close-white" @click="close"></button>
        </div>

        <div class="modal-body">
          <p>{{ message }}</p>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="close">Отмена</button>
          <button class="btn btn-danger" @click="confirm">Удалить</button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineExpose } from "vue"
import bootstrap from "bootstrap/dist/js/bootstrap.bundle.min.js"

const props = defineProps({
  title: { type: String, default: "Подтверждение" },
  message: { type: String, default: "Вы уверены?" }
})

const modalRef = ref(null)
let modalInstance = null
let onConfirmCallback = null

const open = (cb) => {
  onConfirmCallback = cb
  modalInstance = new bootstrap.Modal(modalRef.value)
  modalInstance.show()
}

const confirm = () => {
  if (onConfirmCallback) onConfirmCallback()
  modalInstance.hide()
}

const close = () => {
  modalInstance.hide()
}

defineExpose({ open })
</script>
