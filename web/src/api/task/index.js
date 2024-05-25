import { del, get, post, put } from '../axios/index.js'

export function taskGet(id) {
  return get(`/task/${id}`)
}

export function taskList() {
  return get('/task')
}

export function taskCreate(title, time) {
  return post(`/task`, { title, time })
}

export function taskUpdate(id, title, time) {
  return put(`/task/${id}`, { title, time })
}

export function taskDelete(id) {
  return del(`/task/${id}`)
}

export function taskGetByDate(d) {
  return get(`/task/date?d=${d}`)
}