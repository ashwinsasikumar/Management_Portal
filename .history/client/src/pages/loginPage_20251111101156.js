import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const handleLogin = (e) => {
    e.preventDefault()
    
    // Check credentials
    if (username === 'admin' && password === '987654321') {
      // Clear form and navigate to dashboard
      setUsername('')
      setPassword('')
      setError('')
      navigate('/dashboard')
    } else {
      setError('Invalid username or password')
    }
  }

  return (
    <div className="flex justify-center items-center min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600">
      <div className="bg-white p-10 rounded-2xl shadow-2xl w-full max-w-md">
        <h1 className="text-center text-gray-800 mb-8 text-3xl font-bold">Login</h1>
        
        <form onSubmit={handleLogin}>
          <div className="mb-5">
            <label 
              htmlFor="username" 
              className="block mb-2 text-gray-700 font-medium text-sm"
            >
              Username
            </label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter username"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-lg text-sm transition-all focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          <div className="mb-5">
            <label 
              htmlFor="password" 
              className="block mb-2 text-gray-700 font-medium text-sm"
            >
              Password
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter password"
              required
              className="w-full px-4 py-3 border border-gray-300 rounded-lg text-sm transition-all focus:outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
            />
          </div>

          {error && (
            <div className="text-red-600 text-sm mb-5 p-3 bg-red-50 border-l-4 border-red-600 rounded">
              {error}
            </div>
          )}

          <button 
            type="submit" 
            className="w-full py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white rounded-lg text-base font-semibold cursor-pointer transition-all hover:-translate-y-1 hover:shadow-xl"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  )
}

export default LoginPage