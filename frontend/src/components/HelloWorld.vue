<script setup lang="ts">
defineProps<{ msg: string }>()
</script>

<template>
  <div class="text-center space-y-2">
    <h2 class="text-lg font-medium">{{ msg }}</h2>
  </div>
</template>
