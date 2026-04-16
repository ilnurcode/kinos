<template>
  <nav aria-label="Пагинация">
    <ul class="pagination justify-content-center">
      <li
        class="page-item"
        :class="{ disabled: currentPage === 1 }"
      >
        <a
          class="page-link"
          href="#"
          @click.prevent="changePage(currentPage - 1)"
        >
          ← Пред.
        </a>
      </li>
      
      <li 
        v-for="page in visiblePages" 
        :key="page"
        class="page-item" 
        :class="{ active: page === currentPage }"
      >
        <a
          class="page-link"
          href="#"
          @click.prevent="changePage(page)"
        >
          {{ page }}
        </a>
      </li>
      
      <li
        class="page-item"
        :class="{ disabled: currentPage === totalPages }"
      >
        <a
          class="page-link"
          href="#"
          @click.prevent="changePage(currentPage + 1)"
        >
          След. →
        </a>
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentPage: {
    type: Number,
    default: 1
  },
  totalPages: {
    type: Number,
    default: 1
  }
})

const emit = defineEmits(['page-change'])

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, props.currentPage - 2)
  const end = Math.min(props.totalPages, props.currentPage + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

function changePage(page) {
  if (page >= 1 && page <= props.totalPages) {
    emit('page-change', page)
  }
}
</script>
