import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'

function HonourCardPage() {
  const { id: regulationId, cardId } = useParams()
  const navigate = useNavigate()
  
  const [honourCard, setHonourCard] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showVerticalForm, setShowVerticalForm] = useState(false)
  const [newVerticalName, setNewVerticalName] = useState('')
  const [expandedVertical, setExpandedVertical] = useState(null)
  const [allCourses, setAllCourses] = useState([])
  const [showAddCourse, setShowAddCourse] = useState(null)

  useEffect(() => {
    fetchHonourCard()
    fetchAllCourses()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cardId])

  const fetchHonourCard = async () => {
    try {
      setLoading(true)
      const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/honour-cards`)
      if (!response.ok) {
        throw new Error('Failed to fetch honour cards')
      }
      const data = await response.json()
      const card = data.find(c => c.id === parseInt(cardId))
      setHonourCard(card || null)
      setError('')
    } catch (err) {
      console.error('Error fetching honour card:', err)
      setError('Failed to load honour card')
    } finally {
      setLoading(false)
    }
  }

  const fetchAllCourses = async () => {
    try {
      // Fetch all courses from the database
      // For now, we'll use a simple endpoint - you might need to adjust this
      const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/semesters`)
      if (response.ok) {
        // You may need to create a dedicated endpoint to get all courses
        // For simplicity, we'll implement course search in the UI
      }
    } catch (err) {
      console.error('Error fetching courses:', err)
    }
  }

  const handleCreateVertical = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/honour-card/${cardId}/vertical`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: newVerticalName
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to create vertical')
      }

      setNewVerticalName('')
      setShowVerticalForm(false)
      fetchHonourCard()
    } catch (err) {
      console.error('Error creating vertical:', err)
      setError('Failed to create vertical')
    }
  }

  const handleDeleteVertical = async (verticalId) => {
    if (!window.confirm('Are you sure you want to delete this vertical? All courses in it will be unlinked.')) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/honour-vertical/${verticalId}`, {
        method: 'DELETE',
      })

      if (!response.ok) {
        throw new Error('Failed to delete vertical')
      }

      fetchHonourCard()
    } catch (err) {
      console.error('Error deleting vertical:', err)
      setError('Failed to delete vertical')
    }
  }

  const handleAddCourseToVertical = async (verticalId, courseId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/honour-vertical/${verticalId}/course`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          course_id: courseId
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to add course to vertical')
      }

      setShowAddCourse(null)
      fetchHonourCard()
    } catch (err) {
      console.error('Error adding course:', err)
      setError('Failed to add course')
    }
  }

  const handleRemoveCourseFromVertical = async (verticalId, courseId) => {
    if (!window.confirm('Are you sure you want to remove this course from the vertical?')) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/honour-vertical/${verticalId}/course/${courseId}`, {
        method: 'DELETE',
      })

      if (!response.ok) {
        throw new Error('Failed to remove course')
      }

      fetchHonourCard()
    } catch (err) {
      console.error('Error removing course:', err)
      setError('Failed to remove course')
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center">
        <div className="text-gray-700 text-xl font-medium flex items-center gap-3">
          <svg className="animate-spin h-8 w-8 text-purple-600" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading honour card...
        </div>
      </div>
    )
  }

  if (!honourCard) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 flex items-center justify-center p-6">
        <div className="bg-white rounded-2xl shadow-xl p-12 text-center max-w-md">
          <div className="text-6xl mb-4">❌</div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Honour Card Not Found</h2>
          <p className="text-gray-600 mb-6">The honour card you're looking for doesn't exist.</p>
          <button
            onClick={() => navigate(`/regulation/${regulationId}/curriculum`)}
            className="btn-primary-custom"
          >
            Back to Curriculum
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="card-custom mb-6">
          <div className="flex justify-between items-center">
            <div>
              <button
                onClick={() => navigate(`/regulation/${regulationId}/curriculum`)}
                className="text-purple-600 hover:text-purple-800 font-medium mb-3 flex items-center gap-2 transition-colors"
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M19 12H5M12 19l-7-7 7-7" />
                </svg>
                Back to Curriculum
              </button>
              <div className="flex items-center gap-3">
                <span className="text-5xl">🎖️</span>
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">{honourCard.title}</h1>
                  <p className="text-gray-600 mt-1">Semester {honourCard.semester_number} • Honour Programme</p>
                </div>
              </div>
            </div>
            <button
              onClick={() => setShowVerticalForm(!showVerticalForm)}
              className="btn-primary-custom flex items-center gap-2"
            >
              {showVerticalForm ? (
                <>
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  Cancel
                </>
              ) : (
                <>
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                  Add Vertical
                </>
              )}
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="flex items-start space-x-3 p-4 bg-red-50 border border-red-200 rounded-lg mb-6">
            <svg className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <p className="text-sm font-medium text-red-600">{error}</p>
          </div>
        )}

        {/* Create Vertical Form */}
        {showVerticalForm && (
          <div className="card-custom mb-6">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Add New Vertical</h2>
            <form onSubmit={handleCreateVertical} className="flex gap-4 items-end">
              <div className="flex-1">
                <label className="block text-sm font-semibold text-gray-700 mb-2">Vertical Name</label>
                <input
                  type="text"
                  value={newVerticalName}
                  onChange={(e) => setNewVerticalName(e.target.value)}
                  placeholder="e.g., Data Science Track, AI & ML Specialization"
                  required
                  className="input-custom"
                />
              </div>
              <button type="submit" className="btn-primary-custom">
                Create Vertical
              </button>
              <button
                type="button"
                onClick={() => setShowVerticalForm(false)}
                className="btn-secondary-custom"
              >
                Cancel
              </button>
            </form>
          </div>
        )}

        {/* Verticals List */}
        {!honourCard.verticals || honourCard.verticals.length === 0 ? (
          <div className="card-custom p-12 text-center">
            <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Verticals Yet</h3>
            <p className="text-gray-600 mb-6">Create your first vertical to organize honour programme courses</p>
            <button onClick={() => setShowVerticalForm(true)} className="btn-primary-custom">
              + Add Vertical
            </button>
          </div>
        ) : (
          <div className="space-y-5">
            {honourCard.verticals.map(vertical => (
              <div
                key={vertical.id}
                className="card-custom overflow-hidden hover:shadow-xl transition-all duration-200"
              >
                {/* Vertical Header */}
                <div className="bg-gradient-to-r from-purple-600 to-pink-600 px-6 py-4 flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => setExpandedVertical(expandedVertical === vertical.id ? null : vertical.id)}
                      className="text-white hover:bg-white/20 p-2 rounded-lg transition-all"
                      title={expandedVertical === vertical.id ? "Collapse" : "Expand"}
                    >
                      <svg 
                        width="20" 
                        height="20" 
                        viewBox="0 0 24 24" 
                        fill="none" 
                        stroke="currentColor" 
                        strokeWidth="2"
                        className={`transform transition-transform duration-200 ${expandedVertical === vertical.id ? 'rotate-90' : ''}`}
                      >
                        <path d="M9 18l6-6-6-6" />
                      </svg>
                    </button>
                    <div className="text-3xl">📊</div>
                    <div>
                      <h3 className="text-xl font-bold text-white">{vertical.name}</h3>
                      <span className="text-sm text-purple-100">
                        {vertical.courses?.length || 0} {vertical.courses?.length === 1 ? 'course' : 'courses'}
                      </span>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => setShowAddCourse(showAddCourse === vertical.id ? null : vertical.id)}
                      className="px-4 py-2 bg-white text-purple-600 font-semibold rounded-lg hover:bg-purple-50 transition-all duration-200 flex items-center gap-2"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                      Add Course
                    </button>
                    <button
                      onClick={() => handleDeleteVertical(vertical.id)}
                      className="px-4 py-2 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transition-all duration-200 flex items-center gap-2"
                      title="Delete vertical"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      Delete
                    </button>
                  </div>
                </div>

                {/* Add Course Section */}
                {showAddCourse === vertical.id && (
                  <div className="bg-purple-50 px-6 py-4 border-b border-purple-200">
                    <p className="text-sm text-gray-700 mb-3 font-medium">
                      💡 Enter the Course ID to add it to this vertical
                    </p>
                    <div className="flex gap-3">
                      <input
                        type="number"
                        placeholder="Enter Course ID (e.g., 101)"
                        className="flex-1 input-custom"
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            const courseId = parseInt(e.target.value)
                            if (courseId) {
                              handleAddCourseToVertical(vertical.id, courseId)
                              e.target.value = ''
                            }
                          }
                        }}
                      />
                      <button
                        onClick={() => {
                          const input = document.querySelector('input[placeholder="Enter Course ID (e.g., 101)"]')
                          const courseId = parseInt(input.value)
                          if (courseId) {
                            handleAddCourseToVertical(vertical.id, courseId)
                            input.value = ''
                          }
                        }}
                        className="btn-primary-custom"
                      >
                        Add Course
                      </button>
                    </div>
                  </div>
                )}

                {/* Courses List */}
                {expandedVertical === vertical.id && (
                  <div className="p-6 bg-gray-50">
                    {!vertical.courses || vertical.courses.length === 0 ? (
                      <div className="text-center py-12 bg-white rounded-lg border-2 border-dashed border-gray-300">
                        <svg className="w-16 h-16 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                        <p className="text-gray-500 font-medium">No courses in this vertical yet</p>
                        <p className="text-sm text-gray-400 mt-1">Click "Add Course" to get started</p>
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {vertical.courses.map(course => (
                          <div
                            key={course.id}
                            className="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-all duration-200"
                          >
                            <div className="flex justify-between items-start">
                              <div className="flex-1">
                                <div className="flex items-center gap-3 mb-2">
                                  <span className="px-3 py-1 bg-gradient-to-r from-purple-600 to-pink-600 text-white text-xs font-bold rounded-lg">
                                    {course.course_code}
                                  </span>
                                  <h4 className="font-bold text-gray-900">{course.course_name}</h4>
                                </div>
                                <div className="flex flex-wrap gap-x-6 gap-y-1 text-sm text-gray-600">
                                  <span className="flex items-center gap-1">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    <strong>Credits:</strong> {course.credit}
                                  </span>
                                  <span className="flex items-center gap-1">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                                    </svg>
                                    <strong>Type:</strong> {course.course_type}
                                  </span>
                                  <span className="flex items-center gap-1">
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                                    </svg>
                                    <strong>Category:</strong> {course.category}
                                  </span>
                                </div>
                              </div>
                              <button
                                onClick={() => handleRemoveCourseFromVertical(vertical.id, course.id)}
                                className="ml-4 px-4 py-2 bg-red-50 text-red-600 font-semibold rounded-lg hover:bg-red-100 transition-all duration-200 flex items-center gap-2"
                                title="Remove from vertical"
                              >
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                </svg>
                                Remove
                              </button>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default HonourCardPage
