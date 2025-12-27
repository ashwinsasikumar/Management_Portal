import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function DepartmentOverviewPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [vision, setVision] = useState('')
  const [mission, setMission] = useState([])
  const [peos, setPeos] = useState([])
  const [pos, setPos] = useState([])
  const [psos, setPsos] = useState([])
  
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  // Fetch existing data on mount
  useEffect(() => {
    fetchOverview()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const fetchOverview = async () => {
    try {
      setLoading(true)
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/overview`)
      if (!response.ok) {
        throw new Error('Failed to fetch overview data')
      }
      const data = await response.json()
      
      setVision(data.vision || '')
      // Store full item objects (with id, text, visibility, source_department_id)
      setMission((data.mission || []).map(item => typeof item === 'string' ? { text: item, visibility: 'UNIQUE' } : item))
      setPeos((data.peos || []).map(item => typeof item === 'string' ? { text: item, visibility: 'UNIQUE' } : item))
      setPos((data.pos || []).map(item => typeof item === 'string' ? { text: item, visibility: 'UNIQUE' } : item))
      setPsos((data.psos || []).map(item => typeof item === 'string' ? { text: item, visibility: 'UNIQUE' } : item))
      setError('')
    } catch (err) {
      console.error('Error fetching overview:', err)
      setError('Failed to load overview data')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async (e) => {
    e.preventDefault()
    
    try {
      setSaving(true)
      setError('')
      setSuccess('')
      
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/overview`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          vision,
          mission: mission.filter(item => item.trim() !== '').map(text => ({ text, visibility: 'UNIQUE' })),
          peos: peos.filter(item => item.trim() !== '').map(text => ({ text, visibility: 'UNIQUE' })),
          pos: pos.filter(item => item.trim() !== '').map(text => ({ text, visibility: 'UNIQUE' })),
          psos: psos.filter(item => item.trim() !== '').map(text => ({ text, visibility: 'UNIQUE' })),
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to save overview data')
      }

      setSuccess('Overview saved successfully!')
      setTimeout(() => setSuccess(''), 3000)
    } catch (err) {
      console.error('Error saving overview:', err)
      setError('Failed to save overview data')
    } finally {
      setSaving(false)
    }
  }

  const addItem = (setter, items) => {
    setter([...items, ''])
  }

  const removeItem = (setter, items, index) => {
    setter(items.filter((_, i) => i !== index))
  }

  const updateItem = (setter, items, index, value) => {
    const updated = [...items]
    updated[index] = value
    setter(updated)
  }

  const renderDynamicList = (title, items, setter, placeholder) => (
    <div className="card-custom p-6">
      <div className="flex justify-between items-center mb-5">
        <h3 className="text-lg font-bold text-gray-900">{title}</h3>
        <button
          type="button"
          onClick={() => addItem(setter, items)}
          className="btn-primary-custom text-sm flex items-center space-x-2"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Add Item</span>
        </button>
      </div>
      
      <div className="space-y-3">
        {items.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            <svg className="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p className="text-sm">No items yet. Click "Add Item" to start.</p>
          </div>
        ) : (
          items.map((item, index) => (
            <div key={index} className="flex gap-3">
              <div className="flex-shrink-0 w-8 h-10 flex items-center justify-center">
                <span className="text-sm font-semibold text-gray-500">{index + 1}.</span>
              </div>
              <input
                type="text"
                value={item}
                onChange={(e) => updateItem(setter, items, index, e.target.value)}
                placeholder={`${placeholder} ${index + 1}`}
                className="input-custom flex-1"
              />
              <button
                type="button"
                onClick={() => removeItem(setter, items, index)}
                className="flex-shrink-0 w-10 h-10 flex items-center justify-center bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition-all"
                title="Remove item"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M18 6L6 18M6 6l12 12" />
                </svg>
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  )

  if (loading) {
    return (
      <MainLayout title="Regulation Overview" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading overview data...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout 
      title="Regulation Overview" 
      subtitle={`Vision, Mission, PEOs & POs - Regulation ID: ${id}`}
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
            onClick={() => navigate(`/regulation/${id}/peo-po-mapping`)}
            className="btn-secondary-custom flex items-center space-x-2"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <span>PEO-PO Mapping</span>
          </button>
          <button
            onClick={() => navigate(`/regulation/${id}/curriculum`)}
            className="btn-primary-custom flex items-center space-x-2"
          >
            <span>Manage Curriculum</span>
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
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

        {/* Form */}
        <form onSubmit={handleSave} className="space-y-6">
          {/* Vision */}
          <div className="card-custom p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-5">Vision Statement</h3>
            <textarea
              value={vision}
              onChange={(e) => setVision(e.target.value)}
              placeholder="Enter the department vision statement..."
              rows="4"
              className="input-custom resize-none"
            />
          </div>

          {/* Mission */}
          {renderDynamicList('Mission Statements', mission, setMission, 'Mission statement')}

          {/* PEOs */}
          {renderDynamicList('Program Educational Objectives (PEOs)', peos, setPeos, 'PEO')}

          {/* POs */}
          {renderDynamicList('Program Outcomes (POs)', pos, setPos, 'PO')}

          {/* PSOs */}
          {renderDynamicList('Program Specific Outcomes (PSOs)', psos, setPsos, 'PSO')}

          {/* Save Button */}
          <div className="card-custom p-6">
            <div className="flex justify-end gap-3">
              <button
                type="button"
                onClick={() => navigate('/curriculum')}
                className="btn-secondary-custom"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={saving}
                className="btn-primary-custom disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                {saving ? (
                  <>
                    <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Saving...</span>
                  </>
                ) : (
                  <>
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span>Save Overview</span>
                  </>
                )}
              </button>
            </div>
          </div>
        </form>
      </div>
    </MainLayout>
  )
}

export default DepartmentOverviewPage
