<script setup lang="ts">

import TheHeader from "~/components/TheHeader.vue";

</script>

<template>
  <div class="wrapper">
    <TheHeader></TheHeader>
    <div class="page">
      <slot/>
    </div>
  </div>
</template>

<style lang="scss">
@use '@/assets/scss/main.scss' as *;

.wrapper {
  display: flex !important;
  flex-direction: column;
  flex-wrap: wrap;
  justify-content: flex-start;
  min-width: 100%;
  width: 100%;
  min-height: 100vh;
  max-width: 1920px;
  background-color: $page-background;
}

.page {
  position: relative;
  display: flex !important;
  flex: 1 1 auto;
  margin: 0 auto;
  padding: 25px;
  width: 100%;
  min-height: 100%;
  max-width: 1920px;
}
</style>