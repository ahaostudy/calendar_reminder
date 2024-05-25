<template>
  <div id="task-list">
    <div class="task-list-header">
      <a-typography-title :heading="6" :class="'title'">{{ date.toDateString() }} Tasks (We will notify you via your
        email)
      </a-typography-title>
      <a-button @click="showNewTask">
        <template #icon>
          <icon-plus />
        </template>
        Add Task
      </a-button>
    </div>
    <div class="task-list">
      <div v-if="taskList.length === 0">
        <a-empty></a-empty>
      </div>
      <div class="task-item" v-for="(task, i) in taskList" :key="task.id" @click="showTask(i)">
        <IconClock v-if="task.time * 1000 > now.getTime()"></IconClock>
        <IconTask v-else></IconTask>
        <div class="task-item-content">
          <a-typography-text class="task-item-title">
            {{ task.title }}
            <span class="task-item-time">{{ new Date(task.time * 1000).toLocaleString() }}</span>
          </a-typography-text>
          <a-button type="text" shape="circle" @click.stop="deleteTask(task.id, i)">
            <template #icon>
              <icon-delete size="16" />
            </template>
          </a-button>
        </div>
      </div>
    </div>

    <a-modal class="task-modal" title="Task" v-model:visible="visibleTaskModal" @ok="saveTask">
      <a-form :model="modalTask" :layout="'vertical'">
        <a-form-item field="name" label="Title">
          <a-input v-model="modalTask.title" placeholder="please enter your task title..." />
        </a-form-item>
        <a-form-item field="time" label="Time">
          <a-date-picker v-model="modalTask.time" showTime />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { taskDelete, taskGetByDate } from '@/api/task'
import { Message } from '@arco-design/web-vue'
import { dateString } from '@/utils/date.js'
import IconClock from '@c/icon/IconClock.vue'
import IconTask from '@c/icon/IconTask.vue'
import { taskCreate, taskUpdate } from '../../../api/task/index.js'

const props = defineProps({
  date: {
    default: new Date(),
    type: Date
  }
})

const taskList = reactive([])
const now = ref(new Date())

function viewDate(date) {
  taskGetByDate(dateString(date)).then(res => {
    if (res['status_code'] !== 0) {
      Message.error(res['status_msg'])
      return
    }
    taskList.splice(0, taskList.length, ...res.data)
  })
}

onMounted(() => {
  viewDate(now.value)
})
watch(() => props.date, viewDate)

const visibleTaskModal = ref(false)
const modalTask = ref({})

function showTask(i) {
  let task = taskList[i]
  visibleTaskModal.value = true
  modalTask.value.i = i
  modalTask.value.id = task.id
  modalTask.value.title = task.title
  modalTask.value.time = task.time * 1000
}

function showNewTask() {
  visibleTaskModal.value = true
  modalTask.value.i = -1
  modalTask.value.id = 0
  modalTask.value.title = ''
  modalTask.value.time = Math.max(props.date.getTime(), now.value.getTime() + 60 * 60 * 1000)
}

function saveTask() {
  const mt = modalTask.value
  const time = new Date(mt.time).getTime() / 1000
  if (mt.i === -1) createTask(mt.title, time)
  else updateTask(mt.i, mt.id, mt.title, time)
}


function createTask(title, time) {
  taskCreate(title, time).then(res => {
    if (res['status_code'] !== 0) {
      Message.error(res['status_msg'])
      return
    }
    taskList.push(res.data)
    Message.success('create task successful')
  })
}

function updateTask(i, id, title, time) {
  taskUpdate(id, title, time).then(res => {
    if (res['status_code'] !== 0) {
      Message.error(res['status_msg'])
      return
    }
    taskList[i] = res.data
    Message.success('update task successful')
  })
}

function deleteTask(id, i) {
  taskDelete(id).then(res => {
    if (res['status_code'] !== 0) {
      Message.error(res['status_msg'])
      return
    }
    taskList.splice(i, 1)
    Message.success('the task has been deleted')
  })
}

setInterval(() => {
  now.value = new Date()
}, 1000)
</script>

<style scoped lang="less">
#task-list {
  display: flex;
  flex-direction: column;

  .task-list-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .title {
      margin: 40px 0;
      color: #616468;
    }
  }

  .task-list {
    display: flex;
    flex-direction: column;
    background: #fff;
    border-radius: 10px;
    padding: 12px;
    min-width: 200px;

    .task-item {
      display: flex;
      align-items: center;
      padding: 0 15px 0 25px;
      height: 72px;
      gap: 35px;
      border-radius: 12px;
      cursor: pointer;

      .task-item-content {
        border-top: #dededf 1px solid;
        flex: 1;
        height: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .task-item-time {
          margin-left: 5px;
          color: #666a73;
          font-size: 13px;
        }
      }
    }

    .task-item:hover {
      background-color: #eeeeee;
      border: none;

      .task-item-content {
        border: none;
      }
    }

    .task-item:hover + .task-item {
      .task-item-content {
        border: none;
      }
    }

    .task-item:first-child {
      .task-item-content {
        border: none;
      }
    }
  }
}
</style>

<style>
.task-modal .arco-modal {
  max-width: 90%;
  border-radius: 10px;

  .arco-modal-header {
    border-bottom: none;
    border-top-right-radius: 12px;
    border-top-left-radius: 12px;
    background-color: #fbfbfc;
  }

  .arco-modal-body {
    padding: 24px 36px 4px 36px;
  }

  .arco-modal-footer {
    border-top: none;

    .arco-btn {
      border-radius: 6px;
      content: 'hello';
    }
  }

  .modal-title {
    width: 100%;
  }

}
</style>