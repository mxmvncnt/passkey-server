<script setup lang="ts">
import { ref } from 'vue'
import { startRegistration } from '@simplewebauthn/browser'
import { API_BASE, JSON_HEADERS } from '../config'

const name = ref('')
const busy = ref(false)
const error = ref<string | null>(null)

async function register() {
  error.value = null
  busy.value = true
  try {
    const opts = await fetch(`${API_BASE}/passkey/register/begin`, {
      method: 'POST',
      headers: JSON_HEADERS,
      body: JSON.stringify({ email: name.value }),
    }).then((r) => r.json())
    console.log('register begin', opts)

    const cred = await startRegistration({ optionsJSON: opts })
    console.log('register credential', cred)

    const finishUrl = new URL(`${API_BASE}/passkey/register/finish`)
    finishUrl.searchParams.append('user_id', opts.user.id)

    const result = await fetch(finishUrl, {
      method: 'POST',
      headers: JSON_HEADERS,
      body: JSON.stringify(cred),
    })
    console.log('register finish', result)
    if (!result.ok) throw new Error(`Register finish failed: ${result.status}`)
    alert('c bon tu px te connecter mtn')
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
    console.error(e)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <section>
    <h2>Register</h2>
    <input
      v-model="name"
      type="text"
      placeholder="name"
      :disabled="busy"
    />
    <button type="button" :disabled="busy" @click="register">Register</button>
    <p v-if="error" >{{ error }}</p>
  </section>
</template>
