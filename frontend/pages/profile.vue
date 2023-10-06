<script setup lang="ts">
import AwModal from "~/components/ui-kit/AwModal.vue";
import AwСheckmark from "~/components/ui-kit/AwСheckmark.vue";

definePageMeta({layout: 'default'})
useHead({
  title: 'Личный кабинет'
})

const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
const role = useCookie('role');
const modalActive = ref(false)
const modalActiveSkills = ref(false)
const complete = ref(false)
let stories = ref()
// let stories = [
//   {type: "Операция пополнения", name: "ООО \"ТехноПрогресс\"", price: 15000, date: "2012-04-23T18:25:43.511Z"},
//   {type: "Операция пополнения", name: "ООО \"ТехноПрогресс\"", price: 15000, date: "2012-04-23T18:25:43.511Z"}
// ];

const skills = ref([
  {id: 1, name: "Docker"},
  {id: 2, name: "Python"},
  {id: 3, name: "Rust"}
]);
// let reviews = [
//   {
//     name: "ООО \"ТехноПрогресс\"",
//     desc: "Мы были очень довольны работой разработчика. С его помощью нам удалось создать макет для нашего нового сайта фриланса на высоком уровне."
//   },
//   {
//     name: "ИП \"ВебСтарт\"",
//     desc: "Разработчик продемонстрировал высокие навыки веб-разработки. Специалист хорошо ориентируется в React.JS и Next.JS"
//   },
// ]

function skillsDelete(index: number) {
  skills.value.splice(index, 1);
}

const userData = ref({
  "alias": "",
  "first_name": "",
  "last_name": "",
  "about": "",
  "seller_rating": 0,
  "buyer_rating": 0
})
//
// function unixToDate(unixTime) {
//   const date = new Date(unixTime);
//   const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
//   const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1;
//   const year = date.getFullYear().toString().substr(-2);
//   return `${day}.${month}.${year}`;
// }

function unixToDate(unixTime) {
  const date = new Date(unixTime * 1000);
  const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate();
  const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1;
  const year = date.getFullYear()
  return `${day}.${month}.${year}`;
}

fetch(`${runtimeConfig.public.apiBase}/me`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token.value}`
  }
}).then((response) => {
  return response.json();
}).then((data) => {
  // console.log(data)
  if (data) {
    userData.value = data
  }
}).catch((err) => {
  console.error("Невозможно отправить запрос", err);
});

const balance = ref()

fetch(`${runtimeConfig.public.apiBase}/balance`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token.value}`
  }
}).then((response) => {
  return response.json();
}).then((data) => {
  // console.log(data)
  if (data) {
    balance.value = data.balance
  }
}).catch((err) => {
  console.error("Невозможно отправить запрос", err);
});

function getHistory() {
  fetch(`${runtimeConfig.public.apiBase}/balance/history`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.value}`
    }
  }).then((response) => {
    return response.json();
  }).then((data) => {
    // console.log(data)
    if (data) {
      stories.value = data.history
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

getHistory()

const delta = ref()

function deposit() {
  fetch(`${runtimeConfig.public.apiBase}/balance`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      "operation_type": "deposit",
      "delta": delta.value
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        complete.value = true
        balance.value += delta.value
        delta.value = 0
        getHistory()
        return
      default:
        console.error('Ошибка пополнения счёта')
        break
    }
    return response.json();
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

const aboutNew = ref('');

function profileEdit() {
  fetch(`${runtimeConfig.public.apiBase}/me`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      about: aboutNew.value,
      first_name: userData.value.first_name,
      last_name: userData.value.last_name,
      alias: userData.value.alias
    })
  }).then((response) => {
    switch (response.status) {
      case 400:
        break;
      case 200:
        modalActive.value = false
        return response.json();
      default:
        break;
    }
  }).then((data) => {
    if (data) {
      userData.value.about = data.about
      aboutNew.value = ''
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function withdraw() {
  fetch(`${runtimeConfig.public.apiBase}/balance`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      "operation_type": "withdraw",
      "delta": balance.value
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        complete.value = true
        balance.value = 0
        getHistory()
        return
      default:
        console.error('Ошибка вывода')
        break
    }
    return response.json();
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

const reviews = ref([])

function getReviews() {
  fetch(`${runtimeConfig.public.apiBase}/me/reviews`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token.value}`
    }
  }).then((response) => {
    switch (response.status) {
      case 200:
        return response.json();
        break;
      default:
        console.error('Ошибка вывода');
        break
    }
  }).then((data) => {
    if (data) {
      if (role.value === 'contractor') {
        reviews.value = data.reviews.filter(item => item.review_type === 'seller')
      } else {
        reviews.value = data.reviews.filter(item => item.review_type === 'buyer')
      }
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

getReviews()


</script>

<template>
  <div class="main-content">
    <AwСheckmark :complete="complete" @toggleComplete="complete = !complete"></AwСheckmark>
    <h1 class="page-title">Личный кабинет </h1>
    <div class="container">
      <div class="column">
        <div class="block profile-info">
          <div class="row">
            <span v-if="userData.first_name+userData.last_name === ''">{{ userData.alias }}</span>
            <span v-else>{{ userData.first_name }} {{ userData.last_name }}</span>
            <button class="btn" style="float: right" @click="modalActive = !modalActive">Редактировать</button>
          </div>
          <span>Обо мне:</span>
          <textarea id="about" name="about" placeholder="Заполните информацио о себе" rows="6"
                    v-model="userData.about" disabled></textarea>
        </div>
        <div class="block balance">
          <div class="row">
            <span>Ваш баланс: {{ balance }} ₽</span>
            <button class="btn-primary" style="float: right" @click="withdraw">Вывести</button>
          </div>
          <div class="history-block">
            <span>История</span>
            <div class="history-card" v-for="storie in stories">
              <div class="row">
                <span>{{
                    storie.operation_type == 'deposit' ? 'Пополнение' : storie.operation_type == 'expense' ? 'Расход' : storie.operation_type == 'refund' ? 'Возврат' : storie.operation_type == 'withdraw' ? 'Снятие' : storie.operation_type ? 'Доход' : storie.operation_type == 'income'
                  }}</span>
                <span>{{ storie.delta }} ₽</span>
              </div>
              <div class="row">
                <span class="text-secondary">Операция</span>
                <span class="text-secondary">{{ unixToDate(storie.timestamp) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="column">
        <div class="block skills">
          <div class="row">
            <span>Навыки</span>
            <button class="btn" style="float: right" @click="modalActiveSkills =! modalActiveSkills">Редактировать
            </button>
          </div>
          <ul>
            <li class="skills" v-for="skill in skills">{{ skill.name }}</li>
          </ul>
        </div>
        <div class="block reviews">
          <div class="row">
            <span>Последние отзывы</span>
            <span style="float: right" v-if="role==='contractor'">Ваш рейтинг: {{ userData.seller_rating }}</span>
            <span style="float: right" v-if="role==='customer'">Ваш рейтинг: {{ userData.buyer_rating }}</span>
          </div>
          <div class="reviews-card" v-for="review in reviews">
            <span>{{ review.name }}</span>
            <p style="font-size: 20px">{{ review.text }}</p>
          </div>
        </div>
        <div class="block deposit">
          <h3>Внести деньги</h3>
          <form v-on:submit.prevent="deposit">
            <input type="number" placeholder="0 ₽" required minlength="0" v-model="delta">
            <button class="btn-primary">Внести</button>
          </form>
        </div>
      </div>
    </div>
  </div>
  <AwModal @close="modalActive = !modalActive" :modalActive="modalActive" :title="'Данные профиля'">
    <div class="modal-content">
      <form class="form" v-on:submit.prevent="profileEdit">
        <label for="">Имя</label>
        <input type="text" placeholder="Имя" v-model="userData.first_name">
        <label for="">Фамилия</label>
        <input type="text" placeholder="Фамилия" v-model="userData.last_name">
        <label for="">Псевдоним</label>
        <input type="text" placeholder="Псевдоним" v-model="userData.alias">
        <label for="about">Обо мне</label>
        <textarea id="about" name="about" placeholder="Расскажите о себе" rows="6"
                  v-model="aboutNew"></textarea>
        <button class="btn-primary">Сохранить</button>
      </form>
    </div>
  </AwModal>
  <AwModal @close="modalActiveSkills = !modalActiveSkills" :modalActive="modalActiveSkills" :title="'Навыки'">
    <div class="modal-content">
      <div class="group-column" style="margin-bottom: 10px">
        <label>Укажите навык</label>
        <div class="row" style="justify-content: flex-start; gap: 10px">
          <input placeholder="Навык" type="text" style="width: 40%"/>
          <button class="btn-primary">Добавить</button>
        </div>
      </div>
      <div class="skills">
        <span>Ваши навыки</span>
        <ul>
          <li class="skills-change" v-for="(skill, index) in skills" :key="skill.id">
            {{ skill.name }}
            <div class="icon-close" @click="skillsDelete(index)"></div>
          </li>
        </ul>
      </div>
    </div>
  </AwModal>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.deposit {
  button {
    margin-top: 10px;
  }
}

.balance {
  .history-block {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  & .history-card {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    flex-wrap: wrap;
    width: 100%;
    gap: 15px 0;

    & > span {
      flex: 0 2 40%;
      text-align: left;
      font-size: 20px;
    }
  }

  .text-secondary {
    font-size: 15px !important;
    color: $secondary-text-color;
  }
}

.skills {
  display: flex;
  flex-direction: column;
  gap: 20px;

  ul {
    display: flex;
    flex-wrap: wrap;
    gap: 10px 5px;

    li {
      list-style-type: none;
      border-radius: 30px;
      padding: 5px 15px;
      background-color: $primary-color;
      color: white;
      font-size: 16px;
      width: min-content;
    }

    .skills-change {
      display: flex;
      align-items: center;
      gap: 10px;

      & .icon-close {
        background: #ffffff;
      }
    }
  }
}

.reviews {
  gap: 20px !important;

  .reviews-card {
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding: 25px;
    border-radius: 10px;
    border: 1px solid $primary-color;
    background-color: transparent;
  }
}
</style>