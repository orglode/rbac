import axios, { type AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResult } from '@/shared/types/api'
import { getToken, clearAuth } from '@/shared/utils/auth'
import router from '@/app/router'

const http = axios.create({
  baseURL: '/',
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    const payload = response.data as ApiResult
    if (payload && typeof payload.code === 'number' && payload.code !== 0) {
      ElMessage.error(payload.msg || '请求失败')
      return Promise.reject(payload)
    }
    return response
  },
  (error) => {
    const status = error.response?.status
    const msg = error.response?.data?.msg || error.message || '网络异常'
    if (status === 401) {
      clearAuth()
      if (router.currentRoute.value.path !== '/login') {
        router.replace({ path: '/login', query: { redirect: router.currentRoute.value.fullPath } })
      }
      ElMessage.error(msg || '请重新登录')
    } else {
      ElMessage.error(msg)
    }
    return Promise.reject(error)
  },
)

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const res = await http.request<ApiResult<T>>(config)
  return res.data.data
}

export default http
