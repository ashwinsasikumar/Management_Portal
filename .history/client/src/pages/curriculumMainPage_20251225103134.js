import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function CurriculumMainPage() {
  const navigate = useNavigate()
  const [regulations, setRegulations] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({ name: '', academic_year: '', max_credits: '' })
  const [showLogsModal, setShowLogsModal] = useState(false)
  const [selectedCurriculumId, setSelectedCurriculumId] = useState(null)
  const [logs, setLogs] = useState([])
  const [logsLoading, setLogsLoading] = useState(false)
  const [selectedLog, setSelectedLog] = useState(null)
  const [showDiffModal, setShowDiffModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [editingCurriculum, setEditingCurriculum] = useState(null)
  const [editFormData, setEditFormData] = useState({ name: '', max_credits: '' })

  // Fetch regulations from backend
  useEffect(() => {
    fetchRegulations()
  }, [])

  const fetchRegulations = async () => {
    try {
      setLoading(true)
      const response = await fetch('http://localhost:8080/api/curriculum')
      if (!response.ok) {
        throw new Error('Failed to fetch regulations')
      }
      const data = await response.json()
      setRegulations(data || [])
      setError('')
    } catch (err) {
      console.error('Error fetching regulations:', err)
      setError('Failed to load regulations. Make sure the backend is running.')
      setRegulations([])
    } finally {
      setLoading(false)
    }
  }

  const fetchLogs = async (curriculumId) => {
    try {
      setLogsLoading(true)
      const response = await fetch(`http://localhost:8080/api/curriculum/${curriculumId}/logs`)
      if (!response.ok) {
        throw new Error('Failed to fetch logs')
      }
      const data = await response.json()
      setLogs(data || [])
    } catch (err) {
      console.error('Error fetching logs:', err)
      setLogs([])
    } finally {
      setLogsLoading(false)
    }
  }

  const handleViewLogs = (e, curriculumId) => {
    e.stopPropagation()
    setSelectedCurriculumId(curriculumId)
    setShowLogsModal(true)
    fetchLogs(curriculumId)
  }

  const handleDownloadPDF = async (e, regulationId, regulationName) => {
    e.stopPropagation()
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/pdf`)
      if (!response.ok) {
        throw new Error('Failed to generate PDF')
      }
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${regulationName.replace(/\s+/g, '_')}.pdf`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      console.error('Error downloading PDF:', err)
      alert('Failed to generate PDF. Please try again.')
    }
  }

  const handleAddRegulation = async (e) => {
    e.preventDefault()
    
    if (!formData.name.trim() || !formData.academic_year.trim()) {
      setError('Please fill in all fields')
      return
    }

    try {
      const response = await fetch('http://localhost:8080/api/curriculum/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...formData,
          max_credits: parseInt(formData.max_credits) || 0
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to create regulation')
      }

      // Reset form and refresh list
      setFormData({ name: '', academic_year: '', max_credits: '' })
      setShowForm(false)
      setError('')
      fetchRegulations()
    } catch (err) {
      console.error('Error creating regulation:', err)
      setError('Failed to create regulation')
    }
  }

  const handleEditCurriculum = (e, curriculum) => {
    e.stopPropagation()
    setEditingCurriculum(curriculum)
    setEditFormData({ name: curriculum.name, max_credits: curriculum.max_credits })
    setShowEditModal(true)
  }

  const handleUpdateCurriculum = async (e) => {
    e.preventDefault()
    
    if (!editFormData.name.trim()) {
      setError('Please fill in all fields')
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/curriculum/${editingCurriculum.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: editFormData.name,
          max_credits: parseInt(editFormData.max_credits) || 0
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to update curriculum')
      }

      setShowEditModal(false)
      setEditingCurriculum(null)
      setError('')
      fetchRegulations()
    } catch (err) {
      console.error('Error updating curriculum:', err)
      setError('Failed to update curriculum')
    }
  }

  const handleDeleteRegulation = async (id) => {
    if (!window.confirm('Are you sure you want to delete this regulation?')) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/curriculum/delete?id=${id}`, {
        method: 'DELETE',
      })

      if (!response.ok) {
        throw new Error('Failed to delete regulation')
      }

      setError('')
      fetchRegulations()
    } catch (err) {
      console.error('Error deleting regulation:', err)
      setError('Failed to delete regulation')
    }
  }

  const handleInputChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }))
  }

  return (
    <MainLayout 
      title="Curriculum & Syllabi" 
      subtitle="Manage academic curriculum and syllabi structures"
      actions={
        <div className="flex items-center space-x-3">
          <button 
            onClick={() => navigate('/dashboard')}
            className="btn-secondary-custom flex items-center space-x-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span>Back to Dashboard</span>
          </button>
          <button 
            onClick={() => setShowForm(!showForm)}
            className="btn-primary-custom flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            <span>{showForm ? 'Cancel' : 'Add Regulation'}</span>
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

        {/* Form */}
        {showForm && (
          <div className="card-custom p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-6">Add New Regulation</h2>
            <form onSubmit={handleAddRegulation} className="space-y-5">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Regulation Name</label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  placeholder="e.g., R2024 CSE"
                  required
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Academic Year</label>
                <input
                  type="text"
                  name="academic_year"
                  value={formData.academic_year}
                  onChange={handleInputChange}
                  placeholder="e.g., 2024-2025"
                  required
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Maximum Credits</label>
                <input
                  type="number"
                  name="max_credits"
                  value={formData.max_credits}
                  onChange={handleInputChange}
                  placeholder="e.g., 160"
                  required
                  min="0"
                  className="input-custom"
                />
              </div>

              <div className="flex gap-3 pt-2">
                <button type="submit" className="btn-primary-custom">Create Regulation</button>
                <button type="button" onClick={() => setShowForm(false)} className="btn-secondary-custom">Cancel</button>
              </div>
            </form>
          </div>
        )}

        {/* Loading State */}
        {loading ? (
          <div className="flex justify-center items-center py-20">
            <div className="text-center">
              <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <p className="text-gray-600">Loading regulations...</p>
            </div>
          </div>
        ) : regulations.length === 0 ? (
          <div className="card-custom p-12 text-center">
            <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Regulations Yet</h3>
            <p className="text-gray-600 mb-6">Get started by creating your first regulation</p>
            <button onClick={() => setShowForm(true)} className="btn-primary-custom">
              Add Regulation
            </button>
          </div>
        ) : (
          /* Regulations Grid */
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
            {regulations.map(reg => (
              <div 
                key={reg.id} 
                className="group card-custom p-6 cursor-pointer hover:scale-105 transition-all duration-200"
                onClick={() => navigate(`/regulation/${reg.id}/overview`)}
              >
                {/* Icon */}
                <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-700 rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                  <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>

                {/* Content */}
                <h3 className="text-lg font-bold text-gray-900 mb-2 line-clamp-2">{reg.name}</h3>
                <div className="space-y-2 mb-4">
                  <div className="flex items-center text-sm text-gray-600">
                    <svg className="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    <span>{reg.academic_year || 'N/A'}</span>
                  </div>
                  <div className="flex items-center text-sm text-gray-600">
                    <svg className="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                    </svg>
                    <span>Max Credits: {reg.max_credits || 0}</span>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={(e) => handleEditCurriculum(e, reg)}
                    title="Edit"
                    className="flex-1 px-3 py-2 text-xs font-medium bg-green-50 text-green-700 rounded-lg hover:bg-green-100 transition-colors"
                  >
                    Edit
                  </button>
                  <button
                    onClick={(e) => handleViewLogs(e, reg.id)}
                    title="Logs"
                    className="flex-1 px-3 py-2 text-xs font-medium bg-purple-50 text-purple-700 rounded-lg hover:bg-purple-100 transition-colors"
                  >
                    Logs
                  </button>
                  <button
                    onClick={(e) => handleDownloadPDF(e, reg.id, reg.name)}
                    title="PDF"
                    className="flex-1 px-3 py-2 text-xs font-medium bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors"
                  >
                    PDF
                  </button>
                  <button
                    onClick={(e) => {
                      e.stopPropagation()
                      handleDeleteRegulation(reg.id)
                    }}
                    title="Delete"
                    className="px-3 py-2 text-xs font-medium bg-red-50 text-red-700 rounded-lg hover:bg-red-100 transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}

      {/* Activity Logs Modal */}
      {showLogsModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" onClick={() => setShowLogsModal(false)}>
          <div className="bg-white rounded-2xl shadow-2xl max-w-5xl w-full max-h-[90vh] overflow-hidden" onClick={(e) => e.stopPropagation()}>
            <div className="bg-gradient-to-r from-blue-600 to-blue-700 text-white px-8 py-5 flex items-start justify-between">
              <div>
                <h3 className="text-2xl font-bold mb-1">Activity Timeline</h3>
                <p className="text-sm text-blue-100">Regulation ID: {selectedCurriculumId}</p>
              </div>
              <button 
                onClick={() => setShowLogsModal(false)}
                className="text-white hover:bg-white/20 rounded-lg p-2 transition-all flex-shrink-0 -mt-1"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div className="p-8 overflow-y-auto max-h-[calc(90vh-140px)] bg-gray-50">
              {logsLoading ? (
                <div className="flex justify-center items-center py-20">
                  <div className="text-center">
                    <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <p className="text-gray-600">Loading activity logs...</p>
                  </div>
                </div>
              ) : logs.length === 0 ? (
                <div className="card-custom p-12 text-center">
                  <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">No Activity Yet</h3>
                  <p className="text-gray-600">Activity logs will appear here as changes are made</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {logs.map((log) => {
                    const diffCount = log.diff ? Object.keys(log.diff).length : 0
                    return (
                      <div 
                        key={log.id}
                        className={`card-custom p-5 ${diffCount > 0 ? 'cursor-pointer hover:shadow-lg' : ''}`}
                        onClick={() => {
                          if (diffCount > 0) {
                            setSelectedLog(log)
                            setShowDiffModal(true)
                          }
                        }}
                      >
                            <div className="flex items-start justify-between mb-3">
                          <div className="flex-1">
                            <div className="flex items-center space-x-3 mb-2">
                              <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                                <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                              </div>
                              <div className="flex-1">
                                <h4 className="text-base font-bold text-gray-900">{log.action}</h4>
                                <p className="text-sm text-gray-600">{log.description}</p>
                              </div>
                            </div>
                          </div>
                          {diffCount > 0 && (
                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">{diffCount} change{diffCount > 1 ? 's' : ''}</span>
                          )}
                        </div>
                        
                        <div className="flex items-center gap-6 text-xs text-gray-500 ml-13">
                          <span className="flex items-center gap-1.5">
                            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                            </svg>
                            {log.changed_by}
                          </span>
                          <span className="flex items-center gap-1.5">
                            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clipRule="evenodd" />
                            </svg>
                            {new Date(log.created_at).toLocaleString()}
                          </span>
                          {diffCount > 0 && (
                            <span className="text-blue-600 font-medium">Click to view details →</span>
                          )}
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>

            <div className="bg-white px-8 py-5 flex justify-end border-t border-gray-200">
              <button
                onClick={() => setShowLogsModal(false)}
                className="btn-primary-custom"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Diff Modal */}
      {showDiffModal && selectedLog && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[60] p-4" onClick={() => setShowDiffModal(false)}>
          <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden" onClick={(e) => e.stopPropagation()}>
            <div className="bg-gradient-to-r from-blue-600 to-blue-700 text-white px-8 py-5 flex items-center justify-between">
              <div>
                <h3 className="text-2xl font-bold mb-1">{selectedLog.action}</h3>
                <p className="text-sm text-blue-100">{new Date(selectedLog.created_at).toLocaleString()}</p>
              </div>
              <button 
                onClick={() => setShowDiffModal(false)}
                className="text-white hover:bg-white/20 rounded-xl p-2.5 transition-all"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div className="p-8 overflow-y-auto max-h-[calc(90vh-160px)] bg-gray-50">
              <div className="space-y-5">
                {selectedLog.diff && Object.entries(selectedLog.diff).map(([field, changes]) => {
                  // Ensure changes object exists and has the expected structure
                  if (!changes || typeof changes !== 'object') {
                    console.warn('Invalid changes object for field:', field, changes);
                    return null;
                  }
                  
                  return (
                  <div key={field} className="card-custom p-5">
                    <div className="font-semibold text-gray-800 mb-3 capitalize flex items-center gap-2">
                      <svg className="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                      </svg>
                      {field.replace(/_/g, ' ')}
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <div className="text-xs font-medium text-red-600 mb-1 flex items-center gap-1">
                          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                          </svg>
                          Old Value
                        </div>
                        <div className="bg-red-50 border border-red-200 rounded p-3 text-sm">
                          <pre className="whitespace-pre-wrap break-words font-mono text-xs">
                            {typeof changes.old === 'object' && changes.old !== null 
                              ? JSON.stringify(changes.old, null, 2) 
                              : String(changes.old !== undefined && changes.old !== null ? changes.old : '(empty)')}
                          </pre>
                        </div>
                      </div>
                      <div>
                        <div className="text-xs font-medium text-green-600 mb-1 flex items-center gap-1">
                          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                          </svg>
                          New Value
                        </div>
                        <div className="bg-green-50 border border-green-200 rounded p-3 text-sm">
                          <pre className="whitespace-pre-wrap break-words font-mono text-xs">
                            {typeof changes.new === 'object' && changes.new !== null 
                              ? JSON.stringify(changes.new, null, 2) 
                              : String(changes.new !== undefined && changes.new !== null ? changes.new : '(empty)')}
                          </pre>
                        </div>
                      </div>
                    </div>
                  </div>
                  );
                })}
              </div>
            </div>

            <div className="bg-white px-8 py-5 flex justify-end border-t border-gray-200">
              <button
                onClick={() => setShowDiffModal(false)}
                className="btn-primary-custom"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Curriculum Modal */}
      {showEditModal && editingCurriculum && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[60] p-4" onClick={() => setShowEditModal(false)}>
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full" onClick={(e) => e.stopPropagation()}>
            <div className="bg-gradient-to-r from-green-600 to-green-700 text-white px-8 py-5 flex items-center justify-between rounded-t-2xl">
              <div>
                <h3 className="text-2xl font-bold mb-1">Edit Regulation</h3>
                <p className="text-sm text-green-100">Update regulation details</p>
              </div>
              <button 
                onClick={() => setShowEditModal(false)}
                className="text-white hover:bg-white/20 rounded-xl p-2.5 transition-all"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form onSubmit={handleUpdateCurriculum} className="p-8 space-y-5">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Regulation Name</label>
                <input
                  type="text"
                  value={editFormData.name}
                  onChange={(e) => setEditFormData({ ...editFormData, name: e.target.value })}
                  placeholder="Enter regulation name"
                  required
                  className="input-custom"
                />
              </div>

              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">Maximum Credits</label>
                <input
                  type="number"
                  value={editFormData.max_credits}
                  onChange={(e) => setEditFormData({ ...editFormData, max_credits: e.target.value })}
                  placeholder="e.g., 160"
                  required
                  min="0"
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
                  Update Regulation
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

export default CurriculumMainPage
