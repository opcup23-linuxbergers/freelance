<script setup lang="ts">

const value = ref("");
const formattedValue = ref("");

watch(formattedValue, () => {
  // Remove non-numeric characters and spaces from input value
  const numericValue = formattedValue.value.replace(/[^0-9]/g, "");

  // Add spaces as separators every 3 characters
  const separatedValue = numericValue.replace(/\B(?=(\d{3})+(?!\d))/g, " ");

  // Update the original value with the separated value
  value.value = separatedValue;
});

const formatInput = (event: { target: { value: string; }; }) => {
  // Remove non-numeric characters and spaces from input value
  const numericValue = event.target.value.replace(/[^0-9]/g, "");

  // Add spaces as separators every 3 characters
  const separatedValue = numericValue.replace(/\B(?=(\d{3})+(?!\d))/g, " ");

  // Update the formatted value with the separated value
  formattedValue.value = separatedValue;
};

</script>

<template>
  <input type="text" v-model="formattedValue" @input="formatInput"/>
</template>
