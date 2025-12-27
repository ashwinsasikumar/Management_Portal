import React, { useState, useEffect } from 'react'

function RegulationsPage() {
  const [regulations, setRegulations] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({ name: '', academic_year: '' })

  // Fetch regulations from backend
  useEffect(() => {
    fetchRegulations()
  }, [])

  const fetchRegulations = async () => {
    try {
      setLoading(true)
      const response = await fetch('http://localhost:8080/api/regulations')
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

  const handleDownloadPDF = async (e, regulationId, regulationName) => {
    e.stopPropagation()
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/pdf`)
      if (!response.ok) {
        const errorText = await response.text()
        
        // Check if the error is about Chrome not being installed
        if (errorText.includes('Chrome') || errorText.includes('Chromium')) {
          const useHTML = window.confirm(
            'Chrome is required for PDF generation but is not installed.\n\n' +
            'Would you like to view the HTML preview instead?\n\n' +
            '(You can print it to PDF using your browser\'s print function)'
          )
          
          if (useHTML) {
            // Open HTML preview in new tab
            window.open(`http://localhost:8080/api/regulation/${regulationId}/pdf?preview=html`, '_blank')
            return
          }
        }
        
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
      alert('Failed to generate PDF. Please install Chrome: brew install --cask google-chrome')
    }
  }

  const handleAddRegulation = async (e) => {
    e.preventDefault()
    
    if (!formData.name.trim() || !formData.academic_year.trim()) {
      setError('Please fill in all fields')
      return
    }

    try {
      const response = await fetch('http://localhost:8080/api/regulations/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      if (!response.ok) {
        throw new Error('Failed to create regulation')
      }

      // Reset form and refresh list
      setFormData({ name: '', academic_year: '' })
      setShowForm(false)
      setError('')
      fetchRegulations()
    } catch (err) {
      console.error('Error creating regulation:', err)
      setError('Failed to create regulation')
    }
  }

  const handleDeleteRegulation = async (id) => {
    if (!window.confirm('Are you sure you want to delete this regulation?')) {
      return
    }

    try {
      const response = await fetch(`http://localhost:8080/api/regulations/delete?id=${id}`, {
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
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600 p-10">
      {/* Header */}
      <div className="flex justify-between items-center mb-10 bg-white/95 backdrop-blur-md px-10 py-8 rounded-2xl shadow-xl">
        <h1 className="text-4xl font-bold text-gray-800 tracking-tight">Curriculum Page</h1>
        <button 
          onClick={() => setShowForm(!showForm)}
          className="px-7 py-3.5 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
        >
          {showForm ? 'Cancel' : '+ Add Regulation'}
        </button>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border-l-4 border-red-500 text-red-700 p-4 mb-6 rounded-lg">
          {error}
        </div>
      )}

      {/* Form */}
      {showForm && (
        <div className="bg-white/95 backdrop-blur-md p-10 rounded-2xl shadow-xl mb-10">
          <form onSubmit={handleAddRegulation}>
            <div className="mb-6">
              <label className="block text-gray-700 font-semibold mb-2 text-sm">Regulation Name</label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                placeholder="Enter regulation name"
                required
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
              />
            </div>

            <div className="mb-6">
              <label className="block text-gray-700 font-semibold mb-2 text-sm">Academic Year</label>
              <input
                type="text"
                name="academic_year"
                value={formData.academic_year}
                onChange={handleInputChange}
                placeholder="e.g., 2024-2025"
                required
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
              />
            </div>

            <button 
              type="submit" 
              className="px-8 py-3.5 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
            >
              Create Regulation
            </button>
          </form>
        </div>
      )}

      {/* Loading State */}
      {loading ? (
        <div className="text-center py-16 text-white text-xl font-medium">
          Loading regulations...
        </div>
      ) : regulations.length === 0 ? (
        <div className="text-center py-16 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl text-gray-600 text-lg">
          No regulations found. Add one to get started!
        </div>
      ) : (
        /* Regulations Grid */
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {regulations.map(reg => (
            <div 
              key={reg.id} 
              className="group relative bg-white/95 backdrop-blur-md rounded-xl shadow-md hover:shadow-xl p-5 transition-all duration-300 hover:-translate-y-1 border border-white/50 hover:border-indigo-300 overflow-hidden cursor-pointer"
              onClick={() => window.location.href = `/regulation/${reg.id}/overview`}
            >
              {/* Top colored border on hover */}
              <div className="absolute top-0 left-0 right-0 h-0.5 bg-gradient-to-r from-indigo-500 to-purple-600 transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300" />
              
              {/* Card Content */}
              <div className="flex items-center justify-center min-h-[50px]">
                <h3 className="text-base font-semibold text-gray-800 text-center leading-tight">
                  {reg.name}
                </h3>
              </div>

              {/* Action Buttons */}
              <div className="absolute top-2 right-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                {/* PDF Download Button */}
                <button
                  onClick={(e) => handleDownloadPDF(e, reg.id, reg.name)}
                  title="Download PDF"
                  className="w-7 h-7 flex items-center justify-center bg-blue-100 text-blue-500 rounded-lg hover:bg-blue-500 hover:text-white hover:scale-110 transition-all duration-300"
                >
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
                    <path d="M14 2v6h6M12 18v-6M9 15l3 3 3-3" />
                  </svg>
                </button>

                {/* Delete Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDeleteRegulation(reg.id)
                  }}
                  title="Delete regulation"
                  className="w-7 h-7 flex items-center justify-center bg-red-100 text-red-500 rounded-lg hover:bg-red-500 hover:text-white hover:scale-110 transition-all duration-300"
                >
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2M10 11v6M14 11v6" />
                  </svg>
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default RegulationsPage
