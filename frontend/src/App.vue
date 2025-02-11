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

// const fetchApiData = async () => {
//   try {
//     const token = await getAccessTokenSilently();
//     const response = await fetch("https://your-api-endpoint.com/data", {
//       headers: {
//         Authorization: `Bearer ${token}`,
//       },
//     });
//     const data = await response.json();
//     console.log(data);
//   } catch (error) {
//     console.error(error);
//   }
// };
</script>
