<script setup lang="ts">

import AwButton from "~/components/ui-kit/AwButton.vue"

const props = defineProps({
    title: String,
    modalActive: {
        type: Boolean,
    }
})

const emit = defineEmits(['close'])

const close = () => {
    emit("close")
}

</script>

<template>
    <transition name="modal-animation">
        <div v-show="modalActive" class="modal" @click="close">
            <transition name="modal-animation-inner">
                <div class="modal-inner" @click.stop>
                    <div @click="close" class="icon-close close-button" style="cursor: pointer"></div>
                    <h2 style="margin-bottom: 16px; font-size: 1.5rem">{{ props.title }}</h2>
                    <slot />
                </div>
            </transition>
        </div>
    </transition>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.modal-animation-enter-active,
.modal-animation-leave-active {
    transition: opacity 0.3s cubic-bezier(0.52, 0.02, 0.19, 1.02);
}

.modal-animation-enter-from,
.modal-animation-leave-to {
    opacity: 0;
}

.modal-animation-inner-enter-active {
    transition: all 0.1s cubic-bezier(0.52, 0.02, 0.19, 1.02) 0.15s;
}

.modal-animation-inner-leave-active {
    transition: all 0.1s cubic-bezier(0.52, 0.02, 0.19, 1.02);
}

.modal-animation-inner-enter-from {
    opacity: 0;
    transform: scale(0.8);
}

.modal-animation-inner-leave-to {
    transform: scale(0.8);
}

.modal {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    width: 100vw;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    background-color: $opacity-color;

    .modal-inner {
        position: relative;
        max-width: 650px;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        background: $general-background-light;
        border-radius: 20px;
        padding: 30px;
        max-height: 95vh;
        overflow: auto;

        & .close-button {
            position: absolute;
            top: 15px;
            right: 15px;
            font-size: 20px;
            cursor: pointer;
            transition: 0.4s;

            &:hover {
                transform: scale(1.5);
            }
        }
    }
}
</style>