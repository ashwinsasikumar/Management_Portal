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
      <MainLayout title="Course Mapping" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading mapping data...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  if (cos.length === 0) {
    return (
      <MainLayout title="CO-PO & CO-PSO Mapping" subtitle={`Course ID: ${courseId}`}>
        <div className="card-custom p-12 text-center">
          <svg className="w-20 h-20 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">No Course Outcomes Found</h3>
          <p className="text-gray-600 mb-6">Please add course outcomes in the syllabus page before creating mappings.</p>
          <button onClick={() => navigate(`/course/${courseId}/syllabus`)} className="btn-primary-custom">Go to Syllabus</button>
        </div>
      </MainLayout>
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
    <MainLayout 
      title="CO-PO & CO-PSO Mapping"
      subtitle={`Course ID: ${courseId}`}
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

        {/* CO-PO Mapping Matrix */}
        <div className="card-custom overflow-hidden">
          <div className="px-6 py-4 bg-gradient-to-r from-blue-50 to-blue-100 border-b border-gray-200">
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
        <div className="card-custom overflow-hidden">
          <div className="px-6 py-4 bg-gradient-to-r from-purple-50 to-purple-100 border-b border-gray-200">
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
        <div className="card-custom p-6">
          <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <span className="text-xl">📋</span>
            Course Outcomes Reference
          </h3>
          <div className="space-y-2">
            {cos.map((co, index) => (
              <div key={index} className="flex gap-3 text-sm">
                <span className="font-semibold text-blue-600 min-w-[60px]">CO{index + 1}:</span>
                <span className="text-gray-700">{co}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </MainLayout>
  )
}

export default MappingPage
