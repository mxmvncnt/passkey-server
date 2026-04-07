<template>
  <main style="max-width: 28rem; padding: 1rem">
    <p style="margin: 0 0 1rem">
      <RouterLink to="/" style="color: inherit">← Home</RouterLink>
    </p>
    <h2>Welcome back, {{ user?.Name }}</h2>
    <button @click="addPasskey">Register new passkey</button>
    <ul v-for="credential in credentials">
      <li>
        <h3>Nickname: {{credential.Nickname}}</h3>
        <p>Created at: {{credential.CreatedAt}}</p>
        <p>Last used at: {{credential.LastUsedAt}}</p>
        <p>Authenticator: {{getAuthenticatorName(credential.Aaguid)}}</p>
      </li>
    </ul>
  </main>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue";
import {API_BASE, JSON_HEADERS} from "../config.ts";
import {useRoute} from "vue-router";
import {aaguids} from "../aaguids.ts";
import {startRegistration} from "@simplewebauthn/browser";

export interface ApiUser {
  ID: string
  Name: string
  DisplayName?: string
}

export interface Credential {
  ID: string
  UserID: string
  Nickname: string
  PublicKey: string
  AttestationType: string
  Aaguid: string
  SignCount: number
  Transports: string
  UserPresentFlag: boolean
  UserVerifiedFlag: boolean
  BackupEligibleFlag: boolean
  BackupStateFlag: boolean
  CloneWarning: boolean
  CreatedAt: Date
  LastUsedAt: Date
}

const route = useRoute();
const user = ref<ApiUser>();
const credentials = ref<Credential[]>([]);
const busy = ref(false)
const error = ref<string | null>(null)

function getAuthenticatorName(aaguid: string): string {
  return aaguids[aaguid] || "Unknown";
}

async function addPasskey() {
  error.value = null
  busy.value = true
  try {
    const opts = await fetch(`${API_BASE}/passkey/add/begin`, {
      method: 'POST',
      headers: JSON_HEADERS,
      body: JSON.stringify({ user_id: user.value?.ID }),
    }).then((r) => r.json())
    console.log('add begin', opts)

    const cred = await startRegistration({ optionsJSON: opts })
    console.log('add credential', cred)

    const finishUrl = new URL(`${API_BASE}/passkey/add/finish`, window.location.origin)
    finishUrl.searchParams.append('user_id', opts.user.id)

    const result = await fetch(finishUrl, {
      method: 'POST',
      headers: JSON_HEADERS,
      body: JSON.stringify(cred),
    })
    console.log('add finish', result)
    if (!result.ok) throw new Error(`add finish failed: ${result.status}`)
    alert('c bon tu px te connecter avec mtn')
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
    console.error(e)
  } finally {
    busy.value = false
  }
}


onMounted(async () => {
  const userId = route.params.userId;

  user.value = await fetch(`${API_BASE}/users/${userId}`, {
    method: 'GET',
    headers: JSON_HEADERS,
  }).then((r) => r.json()) as ApiUser;

  credentials.value = await fetch(`${API_BASE}/credentials/${userId}/list`, {
    method: 'GET',
    headers: JSON_HEADERS,
  }).then((r) => r.json()) as Credential[];
});
</script>