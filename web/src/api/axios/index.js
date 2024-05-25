import instance from './instance.js'

export async function get(url, params) {
  const response = await instance.get(url, { params })
  return response.data
}

export async function post(url, data, config) {
  const response = await instance.post(url, data, config)
  return response.data
}

export async function put(url, data, config) {
  const response = await instance.put(url, data, config)
  return response.data
}

export async function del(url, params) {
  const response = await instance.delete(url, { params })
  return response.data
}