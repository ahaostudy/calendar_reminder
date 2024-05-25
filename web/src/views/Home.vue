<template>
  <div id="home">
    <a-page-header
      :style="{ background: 'var(--color-bg-2)' }"
      title="Calendar Reminder"
    >
      <template #back-icon>
        <Logo />
      </template>
      <template #subtitle>
        <a-space>
          <span>Calendar Management Assistant</span>
        </a-space>
      </template>
      <template #extra>
        <div class="home-header-extra">
          <a-typography-text v-if="loggedIn">{{ userinfo['email'] }}</a-typography-text>
          <a-button href="/login" type="text" v-if="!loggedIn">Login</a-button>
          <a-button @click="logout" type="text" v-else>Logout</a-button>
        </div>
      </template>

      <div class="main">
        <div class="calendar">
          <a-typography-title>{{ dateString(calendarDate) }}</a-typography-title>
          <a-calendar v-model="calendarDate" style="width: 100%;" />
        </div>

        <div class="task-list">
          <TaskList :date="calendarDate" />
        </div>
      </div>

    </a-page-header>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import TaskList from '@c/page/home/TaskList.vue'
import { dateString } from '@/utils/date.js'
import Logo from '@c/icon/Logo.vue'

const calendarDate = ref(new Date())

const loggedIn = ref(false)
const userinfo = ref({ email: '' })

onMounted(() => {
  let ui = JSON.parse(localStorage.getItem('userinfo'))
  if (ui !== null) {
    loggedIn.value = true
    userinfo.value = ui['user']
  }
})

function logout() {
  localStorage.removeItem('userinfo')
  loggedIn.value = false
  userinfo.value.email = ''
}
</script>

<style scoped lang="less">
#home {
  height: 100vh;
}

.home-header-extra {
  display: flex;
  gap: 20px;
  align-items: center;
}

.main {
  display: flex;
  height: calc(100vh - 61px);

  .calendar {
    padding: 0 50px;
    width: 35%;
    min-width: 500px;
    max-width: 800px;

    display: flex;
    flex-direction: column;
    align-items: end;
  }

  .task-list {
    padding: 0 50px 50px;
    flex: 1;
    background: #f8f9fa;
    overflow: auto;
  }
}
</style>

<style>
h1.arco-typography {
  margin: 40px 0;
}

.arco-page-header-content {
  padding: 0;
}

.arco-page-header {
  padding-bottom: 0;
}
</style>