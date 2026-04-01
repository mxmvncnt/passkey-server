<template>
  <main style="max-width: 28rem; padding: 1rem">
    <p style="margin: 0 0 1rem">
      <RouterLink to="/" style="color: inherit">← Home</RouterLink>
    </p>
    <UserInfo />
  </main>
</template>

<script setup lang="ts">
import UserInfo from "../components/UserInfo.vue";
</script>