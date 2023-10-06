<script setup lang="ts">
const inputValue = ref('')

const validateInput = () => {
  // Remove non-numeric characters from input value
  const numericValue = inputValue.value.replace(/[^0-5]/g, '')

  // Keep only the first digit if more than one is entered
  inputValue.value = numericValue.slice(0, 1)
}

</script>

<template>
  <input type="text" v-model="inputValue" @input="validateInput" />
</template>