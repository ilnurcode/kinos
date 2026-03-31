<template>
  <footer class="py-3 mt-auto bg-light">
    <div class="container text-center text-muted">
      <small>&copy; {{ year }} Kinos. Все права защищены.</small>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>
