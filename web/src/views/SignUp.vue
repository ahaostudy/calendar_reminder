<template>
  <div id="sign-up">
    <a-typography-title :heading="2">Sign Up</a-typography-title>
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
      <a-form-item field="password_confirm" label="Confirm Password">
        <a-input-password v-model="form.password_confirm" placeholder="please enter your confirm password..." />
      </a-form-item>
      <a-typography-text>
        Already have an account, click here to
        Don't have account, click here to
        <a-link href="/login">Login</a-link>
        .
      </a-typography-text>
      <a-link></a-link>
      <a-form-item>
        <a-button html-type="submit" type="primary" long>Sign Up</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { userRegister } from '@/api/user/index.js'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: '',
  password_confirm: ''
})

const login = () => {
  userRegister(form.email, form.password, form.password_confirm).then(res => {
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
#sign-up {
  display: flex;
  flex-direction: column;
  padding-top: 200px;
  align-items: center;
  gap: 20px;
}
</style>
