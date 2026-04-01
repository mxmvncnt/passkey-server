<template>
  <article>
    <h2 style="font-size: 0.875rem; margin: 0 0 0.75rem">User</h2>
    <dl style="margin: 0; display: grid; gap: 0.35rem 1rem; grid-template-columns: auto 1fr; font-size: 0.875rem">
      <dt style="opacity: 0.75">ID</dt>
      <dd style="margin: 0; word-break: break-all">{{ user.ID }}</dd>
      <dt style="opacity: 0.75">Name</dt>
      <dd style="margin: 0">{{ user.Name ?? '—' }}</dd>
      <dt style="opacity: 0.75">Display name</dt>
      <dd style="margin: 0">{{ user.DisplayName ?? '—' }}</dd>
    </dl>
  </article>
</template>

<script setup lang="ts">

import {onMounted, ref} from "vue";
import {API_BASE, JSON_HEADERS} from "../config.ts";

export interface ApiUser {
  ID: string
  Name: string | null
  DisplayName: string | null
}

const user = ref()
const userId = ref(localStorage.getItem("userId"))

onMounted(async () => {
  user.value = await fetch(`${API_BASE}/users/${userId}`, {
    method: 'POST',
    headers: JSON_HEADERS,
  }).then((r) => r.json()) as ApiUser
  console.log(user.value)
})

</script>