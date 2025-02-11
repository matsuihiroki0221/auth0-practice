<template>
  <div>
    <button @click="login">Log in</button>
    <pre v-if="isAuthenticated">
          <code>{{ user }}</code>
        </pre>
  </div>
</template>

<script setup lang="ts">
import { useAuth0 } from "@auth0/auth0-vue";

const { loginWithRedirect, user, isAuthenticated } = useAuth0();

const login = async () => {
  await loginWithRedirect();
  console.log("login");
};
</script>
