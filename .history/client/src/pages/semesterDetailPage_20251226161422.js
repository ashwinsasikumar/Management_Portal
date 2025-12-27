import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function SemesterDetailPage() {
  const { id, semId } = useParams()
  const navigate = useNavigate()
  
  const [courses, setCourses] = useState([])
  const [semester, setSemester] = useState(null)
  const [curriculum, setCurriculum] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [showAddForm, setShowAddForm] = useState(false)
  const [newCourse, setNewCourse] = useState({
    course_code: '',
    course_name: '',
    course_type: '',
    category: '',
    credit: '',
    lecture_hours: 0,
    theory_hours: 0,
    practical_hours: 0,
    activity_hours: 0,
    tutorial_hours: 0,
    tws_sl_hours: 0,
    total_hours: 0,
    cia_marks: 40,
    see_marks: 60
  })
  const [showEditModal, setShowEditModal] = useState(false)
  const [editingCourse, setEditingCourse] = useState(null)
  const [editCourseData, setEditCourseData] = useState({
    course_code: '',
    course_name: '',
    course_type: '',
    category: '',
    credit: '',
    lecture_hours: 0,
    theory_hours: 0,
    practical_hours: 0,
    activity_hours: 0,
    tutorial_hours: 0,
    tws_sl_hours: 0,
    total_hours: 0,
    cia_marks: 40,
    see_marks: 60
  })

  useEffect(() => {
    fetchCurriculum()
    fetchSemester()
    fetchCourses()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, semId])

  const fetchCurriculum = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/curriculum`)
      if (!response.ok) {
        throw new Error('Failed to fetch curriculum info')
      }
      const data = await response.json()
      const currentCurr = data.find(c => c.id === parseInt(id))
      setCurriculum(currentCurr)
    } catch (err) {
      console.error('Error fetching curriculum:', err)
    }
  }

  const fetchSemester = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semesters`)
      if (!response.ok) {
        throw new Error('Failed to fetch semester info')
      }
      const data = await response.json()
      const currentSem = data.find(s => s.id === parseInt(semId))
      setSemester(currentSem)
    } catch (err) {
      console.error('Error fetching semester:', err)
    }
  }

  const fetchCourses = async () => {
    try {
      setLoading(true)
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semester/${semId}/courses`)
      if (!response.ok) {
        throw new Error('Failed to fetch courses')
      }
      const data = await response.json()
      setCourses(data || [])
      setError('')
    } catch (err) {
      console.error('Error fetching courses:', err)
      setError('Failed to load courses')
    } finally {
      setLoading(false)
    }
  }

  const handleAddCourse = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semester/${semId}/course`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...newCourse,
          credit: parseInt(newCourse.credit),
          theory_hours: parseInt(newCourse.theory_hours) || 0,
          activity_hours: parseInt(newCourse.activity_hours) || 0,
          lecture_hours: parseInt(newCourse.lecture_hours) || 0,
          tutorial_hours: parseInt(newCourse.tutorial_hours) || 0,
          practical_hours: parseInt(newCourse.practical_hours) || 0,
          cia_marks: parseInt(newCourse.cia_marks) || 40,
          see_marks: parseInt(newCourse.see_marks) || 60
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to add course')
      }

      setNewCourse({
        course_code: '',
        course_name: '',
        course_type: '',
        category: '',
        credit: '',
        lecture_hours: 0,
        theory_hours: 0,
        practical_hours: 0,
        activity_hours: 0,
        tutorial_hours: 0,
        tws_sl_hours: 0,
        total_hours: 0,
        cia_marks: 40,
        see_marks: 60
      })
        practical_hours: 0,
        total_hours: 0,
        cia_marks: 40,
        see_marks: 60
      })
      setShowAddForm(false)
      setSuccess('Course added successfully!')
      setTimeout(() => setSuccess(''), 3000)
      fetchCourses()
    } catch (err) {
      console.error('Error adding course:', err)
      setError(err.message || 'Failed to add course')
      setTimeout(() => setError(''), 5000)
    }
  }

  const handleEditCourse = (course) => {
    setEditingCourse(course)
    setEditCourseData({
      course_code: course.course_code,
      course_name: course.course_name,
      course_type: course.course_type,
      category: course.category,
      credit: course.credit,
      lecture_hours: course.lecture_hours,
      theory_hours: course.theory_hours,
      practical_hours: course.practical_hours,
      activity_hours: course.activity_hours,
      tutorial_hours: course.tutorial_hours,
      tws_sl_hours: course.tws_sl_hours || 0,
      total_hours: course.total_hours,
      cia_marks: course.cia_marks,
      see_marks: course.see_marks
    })
    setShowEditModal(true)
  }

  const handleUpdateCourse = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/course/${editingCourse.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...editCourseData,
          credit: parseInt(editCourseData.credit),
          theory_hours: parseInt(editCourseData.theory_hours) || 0,
          activity_hours: parseInt(editCourseData.activity_hours) || 0,
          lecture_hours: parseInt(editCourseData.lecture_hours) || 0,
          tutorial_hours: parseInt(editCourseData.tutorial_hours) || 0,
          practical_hours: parseInt(editCourseData.practical_hours) || 0,
          cia_marks: parseInt(editCourseData.cia_marks) || 40,
          see_marks: parseInt(editCourseData.see_marks) || 60
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to update course')
      }

      setSuccess('Course updated successfully!')
      setTimeout(() => setSuccess(''), 3000)
      setShowEditModal(false)
      setEditingCourse(null)
      fetchCourses()
    } catch (err) {
      console.error('Error updating course:', err)
      setError('Failed to update course')
    }
  }

  const handleRemoveCourse = async (courseId) => {
    if (!window.confirm('Are you sure you want to remove this course from the semester?')) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/semester/${semId}/course/${courseId}`, {
        method: 'DELETE',
      })

      if (!response.ok) {
        throw new Error('Failed to remove course')
      }

      setSuccess('Course removed successfully!')
      setTimeout(() => setSuccess(''), 3000)
      fetchCourses()
    } catch (err) {
      console.error('Error removing course:', err)
      setError('Failed to remove course')
    }
  }

  if (loading) {
    return (
      <MainLayout title="Semester Courses" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading courses...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout 
      title={`Semester ${semester?.semester_number || semId} - Courses`}
      subtitle={
        <div className="flex items-center space-x-4">
          <span>Regulation ID: {id}</span>
          {curriculum && (
            <span className="text-blue-600 font-semibold">
              Total Credits: {courses.reduce((sum, c) => sum + c.credit, 0)} / {curriculum.max_credits || 0}
            </span>
          )}
        </div>
      }
      actions={
        <div className="flex items-center space-x-3">
          <button
            onClick={() => navigate(`/regulation/${id}/curriculum`)}
            className="btn-secondary-custom flex items-center space-x-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span>Back</span>
          </button>
          <button
            onClick={() => setShowAddForm(!showAddForm)}
            className="btn-primary-custom flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            <span>{showAddForm ? 'Cancel' : 'Add Course'}</span>
          </button>
        </div>
      }
    >
      <div className="max-w-7xl mx-auto space-y-6">

        {/* Messages */}
        {error && (
          <div className="flex items-start space-x-3 p-4 bg-red-50 border border-red-200 rounded-lg">
            <svg className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <p className="text-sm font-medium text-red-600">{error}</p>
          </div>
        )}
        
        {success && (
          <div className="flex items-start space-x-3 p-4 bg-green-50 border border-green-200 rounded-lg">
            <svg className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            <p className="text-sm font-medium text-green-600">{success}</p>
          </div>
        )}

        {/* Add Course Form */}
        {showAddForm && (
          <div className="card-custom p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Add New Course</h2>
            <form onSubmit={handleAddCourse} className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
                  <option value="Experiment">Experiment</option>
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
                  min="1"
                  className="input-custom"
                />
              </div>

              {/* Common Fields */}
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Lecture (hrs/week) *</label>
                <input
                  type="number"
                  value={newCourse.lecture_hours}
                  onChange={(e) => setNewCourse({ ...newCourse, lecture_hours: e.target.value })}
                  placeholder="0"
                  required
                  min="0"
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Theory (hrs/week)</label>
                <input
                  type="number"
                  value={newCourse.theory_hours}
                  onChange={(e) => setNewCourse({ ...newCourse, theory_hours: e.target.value })}
                  placeholder="0"
                  min="0"
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Practical (hrs/week)</label>
                <input
                  type="number"
                  value={newCourse.practical_hours}
                  onChange={(e) => setNewCourse({ ...newCourse, practical_hours: e.target.value })}
                  placeholder="0"
                  min="0"
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Activity (hrs/week)</label>
                <input
                  type="number"
                  value={newCourse.activity_hours}
                  onChange={(e) => setNewCourse({ ...newCourse, activity_hours: e.target.value })}
                  placeholder="0"
                  min="0"
                  className="input-custom"
                />
              </div>

              {/* Course Type Specific Fields */}
              {newCourse.course_type === 'Theory' && (
                <>
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Tutorial Hours</label>
                    <input
                      type="number"
                      value={newCourse.tutorial_hours}
                      onChange={(e) => setNewCourse({ ...newCourse, tutorial_hours: e.target.value })}
                      placeholder="0"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Total Hours (Auto)</label>
                    <input
                      type="number"
                      value={(parseInt(newCourse.theory_hours) || 0) + (parseInt(newCourse.tutorial_hours) || 0) + (parseInt(newCourse.activity_hours) || 0)}
                      readOnly
                      className="input-custom bg-gray-100 cursor-not-allowed"
                    />
                  </div>
                </>
              )}

              {newCourse.course_type === 'Experiment' && (
                <>
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">TWS/SL Hours</label>
                    <input
                      type="number"
                      value={newCourse.tws_sl_hours}
                      onChange={(e) => setNewCourse({ ...newCourse, tws_sl_hours: e.target.value })}
                      placeholder="0"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Total Hours (Auto)</label>
                    <input
                      type="number"
                      value={(parseInt(newCourse.practical_hours) || 0) + (parseInt(newCourse.tws_sl_hours) || 0)}
                      readOnly
                      className="input-custom bg-gray-100 cursor-not-allowed"
                    />
                  </div>
                </>
              )}

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">CIA Marks *</label>
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
                <label className="block text-sm font-semibold text-gray-700 mb-2">SEE Marks *</label>
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
                <label className="block text-sm font-semibold text-gray-700 mb-2">Total Score (Auto)</label>
                <input
                  type="number"
                  value={(parseInt(newCourse.cia_marks) || 0) + (parseInt(newCourse.see_marks) || 0)}
                  readOnly
                  className="input-custom bg-gray-100 cursor-not-allowed"
                />
              </div>

              <div className="md:col-span-2">
                <button type="submit" className="w-full btn-primary-custom">Add Course</button>
              </div>
            </form>
          </div>
        )}

        {/* Courses Table */}
        {courses.length === 0 ? (
          <div className="card-custom p-12 text-center">
            <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Courses Yet</h3>
            <p className="text-gray-600 mb-6">Get started by adding your first course</p>
            <button onClick={() => setShowAddForm(true)} className="btn-primary-custom">Add Course</button>
          </div>
        ) : (
          <div className="card-custom overflow-hidden">
            <div className="overflow-x-auto">
              <table className="table-custom">
                <thead>
                  <tr>
                    <th className="text-left">Course Code</th>
                    <th className="text-left">Course Name</th>
                    <th className="text-left">Type</th>
                    <th className="text-left">Category</th>
                    <th className="text-left">Credits</th>
                    <th className="text-left">L-T-P</th>
                    <th className="text-left">Marks</th>
                    <th className="text-center">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {courses.map(course => (
                    <tr key={course.id}>
                      <td className="font-medium">{course.course_code}</td>
                      <td>{course.course_name}</td>
                      <td>
                        <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
                          {course.course_type}
                        </span>
                      </td>
                      <td>
                        <span className="px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs font-medium">
                          {course.category}
                        </span>
                      </td>
                      <td className="font-semibold">{course.credit}</td>
                      <td className="font-mono text-sm">{course.lecture_hours || 0}-{course.tutorial_hours || 0}-{course.practical_hours || 0}</td>
                      <td>
                        <div className="text-xs space-y-1">
                          <div>CIA: <span className="font-semibold">{course.cia_marks || 0}</span></div>
                          <div>SEE: <span className="font-semibold">{course.see_marks || 0}</span></div>
                          <div className="text-blue-600 font-bold">Total: {course.total_marks || 0}</div>
                        </div>
                      </td>
                      <td className="text-center">
                        <div className="flex gap-2 justify-center flex-wrap">
                          <button
                            onClick={() => handleEditCourse(course)}
                            className="px-3 py-1.5 bg-green-600 hover:bg-green-700 text-white text-xs rounded-lg transition-all"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => navigate(`/course/${course.id}/syllabus`)}
                            className="px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs rounded-lg transition-all"
                          >
                            Syllabus
                          </button>
                          <button
                            onClick={() => navigate(`/course/${course.id}/mapping`)}
                            className="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 text-white text-xs rounded-lg transition-all"
                          >
                            Mapping
                          </button>
                          <button
                            onClick={() => handleRemoveCourse(course.id)}
                            className="px-3 py-1.5 bg-red-600 hover:bg-red-700 text-white text-xs rounded-lg transition-all"
                          >
                            Remove
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {/* Edit Course Modal */}
        {showEditModal && editingCourse && (
          <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" onClick={() => setShowEditModal(false)}>
            <div className="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
              <div className="bg-gradient-to-r from-green-600 to-green-700 text-white px-8 py-5 flex items-center justify-between sticky top-0 rounded-t-2xl">
                <div>
                  <h3 className="text-xl font-bold">Edit Course</h3>
                  <p className="text-sm text-green-100">Update course details</p>
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
              
              <form onSubmit={handleUpdateCourse} className="p-8 space-y-5">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Course Code</label>
                    <input
                      type="text"
                      value={editCourseData.course_code}
                      onChange={(e) => setEditCourseData({ ...editCourseData, course_code: e.target.value })}
                      placeholder="e.g., CS101"
                      required
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Course Title</label>
                    <input
                      type="text"
                      value={editCourseData.course_name}
                      onChange={(e) => setEditCourseData({ ...editCourseData, course_name: e.target.value })}
                      placeholder="e.g., Programming Fundamentals"
                      required
                      className="input-custom"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Course Type</label>
                    <select
                      value={editCourseData.course_type}
                      onChange={(e) => setEditCourseData({ ...editCourseData, course_type: e.target.value })}
                      required
                      className="input-custom"
                    >
                      <option value="">Select Type</option>
                      <option value="Theory">Theory</option>
                      <option value="Lab">Lab</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Category</label>
                    <select
                      value={editCourseData.category}
                      onChange={(e) => setEditCourseData({ ...editCourseData, category: e.target.value })}
                      required
                      className="input-custom"
                    >
                      <option value="">Select Category</option>
                      <option value="BS">BS - Basic Sciences</option>
                      <option value="ES">ES - Engineering Sciences</option>
                      <option value="HSS">HSS - Humanities & Social Sciences</option>
                      <option value="PC">PC - Professional Core</option>
                      <option value="PE">PE - Professional Elective</option>
                      <option value="EEC">EEC - Emerging Engineering Courses</option>
                    </select>
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Credits</label>
                    <input
                      type="number"
                      value={editCourseData.credit}
                      onChange={(e) => setEditCourseData({ ...editCourseData, credit: e.target.value })}
                      placeholder="4"
                      required
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Theory Hrs</label>
                    <input
                      type="number"
                      value={editCourseData.theory_hours}
                      onChange={(e) => setEditCourseData({ ...editCourseData, theory_hours: e.target.value })}
                      placeholder="0"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Activity Hrs</label>
                    <input
                      type="number"
                      value={editCourseData.activity_hours}
                      onChange={(e) => setEditCourseData({ ...editCourseData, activity_hours: e.target.value })}
                      placeholder="0"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Lecture (L)</label>
                    <input
                      type="number"
                      value={editCourseData.lecture_hours}
                      onChange={(e) => setEditCourseData({ ...editCourseData, lecture_hours: e.target.value })}
                      placeholder="3"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Tutorial Hrs</label>
                    <input
                      type="number"
                      value={editCourseData.tutorial_hours}
                      onChange={(e) => setEditCourseData({ ...editCourseData, tutorial_hours: e.target.value })}
                      placeholder="0"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Practical (P)</label>
                    <input
                      type="number"
                      value={editCourseData.practical_hours}
                      onChange={(e) => setEditCourseData({ ...editCourseData, practical_hours: e.target.value })}
                      placeholder="2"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">Total Hrs (Auto-calculated)</label>
                    <input
                      type="number"
                      value={(parseInt(editCourseData.theory_hours) || 0) + (parseInt(editCourseData.activity_hours) || 0) + (parseInt(editCourseData.lecture_hours) || 0) + (parseInt(editCourseData.tutorial_hours) || 0) + (parseInt(editCourseData.practical_hours) || 0)}
                      readOnly
                      className="input-custom bg-gray-100 cursor-not-allowed"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">CIA Marks</label>
                    <input
                      type="number"
                      value={editCourseData.cia_marks}
                      onChange={(e) => setEditCourseData({ ...editCourseData, cia_marks: e.target.value })}
                      placeholder="40"
                      min="0"
                      className="input-custom"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">SEE Marks</label>
                    <input
                      type="number"
                      value={editCourseData.see_marks}
                      onChange={(e) => setEditCourseData({ ...editCourseData, see_marks: e.target.value })}
                      placeholder="60"
                      min="0"
                      className="input-custom"
                    />
                  </div>
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
                    Update Course
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

export default SemesterDetailPage
