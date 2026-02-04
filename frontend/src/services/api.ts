import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：处理错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api

// Auth API
export const authAPI = {
  register: (data: { name: string; email: string; password: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
  getMe: () => api.get('/auth/me'),
}

// Friends API
export const friendsAPI = {
  getFriends: () => api.get('/friends'),
  getFriendRequests: () => api.get('/friends/requests'),
  sendFriendRequest: (friendId: number) =>
    api.post('/friends/request', { friend_id: friendId }),
  acceptFriendRequest: (requestId: number) =>
    api.post(`/friends/requests/${requestId}/accept`),
  rejectFriendRequest: (requestId: number) =>
    api.post(`/friends/requests/${requestId}/reject`),
  deleteFriend: (friendId: number) => api.delete(`/friends/${friendId}`),
  blockUser: (userId: number, reason?: string) =>
    api.post(`/friends/${userId}/block`, { user_id: userId, reason }),
  unblockUser: (userId: number) => api.delete(`/friends/${userId}/unblock`),
  getBlockedUsers: () => api.get('/friends/blocked'),
}

// Users API
export const usersAPI = {
  getUser: (userId: number) => api.get(`/users/${userId}`),
  updateUser: (userId: number, data: any) => api.put(`/users/${userId}`, data),
}
