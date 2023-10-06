<script setup lang="ts">
let chatsList = ref([
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
  {name: "Adam"},
  {name: "Заказчик"},
])
let messageList = ref([
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {
    name: "Adam",
    message: "Как у тебя дела?Как у тебя дела?Как у тебя дела?Как у тебя дела?Как у тебя дела?Как у тебя дела?"
  },
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
  {name: "Adam", message: "Привет, друг"},
  {name: "Adam", message: "Как у тебя дела?"},
])


onBeforeMount(async () => {
  await nextTick();
  const messageWindow = ref()
  messageWindow.value.scrollTo({bottom: messageWindow.value.offsetHeight})
})

</script>

<template>
  <div class="main-content">
    <div class="container">
      <div class="column chat-list">
        <span>Чаты</span>
        <div class="chat-elem" v-for="(item, index) in chatsList">
          <span>{{ item.name }}</span>
        </div>
      </div>
      <div class="column" style="flex-basis: 40%">
        <div class="message-window" ref="messageWindow">
          <div class="message-window-content" style="overflow-y: scroll">
            <div class="date"></div>
            //СДЕЛАТЬ ТАК, ЧТОБЫ СООБЩЕНИЯ БЫЛИ В ОДНОМ ОКНЕ, СОРИТИРОВАЛИСЬ ПО ДАТЕ И СТИЛИ ЗАВИСИЛИ ОТ ОТПРАВИТЕЛЯ,
            //ТО ЕСТЬ, ЕСЛИ МОЁ СООБЩЕНИЕ - ТО ОНО СПРАВО.
            <div class="message-other" v-for="(item, index) in messageList">
              <span class="name-reciver">{{ item.name }}</span>
              <span style="transform: translateY(-15px);">{{ item.message }}</span>
            </div>
            <div style="display: flex; flex-direction: column;
align-items: flex-end; width: 100%">
              <div class="message-my">
                <span>fdfkq</span>
              </div>
            </div>
          </div>
          <div>
            <div class="row">
              <input type="text">
              <button class="btn-primary">отправить</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.main-content {
  max-width: 90% !important;
}

.container {
  position: relative;

  & .column {
    min-height: 100%;
  }
}

.chat-elem {
  display: flex;
  flex-wrap: wrap;
  background-color: #00a47f;
  border-radius: 15px;
  color: white;
  padding: 30px 20px;
}

.message-window {
  position: fixed;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 90vh;
  min-height: 85vh;
  max-height: 85vh;
}

.message-window-content {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  gap: 40px;
}

.message-other, .message-my{
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 5px;
  max-width: 60%;
  background-color: #00a47f;
  border-radius: 30px;
  padding: 5px 30px;
  color: white;

  & .name-reciver {
    background-color: #00a47f;
    width: min-content;
    border-radius: 30px;
    padding: 0 10px;
    font-size: 18px;
    position: relative;
    color: #cccaca;
    transform: translateY(-20px);
  }
}

.message-my {

}
</style>