import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function ManageCurriculumPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [semesters, setSemesters] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [newSemester, setNewSemester] = useState({ semester_number: '' })
  const [showEditModal, setShowEditModal] = useState(false)
  const [editingSemester, setEditingSemester] = useState(null)
  const [editSemesterNumber, setEditSemesterNumber] = useState('')

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
        body: JSON.stringify({
          semester_number: parseInt(newSemester.semester_number)
        }),
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

  const handleEditSemester = (semester) => {
    setEditingSemester(semester)
    setEditSemesterNumber(semester.semester_number.toString())
    setShowEditModal(true)
  }

  const handleUpdateSemester = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/semester/${editingSemester.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          semester_number: parseInt(editSemesterNumber)
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to update semester')
      }

      setShowEditModal(false)
      setEditingSemester(null)
      fetchSemesters()
    } catch (err) {
      console.error('Error updating semester:', err)
      setError('Failed to update semester')
    }
  }

  if (loading) {
    return (
      <MainLayout title="Manage Curriculum" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading curriculum...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout 
      title="Manage Curriculum"
      subtitle={`Semesters for Regulation ID: ${id}`}
      actions={
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate('/curriculum')}
            className="btn-secondary-custom flex items-center space-x-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span>Back</span>
          </button>
          <button
            onClick={() => setShowCreateForm(!showCreateForm)}
            className="btn-primary-custom flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            <span>{showCreateForm ? 'Cancel' : 'Add Semester'}</span>
          </button>
        </div>
      }
    >
      <div className="max-w-6xl mx-auto space-y-6">

        {/* Error Message */}
        {error && (
          <div className="flex items-start space-x-3 p-4 bg-red-50 border border-red-200 rounded-lg">
            <svg className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <p className="text-sm font-medium text-red-600">{error}</p>
          </div>
        )}

        {/* Create Semester Form */}
        {showCreateForm && (
          <div className="card-custom p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Add New Semester</h2>
            <form onSubmit={handleCreateSemester} className="flex gap-4 items-end">
              <div className="flex-1">
                <label className="block text-sm font-semibold text-gray-700 mb-2">Semester Number</label>
                <input
                  type="number"
                  value={newSemester.semester_number}
                  onChange={(e) => setNewSemester({ semester_number: e.target.value })}
                  placeholder="e.g., 1"
                  required
                  min="1"
                  className="input-custom"
                />
              </div>
              <button type="submit" className="btn-primary-custom">Create Semester</button>
            </form>
          </div>
        )}

        {/* Semesters Grid */}
        {semesters.length === 0 ? (
          <div className="card-custom p-12 text-center">
            <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Semesters Yet</h3>
            <p className="text-gray-600 mb-6">Get started by creating your first semester</p>
            <button onClick={() => setShowCreateForm(true)} className="btn-primary-custom">Add Semester</button>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
            {semesters.map(sem => (
              <div
                key={sem.id}
                className="group card-custom p-6 cursor-pointer hover:scale-105 transition-all duration-200 relative"
                onClick={() => navigate(`/regulation/${id}/curriculum/semester/${sem.id}`)}
              >
                {/* Edit Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    handleEditSemester(sem)
                  }}
                  className="absolute top-3 right-3 w-8 h-8 flex items-center justify-center bg-green-50 text-green-600 rounded-lg hover:bg-green-100 opacity-0 group-hover:opacity-100 transition-all"
                  title="Edit Semester"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7" />
                    <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />
                  </svg>
                </button>

                <div className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-700 rounded-2xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform">
                    <span className="text-3xl font-bold text-white">{sem.semester_number}</span>
                  </div>
                  <h3 className="text-lg font-bold text-gray-900 mb-1">Semester {sem.semester_number}</h3>
                  <p className="text-sm text-gray-600">Manage courses →</p>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Activity Logs removed - now in Department Overview page */}

        {/* Edit Semester Modal */}
        {showEditModal && editingSemester && (
          <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" onClick={() => setShowEditModal(false)}>
            <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full" onClick={(e) => e.stopPropagation()}>
              <div className="bg-gradient-to-r from-green-600 to-green-700 text-white px-8 py-5 flex items-center justify-between rounded-t-2xl">
                <div>
                  <h3 className="text-xl font-bold">Edit Semester</h3>
                  <p className="text-sm text-green-100">Update semester number</p>
                </div>
                <button 
                  onClick={() => setShowEditModal(false)}
                  className="text-white hover:bg-white/20 rounded-lg p-2 transition-all"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              
              <form onSubmit={handleUpdateSemester} className="p-8 space-y-5">
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-2">Semester Number</label>
                  <input
                    type="number"
                    value={editSemesterNumber}
                    onChange={(e) => setEditSemesterNumber(e.target.value)}
                    placeholder="e.g., 1"
                    required
                    min="1"
                    className="input-custom"
                  />
                </div>

                <div className="flex gap-3 justify-end pt-2">
                  <button
                    type="button"
                    onClick={() => setShowEditModal(false)}
                    className="btn-secondary-custom"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="bg-green-600 hover:bg-green-700 text-white font-medium px-5 py-2.5 rounded-lg transition-all duration-200 shadow-sm hover:shadow-md active:scale-95"
                  >
                    Update Semester
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </MainLayout>
  )
}

export default ManageCurriculumPage
