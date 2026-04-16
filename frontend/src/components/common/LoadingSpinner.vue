<template>
  <div
    class="loading-spinner"
    :class="{ 'fullscreen': fullscreen }"
  >
    <div
      class="spinner-border text-primary"
      role="status"
    >
      <span class="visually-hidden">{{ text }}</span>
    </div>
    <p
      v-if="text"
      class="mt-2 text-muted"
    >
      {{ text }}
    </p>
  </div>
</template>

<script setup>
defineProps({
  fullscreen: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: 'Загрузка...'
  }
})
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.loading-spinner.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  z-index: 9999;
}
</style>
