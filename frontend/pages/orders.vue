<script setup lang="ts">
import AwModal from "~/components/ui-kit/AwModal.vue";

const router = useRouter();
definePageMeta({layout: 'default'})
useHead({
  title: 'Заказы'
})
const role = useCookie('role');
const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');

role.value === 'contractor' ? router.push('/offers') : '';
const orders = ref()
// const offers = [
//   {
//     id: 1,
//     title: 'Установка и настройка 1С:Предприятие 8.3',
//     price: 5000,
//     description: 'Необходимо установить и настроить программу 1С:Предприятие 8.3 на компьютер с операционной системой Windows. Также требуется настройка пользователя и баз данных в соответствии с нашими требованиями. Работа должна быть выполнена качественно и в срок.',
//     rating: 4.99,
//     viewers: 53,
//     author: 'ООО АРУ',
//     status: ''
//   }
// ]

fetch(`${runtimeConfig.public.apiBase}/me/orders`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token.value}`
  }
}).then((response) => {
  return response.json();
}).then((data) => {
  // console.log(data)
  orders.value = data.orders
}).catch((err) => {
  console.error("Невозможно отправить запрос", err);
});

const messageShow = ref(false)

function doneOrder(id) {
  fetch(`${runtimeConfig.public.apiBase}/me/orders/${id}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      status: 'done',
    })
  }).then((response) => {
    switch (response.status) {
      case 400:
        alert('Ошибка')
        return
      case 200:
        window.location.reload()
        return
    }
    return response.json();
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

const showDataModal = ref(false)
const finalComment = ref('')

function showData(final_comment) {
  showDataModal.value = true
  finalComment.value = final_comment.value
  if (final_comment == '') {
    finalComment.value = 'Нет данных'
  }
}

const reviewModal = ref(false)
const reviewText = ref('')
const reviewRating = ref(0)
const reviewOrderId = ref(0)

function sendReview() {
  fetch(`${runtimeConfig.public.apiBase}/reviews/`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      order_id: reviewOrderId.value,
      text: reviewText.value,
      rating: parseInt(reviewRating.value)
    })
  }).then((response) => {
    switch (response.status) {
      case 201:
        alert('Ok')
        return
        break;
      default:
        alert('Ошибка, проверьте отзыв');
        break
    }
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function openReviewModal(id) {
  reviewOrderId.value = id;
  reviewModal.value = true;
}

</script>


<template>
  <div class="main-content">
    <div class="row">
      <h1 class="page-title">Мои заказы (заказчик)</h1>
      <button class="btn-primary new-order">
        <NuxtLink to="/order/new">Создать</NuxtLink>
      </button>
    </div>
    <div class="container">
      <div class="column">
        <div class="offers">
          <div class="offer block" v-for="order in orders">
            <AwModal @close="showDataModal = !showDataModal" :modalActive="showDataModal"
                     :title="'Данные о заказе'">
              <p>Итоговый комментарий:</p>
              <textarea disabled v-model="finalComment"></textarea>
            </AwModal>
            <AwModal @close="reviewModal = !reviewModal" :modalActive="reviewModal"
                     :title="'Отзыв об исполнителе'">
              <p>Оставить отзыв об исполнителе:</p>
              <textarea v-model="reviewText"></textarea>
              <p>Рейтинг (от 1 до 5):</p>
              <input v-model="reviewRating" placeholder="От 1 до 5" min="1" max="5">
              <button @click="sendReview" class="btn-primary" style="margin-top: 10px;">Отправить
              </button>
            </AwModal>

            <AwModal @close="messageShow = !messageShow" :modalActive="messageShow"
                     :title="'Связь с исполнителем'">
              <div class="modal-content">
                <!--                TODO: Тут надо дописать, чтобы почта бралась из order-->
                <a :href="'mailto:'">
                  <button class="btn-primary" style="margin-right: 5px;">Написать на email</button>
                </a>
                <button class="btn-primary">Перейти в чат</button>
              </div>
            </AwModal>
            <div class="offer-info">
              <NuxtLink class="order-name" :to="'/order/' + order.order_id">{{ order.name }}</NuxtLink>
              <div class="row">
                <div class="offers">
                  <div class="icon-offer"></div>
                  <span>{{ order.offer_count }}</span>
                </div>
                <div class="viewers">
                  <div class="icon-eye"></div>
                  <span>{{ order.views }}</span>
                </div>
              </div>
              <button class="btn-primary" v-if="order.status!=='done'" @click="doneOrder(order.order_id)">Завершить
                заказ
              </button>
              <button class="btn-primary" v-if="order.status==='done'" @click="showData(order.final_comment)"
                      style="margin-right: 15px;">
                Посмотреть данные о заказе
              </button>
              <button class="btn-primary" v-if="order.status==='done'" @click="openReviewModal(order.order_id)"
                      style="margin-top: 15px">
                Оставить отзыв об исполнителе
              </button>
            </div>
            <div class="offer-buttons">
              <p class="status"
                 v-bind:class="order.status === 'available'?'available':order.status==='cancelled'?'cancel':order.status==='in_progress'?'in_progress':''">
                {{
                  order.status === 'available' ? 'Поиск исполнителя' : order.status === 'cancelled' ? 'Удалено' : order.status === 'in_progress' ? 'Выполняется' : order.status === 'done' ? 'Завершено' : order.status
                }}
              </p>
              <NuxtLink :to="'/order/' + order.order_id">
                <button class="btn-primary">Перейти к заказу</button>
              </NuxtLink>
              <div class="offer-buttons-group">
                <button class="btn-primary" v-if="order.status !== 'available'" @click="messageShow = !messageShow">
                  Связь с исполнителем
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!--      <div class="column">-->
      <!--        <div class="block">-->
      <!--          asdfads-->
      <!--        </div>-->
      <!--      </div>-->
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.offer-info {
  width: inherit;

  .row {
    justify-content: flex-start;
  }

  .btn-primary {
    width: inherit;
  }
}

.row {
  .rating, .viewers, .offers {
    display: flex;
    align-items: center;
  }
}

.icon-offer {
  width: 25px !important;
  height: 25px !important;
}

.offer .row:last-child {
  justify-content: flex-start
}

.offer {
  margin-bottom: 15px;
  flex-direction: row !important;
  justify-content: space-between;
}

.offer:last-child {
  margin-bottom: 0;
}

.new-order a {
  color: white;
}

a {
  color: black;
}

.offer-buttons *, .offer .btn-primary {
  font-size: medium;
  width: 210px;
  text-align: center;
}

.offer.block {
  height: 230px;
}

.offer-buttons {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

p.status {
  padding: 10px 20px;
  border-radius: 10px;
}

p.status.available {
  background: #ef7b42;
  color: white;
}

p.status.in_progress {
  background-color: #ffc50c;
  color: white;
  //color: white;
}

p.status.done {
  color: $primary-color;
}

p.status.cancel {
  color: brown;
  background: none;
}

</style>