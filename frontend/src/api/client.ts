import axios from 'axios'

// === B-End API Client ===
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach B-end JWT
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('b_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('b_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

// === C-End API Client ===
export const cApi = axios.create({
  baseURL: import.meta.env.VITE_CEND_API_BASE_URL || 'http://localhost:8080',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach C-end wallet JWT
cApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('c_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

cApi.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('c_token')
    }
    return Promise.reject(error)
  },
)
