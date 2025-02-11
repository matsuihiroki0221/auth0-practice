<template>
  <div>
    <button @click="logoutApp">Log out</button>
  </div>
</template>
<script setup lang="ts">
import { useAuth0 } from "@auth0/auth0-vue";
const { logout } = useAuth0();

const logoutApp = () => {
  logout({ logoutParams: { returnTo: window.location.origin } });
};
</script>
