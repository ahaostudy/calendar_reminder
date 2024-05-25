import { post } from '../axios/index.js'

export function userRegister(email, password, password_confirm) {
  return post('/user/register', { email, password, password_confirm })
}

export function userLogin(email, password) {
  return post('/user/login', { email, password })
}
