<script setup lang="ts">
import AwMoneyInput from "~/components/ui-kit/AwMoneyInput.vue";
import AwRaitInput from "~/components/ui-kit/AwRaitInput.vue";

const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
definePageMeta({layout: 'default'})
useHead({
  title: 'Главная'
})
const orders = ref()
fetch(`${runtimeConfig.public.apiBase}/orders`, {
  method: 'GET',
}).then((response) => {
  return response.json();
}).then((data) => {
  if (data.orders) {
    orders.value = data.orders
    orders.value = orders.value.filter(item => item.status === 'available')
  }
}).catch((err) => {
  console.error("Невозможно отправить запрос", err);
});

const text = ref('')
const show = ref(false)

const sort = ref('none')

function getSearch() {
  text.value = text.value.replaceAll(' ', '')
  text.value = text.value.replaceAll(' ', '%20')

  let url;
  if (sort.value === 'none') {
    url = `${runtimeConfig.public.apiBase}/orders?q=${text.value}`
  } else {
    url = `${runtimeConfig.public.apiBase}/orders?q=${text.value}&sort=${sort.value}`
  }
  fetch(url, {
    method: 'GET',
  }).then((response) => {
    return response.json();
  }).then((data) => {
    if (data.orders) {
      show.value = data.orders.length === 0
      orders.value = data.orders
      orders.value = orders.value.filter(item => item.status === 'available')
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}
</script>

<template>
  <div class="main-content">
    <h1 class="page-title">Работа </h1>
    <div class="container">
      <div class="column" style="flex-basis: 60%;">
        <div class="search-block block">
          <div class="title">
            Поиск
          </div>
          <div class="search-group">
            <input placeholder="Ищу..." class="input" type="text" v-model="text" v-on:keyup.enter="getSearch">
            <button class="btn-primary" v-on:click="getSearch">Найти</button>
          </div>
        </div>
        <div class="orders">
          <p style="text-align: center" v-if="show">Заказы не найдены</p>
          <NuxtLink :to="'/order/'+order.order_id" class="order block" v-for="order in orders">
            <div class="row">
              <h3 class="order-title">
                {{ order.name }}
              </h3>
              <p class="price" style="float: right">
                <span v-if="order.price>0">{{ order.price }} ₽/заказ</span>
                <span v-else>Цена договорная</span>
              </p>
            </div>
            <p class="description">{{ order.description }}</p>
            <div class="row">
              <div class="viewers">
                <div class="icon-eye"></div>
                <span>{{ order.views }}</span>
              </div>
              <div class="rating">
                <div class="icon-rating"></div>
                <span>{{ order.rating }}</span>
              </div>
              <div class="offers">
                <div class="icon-offer"></div>
                <span>{{ order.offer_count }}</span>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>
      <div class="column" style="flex-basis: 20%;">
        <div class="block">
          <h4>Фильтр</h4>
          <form @change="getSearch">
            <section style="display: flex; margin-bottom: 15px;">
              <input type="radio" id="offer_count_desc" name="sort" value="offer_count_desc"
                     style="width: min-content; margin-right: 5px;" v-model="sort">
              <label for="offer_count_desc">Больше откликов</label>
            </section>
            <section style="display: flex; margin-bottom: 15px;">
              <input type="radio" id="offer_count_asc" name="sort" value="offer_count_asc"
                     style="width: min-content; margin-right: 5px;" v-model="sort">
              <label for="offer_count_asc">Меньше откликов</label>
            </section>
            <section style="display: flex; margin-bottom: 15px;">
              <input type="radio" id="budget_desc" name="sort" value="budget_desc"
                     style="width: min-content; margin-right: 5px;" v-model="sort">
              <label for="budget_desc">С большей стоимостью</label>
            </section>
            <section style="display: flex; margin-bottom: 15px;">
              <input type="radio" id="budget_asc" name="sort" value="budget_asc"
                     style="width: min-content; margin-right: 5px;" v-model="sort">
              <label for="budget_asc">С меньшей стоимостью</label>
            </section>
            <section style="display: flex; margin-bottom: 15px;">
              <input type="radio" id="none" name="sort" value="none"
                     style="width: min-content; margin-right: 5px;" v-model="sort">
              <label for="none">Без фильтра</label>
            </section>
          </form>
          <!--          <span>Бюджет</span>-->
          <!--          <div class="row">-->
          <!--            <span>От</span>-->
          <!--            <AwMoneyInput/>-->
          <!--            <span>до</span>-->
          <!--            <AwMoneyInput/>-->
          <!--          </div>-->
          <!--          <span>Рейтинг</span>-->
          <!--          <div class="row">-->
          <!--            <span>От</span>-->
          <!--            <AwRaitInput/>-->
          <!--            <span>до</span>-->
          <!--            <AwRaitInput/>-->
          <!--          </div>-->
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.search-group {
  display: flex;
  gap: 5px;
}

.viewers,
.rating,
.offers {
  display: flex;
  align-items: center;
}

.icon-offer {
  width: 25px !important;
  height: 25px !important;
}

.order .row:last-child {
  justify-content: flex-start
}

.order {
  margin-bottom: 15px;
}

.order:last-child {
  margin-bottom: 0;
}

a {
  color: black;
}

</style>