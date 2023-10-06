<script setup lang="ts">
const props = defineProps({
  complete: {
    type: Boolean,
    required: true,
  },
  timeout: {
    type: Number,
    default: 3000
  },
})

const emit = defineEmits(['toggleComplete'])

const toggleComplete = () => {
  emit("toggleComplete")
}

function hideCheckmark() {
  if (props.complete == true) {
    setTimeout(() => {
      toggleComplete()
    }, props.timeout)
  }
}

onMounted(() => {
  hideCheckmark()
})

watch(() => props.complete?.valueOf(), () => {
  hideCheckmark()
});
</script>

<template>
  <transition name="fade">
    <div class="checkmark-wrapper" v-if="props.complete">
      <svg class="checkmark" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 52">
        <circle class="checkmark__circle" cx="26" cy="26" r="24" fill="none"/>
        <path class="checkmark__check" fill="none" d="M14.1 27.2l7.1 7.2 16.7-16.8"/>
      </svg>
    </div>
    </transition>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.checkmark-wrapper {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 10;
}

.checkmark__circle {
  stroke-dasharray: 166;
  stroke-dashoffset: 166;
  stroke-width: 2;
  stroke-miterlimit: 10;
  stroke: $primary-color;
  fill: none;
  animation: stroke 0.6s cubic-bezier(0.65, 0, 0.45, 1) forwards;
}

.checkmark {
  width: 156px;
  height: 156px;
  border-radius: 50%;
  display: block;
  stroke-width: 2;
  stroke: #fff;
  stroke-miterlimit: 10;
  margin: 10% auto;
  box-shadow: $primary-color inset 0 0 0;
  animation: fill .4s ease-in-out .4s forwards, scale .3s ease-in-out .9s both;
}

.checkmark__check {
  transform-origin: 50% 50%;
  stroke-dasharray: 48;
  stroke-dashoffset: 48;
  animation: stroke 0.3s cubic-bezier(0.65, 0, 0.45, 1) 0.8s forwards;
}

@keyframes stroke {
  100% {
    stroke-dashoffset: 0;
  }
}

@keyframes scale {
  0%, 100% {
    transform: none;
  }
  50% {
    transform: scale3d(1.1, 1.1, 1);
  }
}

@keyframes fill {
  100% {
    box-shadow: $primary-color inset 0 0 0 80px;
  }
}
</style>