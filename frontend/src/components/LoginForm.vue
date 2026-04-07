<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { startAuthentication } from '@simplewebauthn/browser'
import type { PublicKeyCredentialRequestOptionsJSON } from '@simplewebauthn/browser'
import { API_BASE, JSON_HEADERS } from '../config'
import type { ApiUser } from '../views/UserView.vue'

const router = useRouter()
const busy = ref(false)
const error = ref<string | null>(null)

interface LoginBeginResponse {
  session_id: string
  options: { publicKey: PublicKeyCredentialRequestOptionsJSON }
}

async function login() {
  error.value = null
  busy.value = true
  try {
    const optionsJSON = (await fetch(`${API_BASE}/passkey/login/begin`, {
      method: 'POST',
      headers: JSON_HEADERS,
    }).then((r) => r.json())) as LoginBeginResponse
    console.log('login begin', optionsJSON)

    const cred = await startAuthentication({
      optionsJSON: optionsJSON.options.publicKey,
    })
    console.log('login credential', cred)

    const finishUrl = new URL(`${API_BASE}/passkey/login/finish`)
    finishUrl.searchParams.append('session_id', optionsJSON.session_id)

    const result = (await fetch(finishUrl, {
      method: 'POST',
      headers: JSON_HEADERS,
      body: JSON.stringify(cred),
    }).then((r) => r.json())) as ApiUser

    console.log('login finish', result)
    await router.push({ path: `user/${result.ID}` })
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
    <h2>Login</h2>
    <button type="button" :disabled="busy" @click="login">Login</button>
    <p v-if="error" >{{ error }}</p>
  </section>
</template>
