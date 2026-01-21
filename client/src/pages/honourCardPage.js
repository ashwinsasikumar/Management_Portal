import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'
import { API_BASE_URL } from '../config'

function HonourCardPage() {
  const { id: curriculumId, cardId } = useParams()
  const navigate = useNavigate()
  
  const [honourCard, setHonourCard] = useState(null)
  const [curriculumTemplate, setCurriculumTemplate] = useState('2026')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showVerticalForm, setShowVerticalForm] = useState(false)
  const [newVerticalName, setNewVerticalName] = useState('')
  const [expandedVertical, setExpandedVertical] = useState(null)
  const [showAddCourse, setShowAddCourse] = useState(null)
  const [newCourse, setNewCourse] = useState({
    course_code: '',
    course_name: '',
    course_type: '',
    category: '',
    credit: '',
    lecture_hrs: 0,
    tutorial_hrs: 0,
    practical_hrs: 0,
    activity_hrs: 0,
    tw_sl_hrs: 0,
    cia_marks: 40,
    see_marks: 60
  })

  useEffect(() => {
    fetchCurriculumTemplate()
    fetchHonourCard()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cardId])

  const fetchCurriculumTemplate = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/curriculum`)
      if (!response.ok) return
      const data = await response.json()
      const curr = data.find(c => c.id === parseInt(curriculumId))
      if (curr && curr.curriculum_template) {
        setCurriculumTemplate(curr.curriculum_template)
      }
    } catch (err) {
      console.error('Error fetching curriculum template:', err)
    }
  }

  const fetchHonourCard = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${API_BASE_URL}/curriculum/${curriculumId}/honour-cards`)
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

  const handleCreateVertical = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`${API_BASE_URL}/honour-card/${cardId}/vertical`, {
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
      const response = await fetch(`${API_BASE_URL}/honour-vertical/${verticalId}`, {
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

  const handleAddCourseToVertical = async (e, verticalId) => {
    e.preventDefault()

    try {
      const lectureHrs = parseInt(newCourse.lecture_hrs) || 0
      const tutorialHrs = parseInt(newCourse.tutorial_hrs) || 0
      const practicalHrs = parseInt(newCourse.practical_hrs) || 0
      const activityHrs = parseInt(newCourse.activity_hrs) || 0
      
      const courseData = {
        ...newCourse,
        credit: parseInt(newCourse.credit),
        lecture_hrs: lectureHrs,
        tutorial_hrs: tutorialHrs,
        practical_hrs: practicalHrs,
        activity_hrs: activityHrs,
        tw_sl_hrs: parseInt(newCourse.tw_sl_hrs) || 0,
        cia_marks: parseInt(newCourse.cia_marks) || 40,
        see_marks: parseInt(newCourse.see_marks) || 60
      }
      
      // Calculate total hours based on course type
      if (newCourse.course_type === 'Lab') {
        courseData.theory_total_hrs = 0
        courseData.tutorial_total_hrs = 0
        courseData.activity_total_hrs = 0
        courseData.practical_total_hrs = practicalHrs * 15
      } else if (newCourse.course_type === 'Theory') {
        courseData.theory_total_hrs = lectureHrs * 15
        courseData.tutorial_total_hrs = tutorialHrs * 15
        courseData.activity_total_hrs = activityHrs * 15
        courseData.practical_total_hrs = 0
      } else if (newCourse.course_type === 'Theory&Lab') {
        courseData.theory_total_hrs = lectureHrs * 15
        courseData.tutorial_total_hrs = tutorialHrs * 15
        courseData.practical_total_hrs = practicalHrs * 15
        courseData.activity_total_hrs = 0
      } else {
        // Default: calculate all
        courseData.theory_total_hrs = lectureHrs * 15
        courseData.tutorial_total_hrs = tutorialHrs * 15
        courseData.practical_total_hrs = practicalHrs * 15
        courseData.activity_total_hrs = activityHrs * 15
      }
      
      const response = await fetch(`${API_BASE_URL}/honour-vertical/${verticalId}/course`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(courseData),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => null)
        throw new Error(errorData?.error || 'Failed to add course to vertical')
      }

      // Reset form state and close panel
      setNewCourse({
        course_code: '',
        course_name: '',
        course_type: '',
        category: '',
        credit: '',
        lecture_hrs: 0,
        tutorial_hrs: 0,
        practical_hrs: 0,
        activity_hrs: 0,
        tw_sl_hrs: 0,
        cia_marks: 40,
        see_marks: 60
      })
      setShowAddCourse(null)
      fetchHonourCard()
    } catch (err) {
      console.error('Error adding course:', err)
      setError(err.message || 'Failed to add course')
    }
  }

  const handleRemoveCourseFromVertical = async (verticalId, courseId) => {
    if (!window.confirm('Are you sure you want to remove this course from the vertical?')) {
      return
    }

    try {
      const response = await fetch(`${API_BASE_URL}/honour-vertical/${verticalId}/course/${courseId}`, {
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
      <MainLayout title="Honour Programme" subtitle="Loading honour card...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading honour card...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  if (!honourCard) {
    return (
      <MainLayout title="Honour Programme" subtitle="Honour card not found">
        <div className="flex justify-center items-center py-20">
          <div className="card-custom max-w-md w-full text-center">
            <div className="text-5xl mb-3">❌</div>
            <h2 className="text-xl font-bold text-gray-900 mb-2">Honour Card Not Found</h2>
            <p className="text-gray-600 mb-6">The honour card you're looking for doesn't exist.</p>
            <button
              onClick={() => navigate(`/curriculum/${curriculumId}/curriculum`)}
              className="btn-primary-custom"
            >
              Back to Curriculum
            </button>
          </div>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout
      title={honourCard.title}
      subtitle={`Semester ${honourCard.semester_number} • Honour Programme`}
      actions={
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate(`/curriculum/${curriculumId}/curriculum`)}
            className="btn-secondary-custom flex items-center space-x-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span>Back</span>
          </button>
          <button
            onClick={() => setShowVerticalForm(!showVerticalForm)}
            className="btn-primary-custom flex items-center space-x-2"
          >
            {showVerticalForm ? (
              <>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
                <span>Cancel</span>
              </>
            ) : (
              <>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                <span>Add Vertical</span>
              </>
            )}
          </button>
        </div>
      }
    >
      <div className="space-y-6">

        {/* Error Message */}
        {error && (
          <div className="flex items-start space-x-3 p-4 bg-red-50 border border-red-200 rounded-lg">
            <svg className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <p className="text-sm font-medium text-red-600">{error}</p>
          </div>
        )}

        {/* Create Vertical Form */}
        {showVerticalForm && (
          <div className="card-custom">
            <h2 className="text-xl font-bold text-gray-900 mb-5">Add New Vertical</h2>
            <form onSubmit={handleCreateVertical} className="flex flex-col sm:flex-row gap-4 items-stretch sm:items-end">
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
              <div className="flex gap-3">
                <button type="submit" className="btn-primary-custom flex-1 sm:flex-none">
                  Create Vertical
                </button>
                <button
                  type="button"
                  onClick={() => setShowVerticalForm(false)}
                  className="btn-secondary-custom flex-1 sm:flex-none"
                >
                  Cancel
                </button>
              </div>
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
          <div className="space-y-6">
            {honourCard.verticals.map(vertical => (
              <div
                key={vertical.id}
                className="card-custom overflow-hidden hover:shadow-xl transition-all duration-200"
              >
                {/* Vertical Header */}
                <div className="bg-gradient-to-r from-purple-600 to-purple-700 px-5 py-5 sm:px-6">
                  <div className="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-4">
                    <div className="flex items-start gap-3 flex-1 min-w-0">
                      <button
                        onClick={() => setExpandedVertical(expandedVertical === vertical.id ? null : vertical.id)}
                        className="text-white hover:bg-white/20 p-2 rounded-lg transition-all flex-shrink-0 mt-1"
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
                      <div className="text-3xl sm:text-4xl flex-shrink-0">📊</div>
                      <div className="flex-1 min-w-0">
                        <h3 className="text-lg sm:text-xl font-bold text-white break-words mb-2">{vertical.name}</h3>
                        <div className="flex items-center gap-2">
                          <span className="inline-flex items-center px-2.5 py-1 bg-pink-500 text-white text-xs font-bold rounded-full">
                            {vertical.courses?.length || 0}
                          </span>
                          <span className="text-sm text-purple-100">
                            {vertical.courses?.length === 1 ? 'course' : 'courses'}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="flex gap-2 flex-wrap">
                      <button
                        onClick={() => {
                          if (showAddCourse === vertical.id) {
                            setShowAddCourse(null)
                          } else {
                            setShowAddCourse(vertical.id)
                          }
                        }}
                        className="px-4 py-2.5 bg-white text-purple-600 font-semibold rounded-lg hover:bg-purple-50 transition-all duration-200 flex items-center gap-2 text-sm"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                        </svg>
                        Add Course
                      </button>
                      <button
                        onClick={() => handleDeleteVertical(vertical.id)}
                        className="px-4 py-2.5 bg-red-500 text-white font-semibold rounded-lg hover:bg-red-600 transition-all duration-200 flex items-center gap-2 text-sm"
                        title="Delete vertical"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                        Delete
                      </button>
                    </div>
                  </div>
                </div>

                {/* Add Course Section */}
                {showAddCourse === vertical.id && (
                  <div className="bg-purple-50 px-5 py-5 sm:px-6 border-b border-purple-200">
                    <p className="text-sm text-gray-700 mb-4 font-medium flex items-center gap-2">
                      <span className="text-lg">💡</span>
                      Add a new course to this vertical (same fields as normal card)
                    </p>
                    <form onSubmit={(e) => handleAddCourseToVertical(e, vertical.id)} className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Course Code *</label>
                        <input
                          type="text"
                          value={newCourse.course_code}
                          onChange={(e) => setNewCourse({ ...newCourse, course_code: e.target.value })}
                          placeholder="e.g., CS101"
                          required
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Course Name *</label>
                        <input
                          type="text"
                          value={newCourse.course_name}
                          onChange={(e) => setNewCourse({ ...newCourse, course_name: e.target.value })}
                          placeholder="e.g., Introduction to Programming"
                          required
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Course Type *</label>
                        <select
                          value={newCourse.course_type}
                          onChange={(e) => setNewCourse({ ...newCourse, course_type: e.target.value })}
                          required
                          className="input-custom"
                        >
                          <option value="">Select Type</option>
                          <option value="Theory">Theory</option>
                          <option value="Lab">Lab</option>
                          <option value="Theory&Lab">Theory&Lab</option>
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Category *</label>
                        <select
                          value={newCourse.category}
                          onChange={(e) => setNewCourse({ ...newCourse, category: e.target.value })}
                          required
                          className="input-custom"
                        >
                          <option value="">Select Category</option>
                          <option value="BS - Basic Sciences">BS - Basic Sciences</option>
                          <option value="ES - Engineering Sciences">ES - Engineering Sciences</option>
                          <option value="HSS - Humanities and Social Sciences">HSS - Humanities and Social Sciences</option>
                          <option value="PC - Professional Core">PC - Professional Core</option>
                          <option value="PE - Professional Elective">PE - Professional Elective</option>
                          <option value="EEC - Employability Enhancement Course">EEC - Employability Enhancement Course</option>
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Credits *</label>
                        <input
                          type="number"
                          value={newCourse.credit}
                          onChange={(e) => setNewCourse({ ...newCourse, credit: e.target.value })}
                          placeholder="e.g., 3"
                          required
                          min="0"
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Lecture (hrs per week) *</label>
                        <input
                          type="number"
                          value={newCourse.lecture_hrs}
                          onChange={(e) => setNewCourse({ ...newCourse, lecture_hrs: e.target.value })}
                          placeholder="0"
                          required
                          min="0"
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">Tutorial (hrs per week)</label>
                        <input
                          type="number"
                          value={newCourse.tutorial_hrs}
                          onChange={(e) => setNewCourse({ ...newCourse, tutorial_hrs: e.target.value })}
                          placeholder="0"
                          min="0"
                          className="input-custom"
                        />
                      </div>

                      {newCourse.course_type !== 'Theory' && (
                        <div>
                          <label className="block text-sm font-semibold text-gray-700 mb-2">Practical (hrs per week)</label>
                          <input
                            type="number"
                            value={newCourse.practical_hrs}
                            onChange={(e) => setNewCourse({ ...newCourse, practical_hrs: e.target.value })}
                            placeholder="0"
                            min="0"
                            className="input-custom"
                          />
                        </div>
                      )}

                      {curriculumTemplate !== '2022' && (
                        <div>
                          <label className="block text-sm font-semibold text-gray-700 mb-2">Activity (hrs per week)</label>
                          <input
                            type="number"
                            value={newCourse.activity_hrs}
                            onChange={(e) => setNewCourse({ ...newCourse, activity_hrs: e.target.value })}
                            placeholder="0"
                            min="0"
                            className="input-custom"
                          />
                        </div>
                      )}

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">CIA *</label>
                        <input
                          type="number"
                          value={newCourse.cia_marks}
                          onChange={(e) => setNewCourse({ ...newCourse, cia_marks: e.target.value })}
                          placeholder="40"
                          required
                          min="0"
                          max="100"
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">SEE *</label>
                        <input
                          type="number"
                          value={newCourse.see_marks}
                          onChange={(e) => setNewCourse({ ...newCourse, see_marks: e.target.value })}
                          placeholder="60"
                          required
                          min="0"
                          max="100"
                          className="input-custom"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL SCORE (Auto)</label>
                        <input
                          type="number"
                          value={(parseInt(newCourse.cia_marks) || 0) + (parseInt(newCourse.see_marks) || 0)}
                          readOnly
                          className="input-custom bg-gray-100 cursor-not-allowed"
                        />
                      </div>

                      {/* Course Type Specific Fields - Total Hours for whole semester */}
                      {newCourse.course_type === 'Theory' && (
                        <>
                          <div className="md:col-span-2">
                            <h3 className="text-sm font-bold text-gray-900 mb-3 mt-4 pt-4 border-t border-gray-200">Total Hours (for whole semester)</h3>
                          </div>
                          
                          {curriculumTemplate === '2026' ? (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">THEORY HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.lecture_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TUTORIAL HOURS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.tutorial_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">ACTIVITY HOURS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.activity_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={((parseInt(newCourse.lecture_hrs) || 0) * 15) + ((parseInt(newCourse.tutorial_hrs) || 0) * 15) + ((parseInt(newCourse.activity_hrs) || 0) * 15)}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          ) : (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">THEORY HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.lecture_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TUTORIAL HOURS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.tutorial_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={((parseInt(newCourse.lecture_hrs) || 0) * 15) + ((parseInt(newCourse.tutorial_hrs) || 0) * 15)}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          )}
                        </>
                      )}

                      {newCourse.course_type === 'Theory&Lab' && (
                        <>
                          <div className="md:col-span-2">
                            <h3 className="text-sm font-bold text-gray-900 mb-3 mt-4 pt-4 border-t border-gray-200">Total Hours (for whole semester)</h3>
                          </div>
                          
                          {curriculumTemplate === '2026' ? (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">THEORY HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.lecture_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TUTORIAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.tutorial_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">PRACTICAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.practical_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={((parseInt(newCourse.lecture_hrs) || 0) * 15) + ((parseInt(newCourse.tutorial_hrs) || 0) * 15) + ((parseInt(newCourse.practical_hrs) || 0) * 15)}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          ) : (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">THEORY HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.lecture_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TUTORIAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.tutorial_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">PRACTICAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.practical_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={((parseInt(newCourse.lecture_hrs) || 0) * 15) + ((parseInt(newCourse.tutorial_hrs) || 0) * 15) + ((parseInt(newCourse.practical_hrs) || 0) * 15)}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          )}
                        </>
                      )}

                      {newCourse.course_type === 'Lab' && (
                        <>
                          <div className="md:col-span-2">
                            <h3 className="text-sm font-bold text-gray-900 mb-3 mt-4 pt-4 border-t border-gray-200">Total Hours (for whole semester)</h3>
                          </div>
                          
                          {curriculumTemplate === '2026' ? (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">PRACTICAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.practical_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TW/SL HRS</label>
                                <input
                                  type="number"
                                  value={newCourse.tw_sl_hrs}
                                  onChange={(e) => setNewCourse({ ...newCourse, tw_sl_hrs: e.target.value })}
                                  placeholder="0"
                                  min="0"
                                  className="input-custom"
                                />
                              </div>

                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">TOTAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={((parseInt(newCourse.practical_hrs) || 0) * 15) + (parseInt(newCourse.tw_sl_hrs) || 0)}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          ) : (
                            <>
                              <div>
                                <label className="block text-sm font-semibold text-gray-700 mb-2">PRACTICAL HRS (Auto)</label>
                                <input
                                  type="number"
                                  value={(parseInt(newCourse.practical_hrs) || 0) * 15}
                                  readOnly
                                  className="input-custom bg-gray-100 cursor-not-allowed"
                                />
                              </div>
                            </>
                          )}
                        </>
                      )}

                      <div className="md:col-span-2 flex justify-end gap-3 mt-2">
                        <button
                          type="button"
                          onClick={() => {
                            setShowAddCourse(null)
                          }}
                          className="btn-secondary-custom"
                        >
                          Cancel
                        </button>
                        <button type="submit" className="btn-primary-custom">
                          Add Course
                        </button>
                      </div>
                    </form>
                  </div>
                )}

                {/* Courses List */}
                {expandedVertical === vertical.id && (
                  <div className="p-5 sm:p-6 bg-gray-50">
                    {!vertical.courses || vertical.courses.length === 0 ? (
                      <div className="text-center py-16 bg-white rounded-xl border-2 border-dashed border-gray-300">
                        <svg className="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                        <p className="text-gray-500 font-semibold text-lg mb-1">No courses in this vertical yet</p>
                        <p className="text-sm text-gray-400">Click "Add Course" to get started</p>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {vertical.courses.map(course => (
                          <div
                            key={course.id}
                            className="bg-white rounded-xl border border-gray-200 p-5 sm:p-6 hover:shadow-lg hover:border-purple-200 transition-all duration-200"
                          >
                            <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start gap-4">
                              <div className="flex-1 min-w-0">
                                <div className="flex flex-wrap items-center gap-3 mb-4">
                                  <span className="inline-flex items-center px-3 py-1.5 bg-gradient-to-r from-purple-600 to-pink-600 text-white text-xs font-bold rounded-lg shadow-sm">
                                    {course.course_code}
                                  </span>
                                  <h4 className="font-bold text-gray-900 text-base sm:text-lg">{course.course_name}</h4>
                                </div>
                                <div className="flex flex-wrap gap-x-6 gap-y-2 text-sm">
                                  <span className="flex items-center gap-2 text-gray-600">
                                    <span className="w-2 h-2 bg-pink-500 rounded-full flex-shrink-0"></span>
                                    <span><strong className="font-semibold text-gray-700">Credits:</strong> {course.credit}</span>
                                  </span>
                                  <span className="flex items-center gap-2 text-gray-600">
                                    <span className="w-2 h-2 bg-pink-500 rounded-full flex-shrink-0"></span>
                                    <span><strong className="font-semibold text-gray-700">Type:</strong> {course.course_type}</span>
                                  </span>
                                  <span className="flex items-center gap-2 text-gray-600">
                                    <span className="w-2 h-2 bg-pink-500 rounded-full flex-shrink-0"></span>
                                    <span><strong className="font-semibold text-gray-700">Category:</strong> {course.category}</span>
                                  </span>
                                </div>
                              </div>
                              <button
                                onClick={() => handleRemoveCourseFromVertical(vertical.id, course.id)}
                                className="px-4 py-2.5 bg-red-50 text-red-600 font-semibold rounded-lg hover:bg-red-100 transition-all duration-200 flex items-center gap-2 flex-shrink-0 self-start lg:self-auto text-sm"
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
    </MainLayout>
  )
}

export default HonourCardPage
