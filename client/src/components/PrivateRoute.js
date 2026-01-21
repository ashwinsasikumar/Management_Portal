import React from 'react'
import { Navigate } from 'react-router-dom'

const PrivateRoute = ({ children }) => {
  const isAuthenticated = () => {
    // Check if user is authenticated by verifying localStorage
    const userId = localStorage.getItem('userId')
    const userRole = localStorage.getItem('userRole')
    return userId && userRole
  }

  return isAuthenticated() ? children : <Navigate to="/" replace />
}

export default PrivateRoute
