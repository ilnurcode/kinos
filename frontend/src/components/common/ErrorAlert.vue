<template>
  <div class="alert alert-danger alert-dismissible fade show" role="alert">
    {{ message }}
    <button type="button" class="btn-close" @click="$emit('dismiss')"></button>
  </div>
</template>

<script setup>
defineProps({
  message: {
    type: String,
    required: true
  }
})

defineEmits(['dismiss'])
</script>
