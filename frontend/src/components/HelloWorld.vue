<script setup lang="ts">
import { ref } from "vue";
import { useAuth0 } from "@auth0/auth0-vue";
import axios from "axios";

defineProps<{ msg: string }>();

const count = ref(0);
const { getAccessTokenSilently } = useAuth0();
const fetchUsers = async () => {
  try {
    const idToken = await getAccessTokenSilently();
    const response = await axios.get("/api/users", {
      headers: {
        Authorization: `Bearer ${idToken}`,
      },
    });
    console.log(response.data);
  } catch (error) {
    console.error("Error fetching users:", error);
  }
};
</script>

<template>
  <h1>{{ msg }}</h1>

  <div class="card">
    <button type="button" @click="count++">count is {{ count }}</button>
    <button type="button" @click="fetchUsers">Fetch Users</button>
    <p>
      Edit
      <code>components/HelloWorld.vue</code> to test HMR
    </p>
  </div>

  <p>
    Check out
    <a href="https://vuejs.org/guide/quick-start.html#local" target="_blank"
      >create-vue</a
    >, the official Vue + Vite starter
  </p>
  <p>
    Learn more about IDE Support for Vue in the
    <a
      href="https://vuejs.org/guide/scaling-up/tooling.html#ide-support"
      target="_blank"
      >Vue Docs Scaling up Guide</a
    >.
  </p>
  <p class="read-the-docs">Click on the Vite and Vue logos to learn more</p>
</template>

<style scoped>
.read-the-docs {
  color: #888;
}
</style>
