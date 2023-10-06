<script setup lang="ts">
import AwModal from "~/components/ui-kit/AwModal.vue";

const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
const role = useCookie('role');
const router = useRouter();

const error = ref('')
const loading = ref(false)
const modalActive = ref(false)
const data = ref({
  email: '',
  password: ''
})

function authorization() {
  loading.value = true
  fetch(`${runtimeConfig.public.apiBase}/login`, {
    method: 'POST',
    body: JSON.stringify({
      email: data.value.email,
      password: data.value.password
    })
  }).then((response) => {
    switch (response.status) {
      case 400:
        error.value = 'Неверный логин или пароль'
        loading.value = false
        break;
      case 200:
        return response.json();
      default:
        error.value = 'Проверьте правильность данных'
        loading.value = false
        return
        break;
    }
  }).then((data) => {
    if (data) {
      try {
        token.value = data.token;
      } catch (e) {
        console.log(e, 'Не получается распарсить данные')
      }
      role.value = 'contractor'
      modalActive.value = true;
      loading.value = false
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
    loading.value = false
  });
}

function changeRole(newRole: string) {
  role.value = newRole
  modalActive.value = false
  router.push('/')
}

</script>

<template>
  <AwModal @close="changeRole('contractor')" :modalActive="modalActive" :title="'Выберите роль'">
    <div class="modal-content">
      <button class="btn-primary" v-on:click="changeRole('contractor')">Исполнитель</button>
      <button class="btn-primary" v-on:click="changeRole('customer')">Заказчик</button>
    </div>
  </AwModal>
  <div class="auth">
    <form class="form" v-on:submit.prevent="authorization">
      <div class="header">Войти</div>
      <div class="form-error" v-if="error">
        {{ error }}
      </div>
      <div class="inputs">
        <input placeholder="Email" class="input" type="email" v-model="data.email" minlength="3" required>
        <input placeholder="Пароль" class="input" type="password" v-model="data.password" minlength="8" autocomplete
               required>
        <button class="sigin-btn btn-primary" v-bind:class="loading?'skeleton':''" :disabled="loading">Войти</button>
        <p class="signup-link">Ещё не зарегистрированы?
          <NuxtLink to="/registration">Регистрация</NuxtLink>
        </p>
      </div>
    </form>
  </div>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.auth {
  margin: 0 auto;
}

.modal-content {
  width: 400px;
  display: flex;
  justify-content: space-around;
  gap: 5px;
}

.form {
  position: relative;
  display: flex;
  flex-direction: column;
  border-radius: 0.75rem;
  background-color: #fff;
  color: rgb(97 97 97);
  box-shadow: 20px 20px 30px rgba(0, 0, 0, .05);
  width: 22rem;
  background-clip: border-box;
}

.header {
  position: relative;
  background-clip: border-box;
  background-color: $primary-color;
  //background-image: linear-gradient(to top right, #1e88e5, #42a5f5);
  margin: 10px;
  border-radius: 0.75rem;
  overflow: hidden;
  color: #fff;
  box-shadow: 0 0 #0000, 0 0 #0000, 0 0 #0000, 0 0 #0000, rgba(33, 150, 243, .4);
  height: 7rem;
  letter-spacing: 0;
  line-height: 1.375;
  font-weight: 600;
  font-size: 1.9rem;
  font-family: Roboto, sans-serif;
  display: flex;
  align-items: center;
  justify-content: center;
}

.inputs {
  padding: 1.5rem;
  gap: 1rem;
  display: flex;
  flex-direction: column;
}

.input-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-width: 200px;
  width: 100%;
  height: 2.75rem;
  position: relative;
}

.input {
  border: 1px solid rgba(128, 128, 128, 0.61);
  outline: 0;
  color: rgb(69 90 100);
  font-weight: 400;
  font-size: .9rem;
  line-height: 1.25rem;
  padding: 0.75rem;
  background-color: transparent;
  border-radius: .375rem;
  width: 100%;
  height: 100%;
}

.input:focus {
  border: 1px solid #1e88e5;
}

.sigin-btn {
  text-transform: uppercase;
  font-weight: 700;
  font-size: .75rem;
  line-height: 1rem;
  text-align: center;
  padding: .75rem 1.5rem;
  background-color: $primary-color;
  //background-image: linear-gradient(to top right, #1e88e5, #42a5f5);
  border-radius: .5rem;
  width: 100%;
  outline: 0;
  border: 0;
  color: #fff;
  cursor: pointer;
}

.signup-link {
  line-height: 1.5;
  font-weight: 300;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.signup-link a, .forget {
  line-height: 1.5;
  font-weight: 700;
  font-size: .875rem;
  margin-left: .25rem;
  color: #1e88e5;
}

</style>