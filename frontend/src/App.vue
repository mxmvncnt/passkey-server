<script setup lang="ts">
import RegisterForm from './components/RegisterForm.vue'
import LoginForm from './components/LoginForm.vue'
</script>

<template>
  <main style="max-width: 22rem; padding: 1rem">
    <h1 style="font-size: 1rem; font-weight: 600; margin: 0 0 1rem">Passkey tester</h1>
    <RegisterForm />
    <hr style="margin: 1.25rem 0; border: none; border-top: 1px solid #ccc" />
    <LoginForm />
  </main>
</template>
