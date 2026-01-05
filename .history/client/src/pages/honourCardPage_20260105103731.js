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
      <div className="min-h-screen bg-gradient-to-br from-purple-500 via-pink-500 to-purple-600 flex items-center justify-center">
        <div className="text-white text-xl font-medium">Loading honour card...</div>
      </div>
    )
  }

  if (!honourCard) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-purple-500 via-pink-500 to-purple-600 flex items-center justify-center">
        <div className="text-white text-xl font-medium">Honour card not found</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-500 via-pink-500 to-purple-600 p-6">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-6 bg-white/95 backdrop-blur-md px-8 py-6 rounded-2xl shadow-xl">
          <div>
            <button
              onClick={() => navigate(`/regulation/${regulationId}/curriculum`)}
              className="text-purple-600 hover:text-purple-800 font-medium mb-2 flex items-center gap-2"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M19 12H5M12 19l-7-7 7-7" />
              </svg>
              Back to Curriculum
            </button>
            <h1 className="text-3xl font-bold text-gray-800 flex items-center gap-3">
              <span className="text-4xl">🎖️</span>
              {honourCard.title}
            </h1>
            <p className="text-gray-600 mt-1">Semester: {honourCard.semester_number}</p>
          </div>
          <button
            onClick={() => setShowVerticalForm(!showVerticalForm)}
            className="px-6 py-3 bg-gradient-to-r from-purple-500 to-pink-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all"
          >
            {showVerticalForm ? 'Cancel' : '+ Add Vertical'}
          </button>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {/* Create Vertical Form */}
        {showVerticalForm && (
          <div className="mb-6 bg-white/95 backdrop-blur-md p-6 rounded-2xl shadow-xl">
            <form onSubmit={handleCreateVertical} className="flex gap-4 items-end">
              <div className="flex-1">
                <label className="block text-gray-700 font-semibold mb-2 text-sm">Vertical Name</label>
                <input
                  type="text"
                  value={newVerticalName}
                  onChange={(e) => setNewVerticalName(e.target.value)}
                  placeholder="e.g., Data Science Track"
                  required
                  className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              </div>
              <button
                type="submit"
                className="px-6 py-2.5 bg-gradient-to-r from-purple-500 to-pink-600 text-white font-semibold rounded-lg hover:shadow-lg transition-all"
              >
                Create Vertical
              </button>
            </form>
          </div>
        )}

        {/* Verticals List */}
        {!honourCard.verticals || honourCard.verticals.length === 0 ? (
          <div className="text-center py-16 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl text-gray-600 text-lg">
            No verticals found. Create one to get started!
          </div>
        ) : (
          <div className="space-y-4">
            {honourCard.verticals.map(vertical => (
              <div
                key={vertical.id}
                className="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl overflow-hidden"
              >
                {/* Vertical Header */}
                <div className="bg-gradient-to-r from-purple-500 to-pink-500 px-6 py-4 flex justify-between items-center">
                  <div className="flex items-center gap-3">
                    <button
                      onClick={() => setExpandedVertical(expandedVertical === vertical.id ? null : vertical.id)}
                      className="text-white hover:text-purple-100 transition-colors"
                    >
                      <svg 
                        width="24" 
                        height="24" 
                        viewBox="0 0 24 24" 
                        fill="none" 
                        stroke="currentColor" 
                        strokeWidth="2"
                        className={`transform transition-transform ${expandedVertical === vertical.id ? 'rotate-90' : ''}`}
                      >
                        <path d="M9 18l6-6-6-6" />
                      </svg>
                    </button>
                    <h3 className="text-xl font-bold text-white">{vertical.name}</h3>
                    <span className="px-3 py-1 bg-white/20 text-white text-sm font-semibold rounded-full">
                      {vertical.courses?.length || 0} courses
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => setShowAddCourse(showAddCourse === vertical.id ? null : vertical.id)}
                      className="px-4 py-2 bg-white text-purple-600 font-semibold rounded-lg hover:bg-purple-50 transition-all"
                    >
                      + Add Course
                    </button>
                    <button
                      onClick={() => handleDeleteVertical(vertical.id)}
                      className="px-4 py-2 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transition-all"
                    >
                      Delete
                    </button>
                  </div>
                </div>

                {/* Add Course Section */}
                {showAddCourse === vertical.id && (
                  <div className="bg-purple-50 px-6 py-4 border-b border-purple-200">
                    <div className="mb-2 text-sm text-gray-600">
                      To add a course, you need to create/select it from the course management section.
                      This feature can be extended to include a course search or creation form.
                    </div>
                    <div className="flex gap-2">
                      <input
                        type="number"
                        placeholder="Enter Course ID"
                        className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
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
                          const input = document.querySelector('input[placeholder="Enter Course ID"]')
                          const courseId = parseInt(input.value)
                          if (courseId) {
                            handleAddCourseToVertical(vertical.id, courseId)
                            input.value = ''
                          }
                        }}
                        className="px-4 py-2 bg-purple-600 text-white font-semibold rounded-lg hover:bg-purple-700 transition-all"
                      >
                        Add
                      </button>
                    </div>
                  </div>
                )}

                {/* Courses List */}
                {expandedVertical === vertical.id && (
                  <div className="p-6">
                    {!vertical.courses || vertical.courses.length === 0 ? (
                      <div className="text-center py-8 text-gray-500">
                        No courses in this vertical yet. Click "Add Course" to get started.
                      </div>
                    ) : (
                      <div className="space-y-3">
                        {vertical.courses.map(course => (
                          <div
                            key={course.id}
                            className="flex justify-between items-center p-4 bg-gradient-to-r from-purple-50 to-pink-50 rounded-lg border border-purple-200 hover:shadow-md transition-all"
                          >
                            <div className="flex-1">
                              <div className="flex items-center gap-3">
                                <span className="px-3 py-1 bg-purple-600 text-white text-xs font-bold rounded">
                                  {course.course_code}
                                </span>
                                <h4 className="font-semibold text-gray-800">{course.course_name}</h4>
                              </div>
                              <div className="mt-2 flex gap-4 text-sm text-gray-600">
                                <span>Credits: {course.credit}</span>
                                <span>Type: {course.course_type}</span>
                                <span>Category: {course.category}</span>
                              </div>
                            </div>
                            <button
                              onClick={() => handleRemoveCourseFromVertical(vertical.id, course.id)}
                              className="px-4 py-2 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transition-all"
                            >
                              Remove
                            </button>
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
