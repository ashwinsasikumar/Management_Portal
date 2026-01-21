import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'
import { API_BASE_URL } from '../config'

function PEOPOMappingPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [peos, setPeos] = useState([])
  const [pos, setPos] = useState([])
  const [matrix, setMatrix] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  useEffect(() => {
    fetchData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id])

  const fetchData = async () => {
    try {
      setLoading(true)
      
      // Fetch department overview (for PEOs and POs)
      const overviewResponse = await fetch(`${API_BASE_URL}/curriculum/${id}/overview`)
      if (!overviewResponse.ok) {
        throw new Error('Failed to fetch department overview')
      }
      const overviewData = await overviewResponse.json()
      setPeos(overviewData.peos || [])
      setPos(overviewData.pos || [])

      // Fetch existing PEO-PO mappings
      const mappingResponse = await fetch(`${API_BASE_URL}/curriculum/${id}/peo-po-mapping`)
      if (!mappingResponse.ok) {
        throw new Error('Failed to fetch PEO-PO mappings')
      }
      const mappingData = await mappingResponse.json()
      setMatrix(mappingData.matrix || {})
      
      setError('')
    } catch (err) {
      console.error('Error fetching data:', err)
      setError('Failed to load data')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      // Convert matrix object to array for backend with 1-based indexing
      const mappings = []
      peos.forEach((_, peoIndex) => {
        pos.forEach((_, poIndex) => {
          const key = `${peoIndex}-${poIndex}`
          const value = matrix[key] || 0
          if (value > 0) { // Only save non-zero values
            mappings.push({
              peo_index: peoIndex + 1,  // Convert to 1-based for database
              po_index: poIndex + 1,    // Convert to 1-based for database
              mapping_value: value
            })
          }
        })
      })

      const response = await fetch(`${API_BASE_URL}/curriculum/${id}/peo-po-mapping`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ mappings }),
      })

      if (!response.ok) {
        throw new Error('Failed to save mappings')
      }

      setSuccess('PEO-PO mappings saved successfully!')
      setTimeout(() => setSuccess(''), 3000)
      setError('')
    } catch (err) {
      console.error('Error saving mappings:', err)
      setError('Failed to save mappings')
    }
  }

  const updateValue = (peoIndex, poIndex, value) => {
    const key = `${peoIndex}-${poIndex}`
    setMatrix({
      ...matrix,
      [key]: parseInt(value)
    })
  }

  const getValue = (peoIndex, poIndex) => {
    const key = `${peoIndex}-${poIndex}`
    return matrix[key] || 0
  }

  if (loading) {
    return (
      <MainLayout title="PEO-PO Mapping" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading PEO-PO mapping...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  if (peos.length === 0 || pos.length === 0) {
    return (
      <MainLayout title="PEO-PO Mapping" subtitle={`Curriculum ID: ${id}`}>
        <div className="card-custom p-12 text-center">
          <svg className="w-20 h-20 text-yellow-400 mx-auto mb-4" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">PEOs or POs Not Found</h3>
          <p className="text-gray-600 mb-6">Please add Program Educational Objectives (PEOs) and Program Outcomes (POs) in the Department Overview page before creating PEO-PO mappings.</p>
          <button onClick={() => navigate(`/curriculum/${id}/overview`)} className="btn-primary-custom">Go to Department Overview</button>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout 
      title="PEO-PO Mapping"
      subtitle={`Curriculum ID: ${id}`}
      actions={
        <div className="flex items-center space-x-3">
          <button onClick={() => navigate(-1)} className="btn-secondary-custom flex items-center space-x-2">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span>Back</span>
          </button>
          <button onClick={handleSave} className="btn-primary-custom flex items-center space-x-2">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
            </svg>
            <span>Save Mapping</span>
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

        {/* Legend */}
        <div className="card-custom p-6">
          <h3 className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
            <span className="text-xl">📊</span>
            Mapping Scale
          </h3>
          <div className="grid grid-cols-4 gap-4">
            <div className="flex items-center gap-2">
              <span className="w-8 h-8 bg-gray-100 rounded flex items-center justify-center font-bold text-gray-600">0</span>
              <span className="text-sm text-gray-600">No Correlation</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-8 h-8 bg-green-100 rounded flex items-center justify-center font-bold text-green-700">1</span>
              <span className="text-sm text-gray-600">Low</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-8 h-8 bg-yellow-100 rounded flex items-center justify-center font-bold text-yellow-700">2</span>
              <span className="text-sm text-gray-600">Medium</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-8 h-8 bg-red-100 rounded flex items-center justify-center font-bold text-red-700">3</span>
              <span className="text-sm text-gray-600">High</span>
            </div>
          </div>
        </div>

        {/* PEO-PO Mapping Matrix */}
        <div className="card-custom overflow-hidden">
          <div className="px-6 py-4 bg-gradient-to-r from-blue-50 to-blue-100 border-b border-gray-200">
            <h2 className="text-xl font-bold text-gray-900 flex items-center gap-2">
              <span className="text-2xl">🎯</span>
              PEO - PO Mapping Matrix
            </h2>
          </div>
          <div className="p-6 overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700 sticky left-0 z-10">
                    PEO / PO
                  </th>
                  {pos.map((_, poIndex) => (
                    <th key={poIndex} className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700">
                      PO{poIndex + 1}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {peos.map((_, peoIndex) => (
                  <tr key={peoIndex} className="hover:bg-gray-50">
                    <td className="border border-gray-300 px-4 py-3 font-semibold text-sm text-gray-700 bg-gray-50 sticky left-0 z-10">
                      PEO{peoIndex + 1}
                    </td>
                    {pos.map((_, poIndex) => (
                      <td key={poIndex} className="border border-gray-300 px-2 py-2">
                        <select
                          value={getValue(peoIndex, poIndex)}
                          onChange={(e) => updateValue(peoIndex, poIndex, e.target.value)}
                          className="w-full px-2 py-1.5 border border-gray-300 rounded focus:border-indigo-500 focus:outline-none text-center font-semibold text-sm"
                        >
                          <option value="0">0</option>
                          <option value="1">1</option>
                          <option value="2">2</option>
                          <option value="3">3</option>
                        </select>
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* PEOs and POs Reference */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="card-custom p-6">
            <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <span className="text-xl">📋</span>
              Program Educational Objectives (PEOs)
            </h3>
            <div className="space-y-2">
              {peos.map((peo, index) => (
                <div key={index} className="flex gap-3 text-sm">
                  <span className="font-semibold text-blue-600 min-w-[70px]">PEO{index + 1}:</span>
                  <span className="text-gray-700">{peo.text || peo}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="card-custom p-6">
            <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <span className="text-xl">🎓</span>
              Program Outcomes (POs)
            </h3>
            <div className="space-y-2">
              {pos.map((po, index) => (
                <div key={index} className="flex gap-3 text-sm">
                  <span className="font-semibold text-purple-600 min-w-[60px]">PO{index + 1}:</span>
                  <span className="text-gray-700">{po.text || po}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  )
}

export default PEOPOMappingPage
