<script setup lang="ts">
const router = useRouter();
const currentPath = ref(router.currentRoute.value.path);
const runtimeConfig = useRuntimeConfig();
const token = useCookie('token');
const role = useCookie('role');

let menu = {
  state: ref(false),
  close: () => {
    menu.state.value = false
  }
}
let notificationShow = ref(false)

onMounted(() => {
  currentPath.value = router.currentRoute.value.path;
  document.addEventListener('click', menu.close)
});

onBeforeUnmount(() => {
  document.removeEventListener('click', menu.close)
})

watch(() => router.currentRoute.value.path, (newVal) => {
  getNotification()
});

function logout() {
  token.value = undefined;
  role.value = undefined;
  router.push({path: "/authorization"});
}

function changeRole() {
  role.value === 'contractor' ? role.value = 'customer' : role.value = 'contractor'
  window.location.reload()
}

const notifications = ref([])
const notificationsCount = ref(0)
const notificationsMark = ref(false)

function openNotification(close = false) {
  if (close) notificationShow.value = !notificationShow.value
  if (!notificationShow.value) {
    return
  } else {
    getNotification()
  }
}

function getNotification() {
  fetch(`${runtimeConfig.public.apiBase}/me/notifications`, {
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
        logout()
        console.log('Нет прав')
        break;
    }
  }).then((data) => {
    if (data) {
      // TODO было бы неплохо перевести на русский
      notifications.value = data.notifications
      notificationsCount.value = data.unread_count
      if (notificationsCount.value > 0) {
        notificationsMark.value = true
      } else
        notificationsMark.value = false
    }
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function readNotification(id) {
  fetch(`${runtimeConfig.public.apiBase}/me/notifications/${id}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${token.value}`
    },
    body: JSON.stringify({
      read: true
    })
  }).then((response) => {
    switch (response.status) {
      case 200:
        getNotification()
        return;
        break;
      case 403:
        logout()
        console.log('Нет прав')
        break;
    }
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

function hideNotification(id) {
  fetch(`${runtimeConfig.public.apiBase}/me/notifications/${id}`, {
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
        getNotification()
        return;
        break;
      case 403:
        logout()
        console.log('Нет прав')
        break;
    }
  }).then((data) => {
  }).catch((err) => {
    console.error("Невозможно отправить запрос", err);
  });
}

if (token.value) {
  getNotification()
}

</script>

<template>
  <header>
    <nav class="nav-menu">
      <ul>
        <NuxtLink to="/">
          <span style="font-size: 34px;">GoJob</span>
        </NuxtLink>
      </ul>
      <ul>
        <li class="link-group">
          <NuxtLink to="/"><span>Найти</span></NuxtLink>
          <NuxtLink to="/orders" v-if="role === 'customer'"><span>Мои заказы</span></NuxtLink>
          <NuxtLink to="/offers" v-if="role === 'contractor'"><span>Мои отклики</span></NuxtLink>
        </li>
        <li class="icon-group">
          <div class="icon-bell" v-if="role" v-on:click="openNotification(true)"></div>
          <span class="notify mark" v-if="notificationsMark"></span>
          <transition name="slide" v-if="role">
            <ul class="ul-menu notifications" v-if="notificationShow">
              <li v-for="notification in notifications">
                <section><span class="notify" v-if="notification.read === false"></span>
                  <p>{{ notification.message }}</p>
                </section>
                <section>
                  <div class="icon-eye" @click="hideNotification(notification.id)"></div>
                  <div class="icon-marked" v-if="!notification.read" @click="readNotification(notification.id)"></div>
                </section>
              </li>
            </ul>
          </transition>
          <div class="icon-profile" @click.stop="menu.state.value=!menu.state.value"></div>
          <transition name="slide" v-if="role">
            <ul class="ul-menu" v-if="menu.state.value">
              <li>
                <NuxtLink to="/profile" v-on:click="menu.state.value=!menu.state.value">Профиль <span class="role" v-if="role==='contractor'">Исполнитель</span><span
                    class="role" v-else>Заказчик</span></NuxtLink>
              </li>
              <li v-on:click="changeRole">Войти как <span class="change-role"
                                                          v-if="role==='customer'">исполнитель</span><span
                  class="change-role" v-else>заказчик</span>
              </li>
              <li v-on:click="logout">Выйти из аккаунта</li>
            </ul>
          </transition>
          <transition name="slide" v-if="!role">
            <ul class="ul-menu" v-if="menu.state.value">
              <li>
                <NuxtLink to="/authorization" v-on:click="menu.state.value=!menu.state.value">Войти</NuxtLink>
              </li>
            </ul>
          </transition>
        </li>
      </ul>
    </nav>
  </header>
</template>

<style scoped lang="scss">
@use '@/assets/scss/main.scss' as *;

.notify.mark {
  position: absolute;
  left: 27px;
  bottom: 4px;
  height: 8px;
  width: 8px;
}

.notify {
  background: #ff7615;
  height: 5px;
  width: 5px;
  border-radius: 50%;
}

.icon-eye {
  height: 20px;
}

header {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: min-content;
  background-color: $primary-color;
  color: white;
  padding: 10px 20px;
}

.nav-menu {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1300px;
  width: 100%;
  margin: auto;
}

.icon-group {
  display: flex;
  align-items: center;
  border-left: 1px solid rgb(0, 203, 144);
  padding-left: 20px;
  gap: 5px;
}

.link-group {
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 10px;
}

ul {
  display: flex;
  align-items: center;
  gap: 10px;
}

a {
  cursor: pointer;
  color: white;
}

.profile-menu {
  position: absolute;
  background: antiquewhite;
  z-index: 3545;
}

li.icon-profile {
  position: relative;
}

li.icon-profile {
  overflow: auto;
}

li.icon-group {
  position: relative;
}

.ul-menu {
  list-style: none;
  z-index: 1;
  position: absolute;
  width: 300px;
  top: 50px;
  right: -5px;
  background: #00A676;
  border-bottom-right-radius: 15px;
  border-bottom-left-radius: 15px;
  flex-direction: column;
  gap: 0;
  //padding: 25px 0;
  overflow: hidden;
}

.ul-menu li {
  padding: 10px;
  cursor: pointer;
}

.ul-menu li:hover {
  background: rgba(0, 0, 0, 0.36);
}

.ul-menu li, .ul-menu a {
  font-size: 17px;
  width: 100%;
}

.change-role {
  font-size: 17px;
}

.ul-menu .role {
  font-size: 15px;
  color: #c3fbcb;
}

ul.ul-menu.notifications {
  background: white;
  border: 1px solid #00A676;
  border-top: none;
  width: 440px;
}

.notifications p {
  font-size: 12pt;
  margin-left: 5px;
}

ul.ul-menu.notifications section {
  display: flex;
  align-items: center;
}

ul.ul-menu.notifications li {
  color: black;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

ul.ul-menu.notifications li:hover {
  background: #0000001f;
}

.slide-enter-active,
.slide-leave-active {
  transition: opacity 0.5s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
}
</style>