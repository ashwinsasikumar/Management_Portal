import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function MappingPage() {
  const { courseId } = useParams()
  const navigate = useNavigate()
  
  const [cos, setCos] = useState([])
  const [coPoMatrix, setCoPoMatrix] = useState({})
  const [coPsoMatrix, setCoPsoMatrix] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const PO_COUNT = 12
  const PSO_COUNT = 3

  useEffect(() => {
    fetchMapping()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [courseId])

  const fetchMapping = async () => {
    try {
      setLoading(true)
      const response = await fetch(`http://localhost:8080/api/course/${courseId}/mapping`)
      if (!response.ok) {
        throw new Error('Failed to fetch mapping data')
      }
      const data = await response.json()
      setCos(data.cos || [])
      setCoPoMatrix(data.co_po_matrix || {})
      setCoPsoMatrix(data.co_pso_matrix || {})
      setError('')
    } catch (err) {
      console.error('Error fetching mapping:', err)
      setError('Failed to load mapping data')
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      // Convert matrix objects to arrays for backend
      const coPoArray = []
      const coPsoArray = []

      // Build CO-PO array
      cos.forEach((_, coIndex) => {
        for (let poIndex = 1; poIndex <= PO_COUNT; poIndex++) {
          const key = `${coIndex}-${poIndex}`
          const value = coPoMatrix[key] || 0
          if (value > 0) { // Only save non-zero values
            coPoArray.push({
              co_index: coIndex,
              po_index: poIndex,
              mapping_value: value
            })
          }
        }
      })

      // Build CO-PSO array
      cos.forEach((_, coIndex) => {
        for (let psoIndex = 1; psoIndex <= PSO_COUNT; psoIndex++) {
          const key = `${coIndex}-${psoIndex}`
          const value = coPsoMatrix[key] || 0
          if (value > 0) { // Only save non-zero values
            coPsoArray.push({
              co_index: coIndex,
              pso_index: psoIndex,
              mapping_value: value
            })
          }
        }
      })

      const response = await fetch(`http://localhost:8080/api/course/${courseId}/mapping`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          co_po_matrix: coPoArray,
          co_pso_matrix: coPsoArray
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to save mapping')
      }

      setSuccess('Mapping saved successfully!')
      setTimeout(() => setSuccess(''), 3000)
      setError('')
    } catch (err) {
      console.error('Error saving mapping:', err)
      setError('Failed to save mapping')
    }
  }

  const updateCoPoValue = (coIndex, poIndex, value) => {
    const key = `${coIndex}-${poIndex}`
    setCoPoMatrix({
      ...coPoMatrix,
      [key]: parseInt(value)
    })
  }

  const updateCoPsoValue = (coIndex, psoIndex, value) => {
    const key = `${coIndex}-${psoIndex}`
    setCoPsoMatrix({
      ...coPsoMatrix,
      [key]: parseInt(value)
    })
  }

  const getCoPoValue = (coIndex, poIndex) => {
    const key = `${coIndex}-${poIndex}`
    return coPoMatrix[key] || 0
  }

  const getCoPsoValue = (coIndex, psoIndex) => {
    const key = `${coIndex}-${psoIndex}`
    return coPsoMatrix[key] || 0
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <div className="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
          <p className="text-gray-600 font-medium">Loading mapping data...</p>
        </div>
      </div>
    )
  }

  if (cos.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="bg-white border-b border-gray-200 shadow-sm">
          <div className="max-w-7xl mx-auto px-6 py-4">
            <button
              onClick={() => navigate(-1)}
              className="flex items-center gap-2 px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all"
            >
              <span className="text-xl">←</span>
              <span className="font-medium">Back</span>
            </button>
          </div>
        </div>
        <div className="max-w-7xl mx-auto px-6 py-16">
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-8 text-center">
            <div className="text-5xl mb-4">⚠️</div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">No Course Outcomes Found</h2>
            <p className="text-gray-600 mb-6">
              Please add course outcomes in the syllabus page before creating CO-PO/PSO mappings.
            </p>
            <button
              onClick={() => navigate(`/course/${courseId}/syllabus`)}
              className="px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-all font-medium"
            >
              Go to Syllabus Page
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-50 shadow-sm">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <button
                onClick={() => navigate(-1)}
                className="flex items-center gap-2 px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all"
              >
                <span className="text-xl">←</span>
                <span className="font-medium">Back</span>
              </button>
              <div className="h-8 w-px bg-gray-300"></div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">CO-PO & CO-PSO Mapping</h1>
                <p className="text-sm text-gray-500">Course ID: {courseId}</p>
              </div>
            </div>
            <button
              onClick={handleSave}
              className="flex items-center gap-2 px-6 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold rounded-lg hover:shadow-lg hover:scale-105 transition-all"
            >
              <span className="text-lg">💾</span>
              Save Mapping
            </button>
          </div>
        </div>
      </div>

      {/* Messages */}
      {(error || success) && (
        <div className="max-w-7xl mx-auto px-6 pt-6">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-800 px-6 py-4 rounded-lg mb-4 flex items-center gap-3">
              <span className="text-xl">⚠️</span>
              <span className="font-medium">{error}</span>
            </div>
          )}
          {success && (
            <div className="bg-green-50 border border-green-200 text-green-800 px-6 py-4 rounded-lg mb-4 flex items-center gap-3">
              <span className="text-xl">✓</span>
              <span className="font-medium">{success}</span>
            </div>
          )}
        </div>
      )}

      <div className="max-w-7xl mx-auto px-6 py-8 space-y-8">
        {/* Legend */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
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

        {/* CO-PO Mapping Matrix */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div className="px-6 py-4 bg-gradient-to-r from-indigo-50 to-purple-50 border-b border-gray-200">
            <h2 className="text-xl font-bold text-gray-900 flex items-center gap-2">
              <span className="text-2xl">🎯</span>
              CO - PO Mapping Matrix
            </h2>
          </div>
          <div className="p-6 overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700 sticky left-0 z-10">
                    CO / PO
                  </th>
                  {Array.from({ length: PO_COUNT }, (_, i) => i + 1).map(poNum => (
                    <th key={poNum} className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700">
                      PO{poNum}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {cos.map((co, coIndex) => (
                  <tr key={coIndex} className="hover:bg-gray-50">
                    <td className="border border-gray-300 px-4 py-3 font-semibold text-sm text-gray-700 bg-gray-50 sticky left-0 z-10">
                      CO{coIndex + 1}
                    </td>
                    {Array.from({ length: PO_COUNT }, (_, i) => i + 1).map(poNum => (
                      <td key={poNum} className="border border-gray-300 px-2 py-2">
                        <select
                          value={getCoPoValue(coIndex, poNum)}
                          onChange={(e) => updateCoPoValue(coIndex, poNum, e.target.value)}
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

        {/* CO-PSO Mapping Matrix */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
          <div className="px-6 py-4 bg-gradient-to-r from-purple-50 to-pink-50 border-b border-gray-200">
            <h2 className="text-xl font-bold text-gray-900 flex items-center gap-2">
              <span className="text-2xl">🎓</span>
              CO - PSO Mapping Matrix
            </h2>
          </div>
          <div className="p-6 overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  <th className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700">
                    CO / PSO
                  </th>
                  {Array.from({ length: PSO_COUNT }, (_, i) => i + 1).map(psoNum => (
                    <th key={psoNum} className="border border-gray-300 bg-gray-100 px-4 py-3 text-sm font-semibold text-gray-700">
                      PSO{psoNum}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {cos.map((co, coIndex) => (
                  <tr key={coIndex} className="hover:bg-gray-50">
                    <td className="border border-gray-300 px-4 py-3 font-semibold text-sm text-gray-700 bg-gray-50">
                      CO{coIndex + 1}
                    </td>
                    {Array.from({ length: PSO_COUNT }, (_, i) => i + 1).map(psoNum => (
                      <td key={psoNum} className="border border-gray-300 px-2 py-2">
                        <select
                          value={getCoPsoValue(coIndex, psoNum)}
                          onChange={(e) => updateCoPsoValue(coIndex, psoNum, e.target.value)}
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

        {/* Course Outcomes Reference */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span className="text-xl">📋</span>
            Course Outcomes Reference
          </h3>
          <div className="space-y-2">
            {cos.map((co, index) => (
              <div key={index} className="flex gap-3 text-sm">
                <span className="font-semibold text-indigo-600 min-w-[60px]">CO{index + 1}:</span>
                <span className="text-gray-700">{co}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export default MappingPage
