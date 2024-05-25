<template>
  <div id="login">
    <a-typography-title :heading="2">Login</a-typography-title>
    <a-form :model="form" :style="{ width: '600px' }" @submit="login" layout="vertical">
      <a-form-item field="email" tooltip="Please enter email" label="Email">
        <a-input
          v-model="form.email"
          placeholder="please enter your email..."
        />
      </a-form-item>
      <a-form-item field="password" label="Password">
        <a-input-password v-model="form.password" placeholder="please enter your password..." />
      </a-form-item>
      <a-typography-text>
        Don't have account, click here to
        <a-link href="/sign-up">Sign Up</a-link>
        .
      </a-typography-text>
      <a-link></a-link>
      <a-form-item>
        <a-button html-type="submit" type="primary" long>Login</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { userLogin } from '@/api/user/index.js'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})

const login = () => {
  userLogin(form.email, form.password).then(res => {
    if (res['status_code'] !== 0) {
      Message.error(res['status_msg'])
      return
    }
    localStorage.setItem('userinfo', JSON.stringify(res.data))
    router.push('/')
  })
}
</script>

<style scoped lang="less">
#login {
  display: flex;
  flex-direction: column;
  padding-top: 200px;
  align-items: center;
  gap: 20px;
}
</style>