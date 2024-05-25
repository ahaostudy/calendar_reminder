import axios from 'axios'
import { Message } from '@arco-design/web-vue'
import router from '../../router'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 60000
})

instance.interceptors.request.use(
  (config) => {
    let token = JSON.parse(localStorage.getItem('userinfo'))?.token
    if (!token) token = ''

    config.headers.Authorization = 'Bearer ' + token
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

instance.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (!error.response) {
      return Promise.reject(error)
    }
    if (error.response.status === 401) {
      Message.error('Please login your account first')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default instance