<script setup lang="ts">

import AwСheckmark from "~/components/ui-kit/AwСheckmark.vue";

definePageMeta({layout: 'default'})
useHead({
  title: 'Заказ'
})

const role = useCookie('role');
const owner = ref(false);
const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
const route = useRoute()
const router = useRouter()
const orders = ref()
const order = ref({
  "order_id": 0,
  "user": "",
  "rating": 0,
  "name": "",
  "description": "",
  "price": 0,
  "status": "",
  "due": 0,
  "published": 0,
  "views": 0,
  "offer_count": 0,
  "final_comment": {
    "text": "",
    "attachments": []
  },
  "attachments": []
})
const contractors = ref()
const offerData = ref({
  due: '',
  price: 0,
  comment: ''
})
const error = ref()
const complete = ref(false)

// function unixToDate(unixTime) {
//   if (unixTime === 0) return
//   const date = new Date(unixTime);
//   const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
//   const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1;
//   const year = date.getFullYear().toString();
//   return `${day}.${month}.${year}`;
// }


function unixToDate(unixTime) {
  const date = new Date(unixTime * 1000);
  const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
  const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1;
  const year = date.getFullYear()
  return `${day}.${month}.${year}`;
}

// function dateToUnix(dateString) {
//   const dateParts = dateString.split('-');
//   const year = dateParts[0];
//   const month = dateParts[1] - 1;
//   const day = dateParts[2];
//   return Date.parse(new Date(year, month, day))
// }

function getOrder() {
  fetch(`${runtimeConfig.public.apiBase}/orders/${route.params.id}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.value}`,
      'cache': "no-cache"
    }
  }).then((response) => {
    return response.json();
  }).then((data) => {
    if (data && !data.uri) {
      order.value = data
    }
    if (data.uri) {
      getMyOrder(data.uri)
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function getMyOrder(uri) {
  fetch(`${runtimeConfig.public.apiBase}${uri}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.value}`
    }
  }).then((response) => {
    return response.json();
  }).then((data) => {
    if (data) {
      order.value = data
      owner.value = true
      getOffers()
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function getOffers() {
  fetch(`${runtimeConfig.public.apiBase}/me/orders/${route.params.id}/offers`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.value}`
    }
  }).then((response) => {
    switch (response.status) {
      case 200:
        return response.json();
        break;
      case 403:
        console.log('Нет прав')
        owner.value = false
        // router.push('/')
        break;
    }
  }).then((data) => {
    if (data) {
      // console.log(data.filter(item => item.status === 'confirmed'))
      if (data.offers.filter(item => item.status === 'confirmed').length !== 1) {
        contractors.value = data.offers.filter(item => item.id !== 0)
      }
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function sendOffer() {
  fetch(`${runtimeConfig.public.apiBase}/orders/${route.params.id}/offers`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      due: Date.parse(offerData.value.due) / 1000,
      price: offerData.value.price,
      comment: offerData.value.comment
    })
  }).then((response) => {
    switch (response.status) {
      case 409:
        error.value = 'Вы уже откликались на этот заказ'
        break;
      case 400:
        error.value = 'Проверьте корректность данных'
        break;
      case 201:
        complete.value = true
        router.push('/offers')
        return
        break;
      default:
        error.value = 'Непредвиденные обстоятельства'
        break;
    }
    return response.json();
  }).then((data) => {
    // console.log(data)
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
    error.value = 'Сервер не отвечает'
  });
}

function deleteOrder() {
  fetch(`${runtimeConfig.public.apiBase}/me/orders/${route.params.id}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token.value}`
    }
  }).then((response) => {
    switch (response.status) {
      case 200:
        router.push('/')
        break;
      case 403:
        console.log('Статьи не существует')
        break;
    }
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

getOrder()

function acceptOffer(order_id, offer_id) {
  fetch(`${runtimeConfig.public.apiBase}/me/orders/${order_id}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      offer_id: offer_id
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        alert('Ok')
        return
        break;
      case 404:
        alert('Данного заказа не существует')
        break;
      case 409:
        alert('В данном заказе уже присутсвует исполнитель')
        break;
      default:
        break;
    }
    return response.json();
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
    error.value = 'Сервер не отвечает'
  });
}

function hideOffer(id) {
  fetch(`${runtimeConfig.public.apiBase}/offers/${id}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      hidden: true
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        getOffers()
        return;
      case 404:
        alert('Данного оффера не существует')
        return;
      default:
        alert('Ошибка')
        break;
    }
    return response.json();
  }).then((data) => {
    console.log(data)
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

</script>

<template>
  <div class="main-content">
    <AwСheckmark :complete="complete" @toggleComplete="complete = !complete"></AwСheckmark>
    <section class="row">
      <h1 class="page-title">{{ order.name }}</h1>
      <section>
        <button class="btn-primary delete" v-if="owner" v-on:click="deleteOrder">Удалить</button>
        <NuxtLink :to="'/order/edit/'+order.order_id">
          <button class="btn-primary" v-if="owner">Редактировать</button>
        </NuxtLink>
      </section>
    </section>
    <div class="container" v-if="order.status !== 'confirmed'">
      <div class="column" style="flex-basis: 30%;">
        <div class="description block">
          <h3>Описание</h3>
          <p>{{ order.description }}</p>
        </div>
        <!--        <div class="files block">-->
        <!--          <h3>Приложения</h3>-->
        <!--          <ul>-->
        <!--            <li>ТЗ.docx</li>-->
        <!--          </ul>-->
        <!--        </div>-->
        <div class="block" v-if="order.final_comment.text">
          <label style="font-weight: bold">Итоговый комментарий</label>
          <p> {{ order.final_comment.text }}</p>
        </div>
        <div class="respond block" v-if="role==='contractor' && !owner">
          <h3>Ответить на заказ</h3>
          <form v-on:submit.prevent="sendOffer">
            <div class="form-error" v-if="error">
              {{ error }}
            </div>
            <div class="input-group">
              <label for="days">Срок исполнения: </label>
              <input id="days" name="days" type="date" minlength="1" required v-model="offerData.due">
            </div>
            <div class="input-group">
              <label for="price">Ваша стоимость:</label>
              <input id="price" name="price" type="number" minlength="1" required v-model="offerData.price">
            </div>
            <div class="input-group">
              <label for="comment">Комментарий:</label>
              <textarea id="comment" name="comment" placeholder="Комментарий" rows="6"
                        v-model="offerData.comment"></textarea>
            </div>
            <button class="btn-primary">
              Откликнуться
            </button>
          </form>
        </div>
      </div>
      <div class="column">
        <div class="customer-info block">
          <h4>О заказчике</h4>
          <p>Заказчик: {{ order.user }}</p>
          <p>Рейтинг: {{ order.rating }}</p>
        </div>
        <div class="order-info block">
          <h4>О заказе</h4>
          <p>Крайний срок: {{ unixToDate(order.due) }}</p>
          <p>Опубликован: {{ unixToDate(order.published) }}</p>
          <p>Стоимость:
            <span v-if="order.price>0">{{ order.price }} ₽/заказ</span>
            <span v-else>договорная</span>
          </p>
        </div>
      </div>
    </div>
    <div class="container" v-if="contractors">
      <div class="column">
        <div class="requests">
          <div class="contractors block" v-for="contractor in contractors">
            <div class="row">
              <h3 class="name">{{ contractor.user }}</h3>
              <div class="rating">
                <div class="icon-rating"></div>
                <span>{{ contractor.rating }}</span>
              </div>
            </div>
            <div class="request-info">
              <p>Выполню ваш заказ до {{ unixToDate(contractor.due) }}</p>
              <p>Стоимость: {{ contractor.price }} ₽</p>
            </div>
            <p>{{ contractor.comment }}</p>
            <div class="control-panel">
              <button class="btn-primary accept" v-on:click="acceptOffer(contractor.order.id,contractor.id)">Принять
              </button>
              <NuxtLink to="/chats">
                <button class="btn-primary message">Связаться</button>
              </NuxtLink>
              <button class="btn-primary message cancel" v-on:click="hideOffer(contractor.id)">Скрыть</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="container" v-if="!role">
      <div class="column">
        <div class="block get-login">
          <p>Хотите откликнуться?</p>
          <NuxtLink to="/authorization">Войдите</NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.form-error {
  margin: 0;
}

.btn-primary.delete {
  background: none;
  color: red;
}

.rating {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 28px;
}

ul {
  list-style: none;
}

form > * {
  margin-bottom: 15px;
}

form > *:last-child {
  margin-bottom: 0;
}

.input-group {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
}

.contractors button {
  font-size: medium;
  width: 210px;
  text-align: center;
}

.message {
  background: none;
  color: $primary-color;
  width: 90px !important;
}

.message:hover {
  background: none;
}

.contractors {
  border: 1px solid $primary-color;
  padding: 15px;
  border-radius: 10px;
  margin-bottom: 15px;
}

.contractors:last-child {
  margin-bottom: 0;
}

.contractors > * {
  margin-bottom: 15px;
}

.contractors > *:last-child {
  margin-bottom: 0;
}

.cancel {
  color: red;
}

.get-login a {
  color: $primary-color;
}
</style>