import { useState, useEffect } from 'react'

export interface User {
  id: number
  name: string
  email: string
  username?: string
  avatar_url?: string
  is_vip: boolean
  is_online: boolean
}

export function useAuth() {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    if (token) {
      // 从 localStorage 加载用户信息
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        setUser(JSON.parse(storedUser))
      }
    }
  }, [token])

  const login = (newToken: string, userData: User) => {
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(userData))
    setToken(newToken)
    setUser(userData)
  }

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }

  return {
    token,
    user,
    login,
    logout,
    isAuthenticated: !!token,
  }
}
