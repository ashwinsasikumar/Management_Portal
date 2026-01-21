import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'
import { API_BASE_URL } from '../config'

function Dashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState({
    totalCurriculum: 0,
    activeCurriculum: 0,
    totalCourses: 0,
    recentActivities: 0
  })

  useEffect(() => {
    // Fetch dashboard stats
    fetchDashboardStats()
  }, [])

  const fetchDashboardStats = async () => {
    try {
      // Fetch actual curriculum count from API
      const response = await fetch(`${API_BASE_URL}/curriculum`)
      if (response.ok) {
        const data = await response.json()
        setStats({
          totalCurriculum: data.length || 0,
          activeCurriculum: data.length || 0,
          totalCourses: 0, // Can be calculated from all courses across semesters
          recentActivities: 0 // Can be fetched from logs API
        })
      } else {
        // Fallback to default values if API fails
        setStats({
          totalCurriculum: 0,
          activeCurriculum: 0,
          totalCourses: 0,
          recentActivities: 0
        })
      }
    } catch (error) {
      console.error('Error fetching dashboard stats:', error)
      setStats({
        totalCurriculum: 0,
        activeCurriculum: 0,
        totalCourses: 0,
        recentActivities: 0
      })
    }
  }

  const statCards = [
    {
      title: 'Total Curriculum',
      value: stats.totalCurriculum,
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      ),
      color: 'from-blue-500 to-blue-600',
      bgColor: 'bg-blue-50',
      textColor: 'text-blue-600',
      customColor: 'rgb(255, 195, 0)',
      customBg: 'rgba(255, 195, 0, 0.1)'
    },
    {
      title: 'Active Curriculum',
      value: stats.activeCurriculum,
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
      ),
      color: 'from-green-500 to-green-600',
      bgColor: 'bg-green-50',
      textColor: 'text-green-600'
    },
    {
      title: 'Total Courses',
      value: stats.totalCourses,
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      ),
      color: 'from-purple-500 to-purple-600',
      bgColor: 'bg-purple-50',
      textColor: 'text-purple-600'
    },
    {
      title: 'Recent Activities',
      value: stats.recentActivities,
      icon: (
        <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
      color: 'from-orange-500 to-orange-600',
      bgColor: 'bg-orange-50',
      textColor: 'text-orange-600'
    }
  ]

  const quickActions = [
    // {
    //   title: 'Manage Regulations',
    //   description: 'View and manage all regulations',
    //   icon: (
    //     <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
    //       <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
    //     </svg>
    //   ),
    //   action: () => navigate('/regulations')
    // },
    {
      title: 'View Curriculum',
      description: 'Browse all curriculum structures',
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
      ),
      action: () => navigate('/curriculum')
    },
    {
      title: 'Manage Clusters',
      description: 'Create and manage department clusters',
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      ),
      action: () => navigate('/clusters')
    },
    {
      title: 'Manage Sharing',
      description: 'Control content sharing between cluster departments',
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      ),
      action: () => navigate('/sharing')
    }
  ]

  // Add Users Management action for admin users only
  const userRole = localStorage.getItem('userRole')
  if (userRole === 'admin') {
    quickActions.push({
      title: 'User Management',
      description: 'Manage system users and permissions',
      icon: (
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
      ),
      action: () => navigate('/users')
    })
  }

  return (
    <MainLayout 
      title="Dashboard" 
      subtitle="Welcome back! Here's what's happening with your curriculum"
    >
      <div className="space-y-8">
        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {statCards.map((stat, index) => (
            <div key={index} className="card-custom p-6 hover:scale-105 transition-transform duration-200">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-sm font-medium text-gray-600 mb-1">{stat.title}</p>
                  <p className="text-3xl font-bold text-gray-900">{stat.value}</p>
                </div>
                <div className={stat.customBg ? 'p-3 rounded-xl' : `${stat.bgColor} p-3 rounded-xl`} style={stat.customBg ? {backgroundColor: stat.customBg} : {}}>
                  <div className={stat.customColor ? '' : stat.textColor} style={stat.customColor ? {color: stat.customColor} : {}}>
                    {stat.icon}
                  </div>
                </div>
              </div>
              <div className="mt-4 flex items-center">
                <div className={stat.customColor ? 'w-full h-1 rounded-full' : `w-full h-1 bg-gradient-to-r ${stat.color} rounded-full`} style={stat.customColor ? {backgroundColor: stat.customColor} : {}}></div>
              </div>
            </div>
          ))}
        </div>

        {/* Quick Actions */}
        <div className="card-custom p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-6">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {quickActions.map((action, index) => (
              <button
                key={index}
                onClick={action.action}
                className="flex items-start space-x-4 p-5 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl transition-all duration-200 hover:scale-105 border border-gray-200 group"
                style={{
                  '--hover-from': 'rgba(67, 113, 229, 0.05)',
                  '--hover-to': 'rgba(67, 113, 229, 0.1)'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'linear-gradient(to bottom right, rgba(67, 113, 229, 0.05), rgba(67, 113, 229, 0.1))'
                  e.currentTarget.style.borderColor = 'rgba(67, 113, 229, 0.3)'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'linear-gradient(to bottom right, rgb(249, 250, 251), rgb(243, 244, 246))'
                  e.currentTarget.style.borderColor = 'rgb(229, 231, 235)'
                }}
              >
                <div className="flex-shrink-0 w-12 h-12 bg-white rounded-lg flex items-center justify-center transition-all duration-200 shadow-sm group-hover:text-white" style={{color: 'rgb(67, 113, 229)'}} onMouseEnter={(e) => {
                  const parent = e.currentTarget.parentElement
                  if (parent?.matches(':hover')) {
                    e.currentTarget.style.background = 'rgb(67, 113, 229)'
                    e.currentTarget.style.color = 'white'
                  }
                }} onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'white'
                  e.currentTarget.style.color = 'rgb(67, 113, 229)'
                }}>
                  {action.icon}
                </div>
                <div className="flex-1 text-left">
                  <h3 className="text-base font-semibold text-gray-900 mb-1">{action.title}</h3>
                  <p className="text-sm text-gray-600">{action.description}</p>
                </div>
                <svg className="w-5 h-5 text-gray-400 transition-colors group-hover:text-[rgb(67,113,229)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            ))}
          </div>
        </div>

        {/* Welcome Message */}
        <div className="card-custom p-8 text-white" style={{background: 'linear-gradient(to bottom right, rgb(67, 113, 229), rgb(47, 93, 209))'}}>
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <h2 className="text-2xl font-bold mb-3">Welcome to Curriculum Management System</h2>
              <p className="mb-6 max-w-2xl" style={{color: 'rgba(255, 255, 255, 0.9)'}}>
                Streamline your academic planning with our comprehensive curriculum management platform. 
                Create, manage, and track curriculum structures, courses, and mappings all in one place.
              </p>
              <button
                onClick={() => navigate('/curriculum')}
                className="bg-white px-6 py-3 rounded-lg font-semibold hover:shadow-lg transition-all duration-200 hover:scale-105 active:scale-95"
                style={{color: 'rgb(67, 113, 229)'}}
              >
                Get Started
              </button>
            </div>
            <div className="hidden lg:block">
              <svg className="w-32 h-32 opacity-50" style={{color: 'rgba(255, 255, 255, 0.4)'}} fill="currentColor" viewBox="0 0 20 20">
                <path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  )
}

export default Dashboard
