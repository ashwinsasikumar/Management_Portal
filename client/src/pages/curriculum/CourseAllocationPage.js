import React, { useState, useEffect } from 'react'
import MainLayout from '../../components/MainLayout'
import { API_BASE_URL } from '../../config'
import './CourseAllocationPage.css'

function CourseAllocationPage() {
  const [curriculums, setCurriculums] = useState([])
  const [semesters, setSemesters] = useState([])
  const [teachers, setTeachers] = useState([])
  const [courses, setCourses] = useState([])
  const [summary, setSummary] = useState(null)
  
  const [filters, setFilters] = useState({
    curriculum_id: '',
    semester_id: '',
    academic_year: '2025-2026'
  })

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  
  // Modal State
  const [showAddModal, setShowAddModal] = useState(false)
  const [selectedCourse, setSelectedCourse] = useState(null)
  const [newAlloc, setNewAlloc] = useState({
    teacher_id: '',
    section: 'A',
    role: 'Primary',
    allocation_id: null
  })

  useEffect(() => {
    fetchCurriculums()
    fetchTeachers()
  }, [])

  useEffect(() => {
    if (filters.curriculum_id) {
      fetchSemesters(filters.curriculum_id)
    } else {
      setSemesters([])
      setFilters(prev => ({ ...prev, semester_id: '' }))
    }
  }, [filters.curriculum_id])

  useEffect(() => {
    if (filters.semester_id && filters.academic_year) {
      fetchAllocations()
      fetchSummary()
    } else {
      setCourses([])
      setSummary(null)
    }
  }, [filters.semester_id, filters.academic_year])

  const fetchCurriculums = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/curriculum`)
      const data = await res.json()
      setCurriculums(Array.isArray(data) ? data : [])
      // Automatically select first curriculum and fetch its semesters
      if (data && data.length > 0) {
        const firstId = data[0].id
        setFilters(prev => ({ ...prev, curriculum_id: firstId }))
        // Immediately fetch semesters
        fetchSemesters(firstId)
      }
    } catch (err) {
      console.error('Error fetching curriculums:', err)
    }
  }

  const fetchSemesters = async (curriculumId) => {
    try {
      const res = await fetch(`${API_BASE_URL}/curriculum/${curriculumId}/semesters`)
      const data = await res.json()
      setSemesters(Array.isArray(data) ? data : [])
      // Automatically select first semester
      if (data && data.length > 0) {
        setFilters(prev => ({ ...prev, semester_id: data[0].id }))
      }
    } catch (err) {
      console.error('Error fetching semesters:', err)
      setSemesters([])
    }
  }

  const fetchTeachers = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/teachers`)
      const data = await res.json()
      setTeachers(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Error fetching teachers:', err)
    }
  }

  const fetchAllocations = async () => {
    setLoading(true)
    setError('')
    try {
      const res = await fetch(`${API_BASE_URL}/allocations?semester_id=${filters.semester_id}&academic_year=${filters.academic_year}`)
      if (!res.ok) throw new Error('Failed to fetch allocations')
      const data = await res.json()
      setCourses(data || [])
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const fetchSummary = async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/allocations/summary?semester_id=${filters.semester_id}&academic_year=${filters.academic_year}`)
      if (res.ok) {
        const data = await res.json()
        setSummary(data)
      }
    } catch (err) {
      console.error('Error fetching summary:', err)
    }
  }

  const handleAddFaculty = (course) => {
    setSelectedCourse(course)
    setNewAlloc({ teacher_id: '', section: 'A', role: 'Primary', allocation_id: null })
    setShowAddModal(true)
  }

  const handleEditAllocation = (course, allocation) => {
    setSelectedCourse(course)
    setNewAlloc({
      teacher_id: allocation.teacher_id.toString(),
      section: allocation.section,
      role: allocation.role,
      allocation_id: allocation.id
    })
    setShowAddModal(true)
  }

  const saveAllocation = async (e) => {
    e.preventDefault()
    try {
      // Find the selected semester to get its semester_number
      const selectedSem = semesters.find(s => s.id === parseInt(filters.semester_id))
      const semesterNumber = selectedSem?.semester_number || 1

      const payload = {
        course_id: selectedCourse.id,
        teacher_id: parseInt(newAlloc.teacher_id),
        academic_year: filters.academic_year,
        semester: semesterNumber,
        section: newAlloc.section,
        role: newAlloc.role
      }

      let res
      if (newAlloc.allocation_id) {
        // Update existing
        res = await fetch(`${API_BASE_URL}/allocations/${newAlloc.allocation_id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      } else {
        // Create new
        res = await fetch(`${API_BASE_URL}/allocations`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      }

      if (!res.ok) throw new Error('Failed to save allocation')
      
      setSuccess(newAlloc.allocation_id ? 'Allocation updated successfully!' : 'Allocation saved successfully!')
      setShowAddModal(false)
      fetchAllocations()
      fetchSummary()
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      alert(err.message)
    }
  }

  const removeAllocation = async (allocId) => {
    if (!window.confirm('Remove this faculty assignment?')) return
    try {
      const res = await fetch(`${API_BASE_URL}/allocations/${allocId}`, {
        method: 'DELETE'
      })
      if (!res.ok) throw new Error('Failed to remove allocation')
      fetchAllocations()
      fetchSummary()
      setSuccess('Allocation removed successfully!')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      alert(err.message)
    }
  }

  return (
    <MainLayout title="Course Allocation" subtitle="Assign faculty to courses">
      <div className="space-y-6">
        {/* Summary Cards */}
        {summary && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
              <div className="text-3xl font-bold text-gray-900">{summary.total_courses}</div>
              <div className="text-sm text-gray-500 font-medium mt-1">Total Courses</div>
            </div>
            <div className="bg-white p-6 rounded-xl shadow-sm border-l-4 border-green-500">
              <div className="text-3xl font-bold text-green-600">{summary.assigned_courses}</div>
              <div className="text-sm text-gray-500 font-medium mt-1">Assigned</div>
            </div>
            <div className="bg-white p-6 rounded-xl shadow-sm border-l-4 border-orange-500">
              <div className="text-3xl font-bold text-orange-600">{summary.unassigned_courses}</div>
              <div className="text-sm text-gray-500 font-medium mt-1">Unassigned</div>
            </div>
            <div className="bg-white p-6 rounded-xl shadow-sm border-l-4 border-blue-500">
              <div className="text-3xl font-bold text-blue-600">{summary.active_teachers}/{summary.total_teachers}</div>
              <div className="text-sm text-gray-500 font-medium mt-1">Active Teachers</div>
            </div>
          </div>
        )}

        {/* Filters */}
        <div className="card-custom p-6 bg-white shadow-sm border border-gray-100">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Curriculum</label>
              <select
                value={filters.curriculum_id}
                onChange={(e) => setFilters({ ...filters, curriculum_id: e.target.value })}
                className="input-custom"
              >
                <option value="">Select Curriculum</option>
                {curriculums.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Semester</label>
              <select
                value={filters.semester_id}
                onChange={(e) => setFilters({ ...filters, semester_id: e.target.value })}
                className="input-custom"
                disabled={!filters.curriculum_id}
              >
                <option value="">Select Semester</option>
                {semesters.map(s => (
                  <option key={s.id} value={s.id}>
                    {s.card_type === 'semester' ? `Semester ${s.semester_number}` : s.card_type}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Academic Year</label>
              <select
                value={filters.academic_year}
                onChange={(e) => setFilters({ ...filters, academic_year: e.target.value })}
                className="input-custom"
              >
                <option value="2024-2025">2024-2025</option>
                <option value="2025-2026">2025-2026</option>
                <option value="2026-2027">2026-2027</option>
              </select>
            </div>
          </div>
        </div>

        {success && (
          <div className="p-4 bg-green-50 border border-green-200 text-green-700 rounded-lg">
            {success}
          </div>
        )}

        {/* Allocation Table */}
        <div className="card-custom overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Course</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Details</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Faculty Allocation</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr><td colSpan="3" className="px-6 py-10 text-center text-gray-500">Loading courses...</td></tr>
              ) : courses.length === 0 ? (
                <tr><td colSpan="3" className="px-6 py-10 text-center text-gray-500">Select filter to view courses</td></tr>
              ) : (
                courses.map(course => (
                  <tr key={course.id}>
                    <td className="px-6 py-4">
                      <div className="text-sm font-bold text-gray-900">{course.course_code}</div>
                      <div className="text-sm text-gray-600">{course.course_name}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {course.course_type} | {course.credit} Credits
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-2 items-center">
                        {course.allocations && course.allocations.map(alloc => (
                          <div key={alloc.id} className={`inline-flex items-center px-3 py-1.5 rounded-lg text-sm font-medium border ${alloc.role === 'Primary' ? 'bg-blue-100 text-blue-800 border-blue-200' : 'bg-amber-50 text-amber-800 border-amber-200'}`}>
                            <span className="mr-1.5 font-bold text-xs bg-white bg-opacity-50 px-1.5 py-0.5 rounded">{alloc.section}</span>
                            <span className="mr-2">{alloc.teacher_name}</span>
                            <span className="text-xs opacity-75">({alloc.role})</span>
                            <button
                              onClick={() => handleEditAllocation(course, alloc)}
                              className="ml-2 text-blue-500 hover:text-blue-700 transition-colors"
                              title="Edit"
                            >
                              ✏️
                            </button>
                            <button
                              onClick={() => removeAllocation(alloc.id)}
                              className="ml-1 text-blue-400 hover:text-red-500 transition-colors"
                              title="Remove"
                            >
                              ✕
                            </button>
                          </div>
                        ))}
                        {(!course.allocations || course.allocations.length === 0) && (
                          <span className="text-orange-600 font-medium text-sm">⚠️ No teacher assigned</span>
                        )}
                        <button
                          onClick={() => handleAddFaculty(course)}
                          className="px-3 py-1.5 text-xs font-medium text-blue-600 border border-dashed border-blue-300 rounded-md hover:bg-blue-50 transition-colors"
                        >
                          + Add Faculty
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Add/Edit Allocation Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">
              {newAlloc.allocation_id ? 'Edit Faculty Assignment' : 'Assign Faculty'}: {selectedCourse?.course_name}
            </h3>
            <form onSubmit={saveAllocation} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Select Teacher</label>
                <select
                  value={newAlloc.teacher_id}
                  onChange={(e) => setNewAlloc({ ...newAlloc, teacher_id: e.target.value })}
                  className="input-custom"
                  required
                >
                  <option value="">Select Teacher</option>
                  {teachers.map(t => (
                    <option key={t.id} value={t.id}>{t.name} ({t.desg})</option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Section</label>
                  <select
                    value={newAlloc.section}
                    onChange={(e) => setNewAlloc({ ...newAlloc, section: e.target.value })}
                    className="input-custom"
                  >
                    <option value="A">A</option>
                    <option value="B">B</option>
                    <option value="C">C</option>
                    <option value="D">D</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Role</label>
                  <select
                    value={newAlloc.role}
                    onChange={(e) => setNewAlloc({ ...newAlloc, role: e.target.value })}
                    className="input-custom"
                  >
                    <option value="Primary">Primary</option>
                    <option value="Assistant">Assistant</option>
                  </select>
                </div>
              </div>
              <div className="flex justify-end space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowAddModal(false)}
                  className="btn-secondary-custom"
                >
                  Cancel
                </button>
                <button type="submit" className="btn-primary-custom">
                  {newAlloc.allocation_id ? 'Update Assignment' : 'Confirm Assignment'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </MainLayout>
  )
}

export default CourseAllocationPage
