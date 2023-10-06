<script setup lang="ts">
import AwModal from "~/components/ui-kit/AwModal.vue";

definePageMeta({layout: 'default'})
useHead({
  title: 'Главная'
})
const role = useCookie('role');
const router = useRouter();
role.value === 'customer' ? router.push('/orders') : '';
const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
const offers = ref()

fetch(`${runtimeConfig.public.apiBase}/me/offers`, {
  method: 'GET',
  cache: 'no-cache',
  headers: {
    'Authorization': `Bearer ${token.value}`,
    'cache': "no-store"
  }
}).then((response) => {
  return response.json();
}).then((data) => {
  offers.value = data.offers.filter(item => item.id !== 0)
}).catch((err) => {
  console.error("Невозможно отправить запрос", err);
});

const finalComment = ref('')

function doneOffer(id) {
  fetch(`${runtimeConfig.public.apiBase}/offers/${id}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      final_comment: {
        text: finalComment.value,
        attachments: attachments.value
      }
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        router.push('/offers')
        finalCommentShow.value = false
        return
        break;
      default:
        alert('Ошибка, попробуйте ещё раз')
        break;
        return;
    }
  }).then((data) => {
    console.log(data)
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

const finalCommentShow = ref(false)
const attachments = ref([])

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
        alert('Файл не прикрепился, попробуйте ещё раз')
        break;
    }
    return response.json();
  }).then((data) => {
    attachments.value.push(data.uri)
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

const messageShow = ref(false)

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
    <h1 class="page-title">Мои заказы (исполнитель)</h1>
    <div class="container">
      <div class="column">
        <div class="offers">
          <div class="offer block" v-for="offer in offers">
            <div class="offer-info">
              <div class="order-name">
                <NuxtLink :to="'/order/'+offer.order.id">{{ offer.order.name }}</NuxtLink>
              </div>
              <div class="order-comment" v-if="offer.status==='submitted'" style="margin-top: 15px">
                Ваш комментарий к заказу: {{ offer.comment }}
              </div>
              <button class="btn-primary" @click="finalCommentShow = !finalCommentShow"
                      v-if="offer.status==='confirmed'" style="margin-top: 15px">Создать/изменить ответ
              </button>

              <button class="btn-primary" v-if="offer.status==='done'" @click="openReviewModal(offer.order.id)"
                      style="margin-top: 15px">
                Оставить отзыв об исполнителе
              </button>

              <AwModal @close="finalCommentShow = !finalCommentShow" :modalActive="finalCommentShow"
                       :title="'Итоговый комментарий'">
                <div class="modal-content">
                  <textarea class="order-final-comment" v-if="offer.status==='confirmed'"
                            v-model="finalComment"></textarea>
                  <b>Прикрепить материалы:</b>
                  <div class="applications">
                    <label class="input-file">
                      <input type="file" name="file" @change="load">
                    </label>
                  </div>
                  <button class="btn-primary" style="margin-top: 10px" @click="doneOffer(offer.id)">Сдать заказ</button>
                </div>
              </AwModal>
              <AwModal @close="reviewModal = !reviewModal" :modalActive="reviewModal"
                       :title="'Отзыв об исполнителе'">
                <p>Оставить отзыв об исполнителе:</p>
                <textarea v-model="reviewText"></textarea>
                <p>Рейтинг (от 0 до 5):</p>
                <input v-model="reviewRating" placeholder="От 0 до 5" min="0" max="5">
                <button @click="sendReview()" class="btn-primary" style="margin-top: 10px">Отправить
                </button>
              </AwModal>
              <AwModal @close="messageShow = !messageShow" :modalActive="messageShow"
                       :title="'Связь с заказчиком'">
                <div class="modal-content">
                  <a :href="'mailto:'+offer.order.email">
                    <button class="btn-primary" style="margin-right: 5px;">Написать на email</button>
                  </a>
                  <button class="btn-primary">Перейти в чат</button>
                </div>
              </AwModal>
            </div>
            <div class="offer-buttons">
              <p class="status processing"
                 v-bind:class="offer.status === 'submitted'?'submitted':offer.status==='confirmed'?'confirmed':offer.status==='in_progress'?'in_progress':''">
                {{
                  offer.status === 'submitted' ? 'Отправлено' : offer.status === 'confirmed' ? 'Утверждено' : offer.status === 'done' ? 'Выполнено' : 'Ошибка'
                }}
              </p>
              <div class="offer-buttons-group">
                <NuxtLink :to="'/order/'+offer.order.id" class="btn-primary">
                  Перейти к заказу
                </NuxtLink>
                <button class="btn-primary" @click="messageShow = !messageShow">
                  Связь с заказчиком
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
}

.offer-info .btn-primary {
  width: inherit !important;
}

.offer-buttons-group > * {
  margin-bottom: 5px;
}

.offer-buttons-group > *:last-child {
  margin: 0;
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

a {
  color: black;
}

.offer-buttons *, .offer .btn-primary {
  font-size: medium;
  width: 210px;
  text-align: center;
}

.offer.block {
  min-width: 230px;
}

@media (max-width: 470px) {
  .offer.block {
    flex-wrap: wrap;
  }
}

.offer-buttons {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 5px;
}

p.status {
  padding: 10px 20px;
  border-radius: 10px;
}

p.status.processing {
  background: #ef7b42;
  color: white;
}

p.status.done {
  color: $primary-color;
}

p.status.cancel {
  color: white;
  background: orange;
}

</style>