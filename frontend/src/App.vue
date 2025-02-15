<template>
  <div v-if="isAuthenticated">
    <HelloWorld msg="Hello Vue 3 + Vite + TypeScript + Auth0" />
    <code>{{ user }}</code>
    <LogOut />
  </div>
  <div v-else>
    <LogIn />
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from "@auth0/auth0-vue";
import HelloWorld from "./components/HelloWorld.vue";
import LogIn from "./components/LogIn.vue";
import LogOut from "./components/LogOut.vue";

const { isAuthenticated, user } = useAuth0();
</script>
