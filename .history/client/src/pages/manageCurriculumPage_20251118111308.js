import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'

function ManageCurriculumPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [semesters, setSemesters] = useState([])
  const [logs, setLogs] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [newSemester, setNewSemester] = useState({ semester_number: '' })
  const [activeTab, setActiveTab] = useState('semesters')
  const [selectedLog, setSelectedLog] = useState(null)
  const [showDiffModal, setShowDiffModal] = useState(false)

  useEffect(() => {
    fetchSemesters()
    fetchLogs()
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

  const fetchLogs = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/curriculum/${id}/logs`)
      if (!response.ok) {
        throw new Error('Failed to fetch logs')
      }
      const data = await response.json()
      setLogs(data || [])
    } catch (err) {
      console.error('Error fetching logs:', err)
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
              onClick={() => navigate('/curriculum')}
              className="text-indigo-600 hover:text-indigo-800 font-medium mb-2 flex items-center gap-2"
            >
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M19 12H5M12 19l-7-7 7-7" />
              </svg>
              Back to Curriculum
            </button>
            <h1 className="text-3xl font-bold text-gray-800">Manage Curriculum</h1>
            <p className="text-gray-600 mt-1">Regulation ID: {id}</p>
          </div>
          <button
            onClick={() => setShowCreateForm(!showCreateForm)}
            className="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all"
          >
            {showCreateForm ? 'Cancel' : '+ Create Semester'}
          </button>
        </div>

        {/* Tabs */}
        <div className="mb-6 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl overflow-hidden">
          <div className="flex border-b border-gray-200">
            <button
              onClick={() => setActiveTab('semesters')}
              className={`flex-1 px-6 py-4 text-center font-semibold transition-all ${
                activeTab === 'semesters'
                  ? 'bg-gradient-to-r from-indigo-50 to-purple-50 text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <span className="text-xl">📚</span>
                <span>Semesters</span>
              </div>
            </button>
            <button
              onClick={() => setActiveTab('logs')}
              className={`flex-1 px-6 py-4 text-center font-semibold transition-all ${
                activeTab === 'logs'
                  ? 'bg-gradient-to-r from-indigo-50 to-purple-50 text-indigo-600 border-b-2 border-indigo-600'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <span className="text-xl">📝</span>
                <span>Activity Logs</span>
              </div>
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {/* Semesters Tab Content */}
        {activeTab === 'semesters' && (
          <>
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
          </>
        )}

        {/* Logs Tab Content */}
        {activeTab === 'logs' && (
          <div className="bg-white/95 backdrop-blur-md p-8 rounded-2xl shadow-xl">
            <h2 className="text-2xl font-bold text-gray-800 mb-6">Activity Timeline</h2>
            {logs.length === 0 ? (
              <div className="text-center py-16 text-gray-500">
                <div className="text-5xl mb-4">📝</div>
                <p className="text-lg">No activity logs yet</p>
              </div>
            ) : (
              <ul className="timeline timeline-vertical">
                {logs.map((log, index) => {
                  const diffCount = log.diff ? Object.keys(log.diff).length : 0
                  return (
                    <li key={log.id}>
                      {index > 0 && <hr className="bg-indigo-200" />}
                      <div 
                        className={`timeline-start timeline-box bg-gradient-to-r from-indigo-50 to-purple-50 border-2 border-indigo-200 ${diffCount > 0 ? 'cursor-pointer hover:shadow-lg transition-shadow' : ''}`}
                        onClick={() => {
                          if (diffCount > 0) {
                            setSelectedLog(log)
                            setShowDiffModal(true)
                          }
                        }}
                      >
                        <div className="font-bold text-indigo-700 mb-1 flex items-center justify-between">
                          <span>{log.action}</span>
                          {diffCount > 0 && (
                            <span className="badge badge-sm badge-info">{diffCount} field{diffCount > 1 ? 's' : ''} changed</span>
                          )}
                        </div>
                        <div className="text-sm text-gray-600 mb-2">{log.description}</div>
                        <div className="text-xs text-gray-500 flex items-center gap-4">
                          <span className="flex items-center gap-1">
                            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
                            </svg>
                            {log.changed_by}
                          </span>
                          <span className="flex items-center gap-1">
                            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clipRule="evenodd" />
                            </svg>
                            {new Date(log.created_at).toLocaleString()}
                          </span>
                        </div>
                        {diffCount > 0 && (
                          <div className="text-xs text-indigo-500 mt-2">
                            Click to view changes
                          </div>
                        )}
                      </div>
                      <div className="timeline-middle">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-5 h-5 text-indigo-600">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clipRule="evenodd" />
                        </svg>
                      </div>
                      <hr className="bg-indigo-200" />
                    </li>
                  )
                })}
              </ul>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default ManageCurriculumPage
