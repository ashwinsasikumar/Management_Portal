import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'

function CurriculumPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [semesters, setSemesters] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [newSemester, setNewSemester] = useState({ semester_number: '' })

  useEffect(() => {
    fetchSemesters()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const fetchSemesters = async () => {
    try {
      setLoading(true)
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semesters`)
      if (!response.ok) {
        throw new Error('Failed to fetch semesters')
      }
      const data = await response.json()
      setSemesters(data || [])
      setError('')
    } catch (err) {
      console.error('Error fetching semesters:', err)
      setError('Failed to load semesters')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateSemester = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semester`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newSemester),
      })

      if (!response.ok) {
        throw new Error('Failed to create semester')
      }

      setNewSemester({ semester_number: '' })
      setShowCreateForm(false)
      fetchSemesters()
    } catch (err) {
      console.error('Error creating semester:', err)
      setError('Failed to create semester')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600 flex items-center justify-center">
        <div className="text-white text-xl font-medium">Loading curriculum...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600 p-6">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-6 bg-white/95 backdrop-blur-md px-8 py-6 rounded-2xl shadow-xl">
          <div>
            <button
              onClick={() => navigate('/regulations')}
              className="text-indigo-600 hover:text-indigo-800 font-medium mb-2 flex items-center gap-2"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M19 12H5M12 19l-7-7 7-7" />
              </svg>
              Back to Regulations
            </button>
            <h1 className="text-3xl font-bold text-gray-800">Curriculum Builder</h1>
            <p className="text-gray-600 mt-1">Regulation ID: {id}</p>
          </div>
          <button
            onClick={() => setShowCreateForm(!showCreateForm)}
            className="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all"
          >
            {showCreateForm ? 'Cancel' : '+ Create Semester'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {/* Create Semester Form */}
        {showCreateForm && (
          <div className="mb-6 bg-white/95 backdrop-blur-md p-6 rounded-2xl shadow-xl">
            <form onSubmit={handleCreateSemester} className="flex gap-4 items-end">
              <div className="flex-1">
                <label className="block text-gray-700 font-semibold mb-2 text-sm">Semester Number</label>
                <input
                  type="number"
                  value={newSemester.semester_number}
                  onChange={(e) => setNewSemester({ semester_number: e.target.value })}
                  placeholder="e.g., 1"
                  required
                  min="1"
                  className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </div>
              <button
                type="submit"
                className="px-6 py-2.5 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg hover:shadow-lg transition-all"
              >
                Create Semester
              </button>
            </form>
          </div>
        )}

        {/* Semesters Grid */}
        {semesters.length === 0 ? (
          <div className="text-center py-16 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl text-gray-600 text-lg">
            No semesters found. Create one to get started!
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {semesters.map(sem => (
              <div
                key={sem.id}
                onClick={() => navigate(`/regulation/${id}/curriculum/semester/${sem.id}`)}
                className="group bg-white/95 backdrop-blur-md rounded-xl shadow-md hover:shadow-xl p-6 transition-all duration-300 hover:-translate-y-1 cursor-pointer border border-white/50 hover:border-indigo-300"
              >
                <div className="text-center">
                  <div className="text-5xl font-bold text-indigo-600 mb-2">{sem.semester_number}</div>
                  <h3 className="text-lg font-semibold text-gray-800">Semester {sem.semester_number}</h3>
                  <p className="text-sm text-gray-500 mt-2">Click to manage courses</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default CurriculumPage
