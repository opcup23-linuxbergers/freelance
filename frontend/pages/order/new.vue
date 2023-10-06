<script setup lang="ts">
import AwСheckmark from "~/components/ui-kit/AwСheckmark.vue";

const runtimeConfig = useRuntimeConfig();
const due = ref()
const token = useCookie('token');
const router = useRouter();
const orderData = ref({
  name: "",
  description: "",
  price: 0
})
const error = ref('')
let attachments = ref([])

function newOrder() {
  fetch(`${runtimeConfig.public.apiBase}/orders`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      name: orderData.value.name,
      description: orderData.value.description,
      price: orderData.value.price,
      due: Date.parse(due.value) / 1000,
      attachments: attachments.value,
    })
  }).then((response) => {
    switch (response.status) {
      case 400:
        error.value = 'Неверные данные, проверьте свой балланс'
        return
      case 201:
        router.push({path: "/orders"});
        return
    }
    return response.json();
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function load() {
  const input = document.querySelector('input[type="file"]')
  let data = new FormData()
  data.append('uploadedfile', input.files[0])
  fetch(`${runtimeConfig.public.apiBase}/uploads`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`,
    },
    body: data
  }).then((response) => {
    switch (response.status) {
      case 200:
        alert('Файл загружен успешно')
        break;
      case 412:
        alert('Файл слишком большой')
        break;
      case 500:
        alert('Файл не выбран')
        break;
    }
    return response.json();
  }).then((data) => {
    attachments.value.push(data.uri)
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

</script>

<template>
  <div class="main-content">
    <h1 class="page-title">Создание заказа</h1>
    <div class="container">
      <div class="column">
        <form class="border-info block" v-on:submit.prevent="newOrder">
          <div class="form-error" v-if="error">
            {{ error }}
          </div>
          <label for="">Название заказа</label>
          <input type="text" v-model="orderData.name">
          <label for="">Описание задачи</label>
          <textarea name="" id="" cols="30" rows="10" v-model="orderData.description"></textarea>
          <!--          <label for="">Приложения</label>-->
          <!--          <input type="file">-->
          <label for="">Крайний срок</label>
          <input type="date" v-model="due">
          <label for="">Стоимость</label>
          <input type="number" required min="0" minlength="1" v-model="orderData.price">

          <b>Прикрепить материалы:</b>
          <div class="applications">
            <label class="input-file">
              <input type="file" name="file" @change="load">
              <!--            <input type="file" name="file" @change="changeMaterial(id)">-->
              <!--            <span class="input-file-btn">Выберите файл</span>-->
              <!--            <span class="input-file-text">{{  }}</span>-->
            </label>
          </div>

          <h4>Вложения:</h4>
          <ul>
            <li v-for="item in attachments">{{ item}}</li>
          </ul>

          <button class="btn-primary">Создать заказ</button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.form-error {
  margin: 0;
}

/* Focus */
.input-file input[type=file]:focus + .input-file-btn {
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, .25);
}

.input-file:active .input-file-btn {
  background-color: #2E703A;
}

/* Disabled */
.input-file input[type=file]:disabled + .input-file-btn {
  background-color: #eee;
}

</style>