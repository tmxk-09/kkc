import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 300000, // 5分钟超时（大文件上传）
})

// 请求拦截器 — 自动加 Token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器
api.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code === 401) {
      localStorage.removeItem('token')
      window.location.hash = '#/login'
      return Promise.reject(new Error('未登录'))
    }
    return res
  },
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.hash = '#/login'
    }
    return Promise.reject(error)
  }
)

// ========== 认证 ==========
export function login(data) {
  return api.post('/auth/login', data)
}

export function register(data) {
  return api.post('/auth/register', data)
}

// ========== 上传 ==========
export function checkHash(fileHash) {
  return api.post('/upload/check', { fileHash })
}

export function instantUpload(data) {
  return api.post('/upload/instant', data)
}

export function initUpload(data) {
  return api.post('/upload/init', data)
}

export function uploadChunk(uploadId, chunkIndex, file, onProgress) {
  const formData = new FormData()
  formData.append('uploadId', uploadId)
  formData.append('chunkIndex', chunkIndex)
  formData.append('file', file)
  return api.post('/upload/chunk', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress,
    timeout: 600000, // 10分钟
  })
}

/**
 * 带实时进度回调的分片上传
 * @param {Function} onProgressBytes - (loaded: number) => void 已上传字节数
 */
export function uploadChunkWithProgress(uploadId, chunkIndex, file, onProgressBytes) {
  const formData = new FormData()
  formData.append('uploadId', uploadId)
  formData.append('chunkIndex', chunkIndex)
  formData.append('file', file)
  return api.post('/upload/chunk', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (onProgressBytes && e.loaded) {
        onProgressBytes(e.loaded)
      }
    },
    timeout: 600000,
  })
}

export function completeUpload(uploadId, expireDays = 1) {
  return api.post('/upload/complete', { uploadId, expireDays })
}

export function getProgress(uploadId) {
  return api.post('/upload/progress', { uploadId })
}

export function getUploadTasks() {
  return api.post('/upload/tasks')
}

export function deleteUploadTask(uploadId) {
  return api.post('/upload/task/delete', { uploadId })
}

// ========== 文件管理 ==========
export function getFileList(data) {
  return api.post('/file/list', data)
}

export function deleteFile(id) {
  return api.post('/file/delete', { id })
}

export function getDownloadUrl(id) {
  return `/api/file/download/${id}`
}

export function checkPickupCode(code) {
  return api.get(`/file/pickup/check/${code}`)
}

export function getPickupUrl(code) {
  return `/api/file/pickup/${code}`
}

// ========== 用户管理 (Admin & User) ==========
export function getUserInfo() {
  return api.get('/user/info')
}
export function getUserList() {
  return api.get('/user/list')
}
export function updateUserStatus(data) {
  return api.post('/user/status', data)
}
export function updateUserSpace(data) {
  return api.post('/user/space', data)
}
export function updateUserPassword(data) {
  return api.post('/user/password', data)
}

// ========== 空间赠送 (数字资产流转) ==========
export function giftSpace(data) {
  return api.post('/space/gift', data)
}

// ========== 工单管理 ==========
export function createTicket(data) {
  return api.post('/ticket/create', data)
}
export function getTicketList(data) {
  return api.post('/ticket/list', data || {})
}
export function replyTicket(data) {
  return api.post('/ticket/reply', data)
}

// ========== 系统通知 ==========
export function getMessageList() {
  return api.post('/message/list')
}
export function readMessage(id) {
  return api.post('/message/read', { id })
}
export function readAllMessage() {
  return api.post('/message/read_all')
}
export function sendMessage(data) {
  return api.post('/message/send', data)
}

export default api
