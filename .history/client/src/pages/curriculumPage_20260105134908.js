import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'

function CurriculumPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  
  const [semesters, setSemesters] = useState([])
  const [honourCards, setHonourCards] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [cardType, setCardType] = useState('normal') // 'normal' or 'honour'
  const [showDropdown, setShowDropdown] = useState(false)
  const [newSemester, setNewSemester] = useState({ semester_number: '', max_credits: '' })
  const [newHonourCard, setNewHonourCard] = useState({ title: '', semester_number: '' })

  useEffect(() => {
    fetchSemesters()
    fetchHonourCards()
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

  const fetchHonourCards = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/honour-cards`)
      if (!response.ok) {
        throw new Error('Failed to fetch honour cards')
      }
      const data = await response.json()
      setHonourCards(data || [])
    } catch (err) {
      console.error('Error fetching honour cards:', err)
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
          semester_number: parseInt(newSemester.semester_number),
          max_credits: parseInt(newSemester.max_credits) || 0
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to create semester')
      }

      setNewSemester({ semester_number: '', max_credits: '' })
      setShowCreateForm(false)
      setShowDropdown(false)
      fetchSemesters()
    } catch (err) {
      console.error('Error creating semester:', err)
      setError('Failed to create semester')
    }
  }

  const handleCreateHonourCard = async (e) => {
    e.preventDefault()
    
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${id}/honour-card`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: newHonourCard.title,
          semester_number: parseInt(newHonourCard.semester_number)
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to create honour card')
      }

      setNewHonourCard({ title: '', semester_number: '' })
      setShowCreateForm(false)
      setShowDropdown(false)
      fetchHonourCards()
    } catch (err) {
      console.error('Error creating honour card:', err)
      setError('Failed to create honour card')
    }
  }

  const handleAddClick = (type) => {
    setCardType(type)
    setShowCreateForm(true)
    setShowDropdown(false)
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
              onClick={() => navigate('/regulations')}
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
          <div className="flex items-center gap-3 relative">
            {/* Info Icon with Tooltip */}
            <div className="group relative">
              <div className="w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center cursor-help">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-indigo-600">
                  <circle cx="12" cy="12" r="10"></circle>
                  <path d="M12 16v-4"></path>
                  <path d="M12 8h.01"></path>
                </svg>
              </div>
              <div className="invisible group-hover:visible absolute right-0 top-10 w-80 bg-gray-900 text-white text-xs rounded-lg p-4 shadow-xl z-50">
                <div className="mb-2">
                  <span className="font-semibold text-indigo-300">Normal Card:</span>
                  <span className="ml-1">semester, electives, verticals, open elective courses, one credit courses</span>
                </div>
                <div>
                  <span className="font-semibold text-purple-300">Honour Card:</span>
                  <span className="ml-1">honour card with verticals inside</span>
                </div>
              </div>
            </div>

            {/* Dropdown Button */}
            <div className="relative">
              <button
                onClick={() => setShowDropdown(!showDropdown)}
                className="px-6 py-3 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-xl shadow-lg hover:shadow-xl hover:-translate-y-0.5 transition-all flex items-center gap-2"
              >
                + Add
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M6 9l6 6 6-6" />
                </svg>
              </button>
              
              {showDropdown && (
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-xl border border-gray-200 z-50">
                  <button
                    onClick={() => handleAddClick('normal')}
                    className="w-full text-left px-4 py-3 hover:bg-indigo-50 text-gray-700 font-medium rounded-t-lg transition-colors"
                  >
                    Normal Card
                  </button>
                  <button
                    onClick={() => handleAddClick('honour')}
                    className="w-full text-left px-4 py-3 hover:bg-purple-50 text-gray-700 font-medium rounded-b-lg transition-colors"
                  >
                    Honour Card
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 text-red-700 p-4 rounded-lg">
            {error}
          </div>
        )}

        {/* Create Form */}
        {showCreateForm && (
          <div className="mb-6 bg-white/95 backdrop-blur-md p-6 rounded-2xl shadow-xl">
            {cardType === 'normal' ? (
              <form onSubmit={handleCreateSemester} className="flex gap-4 items-end">
                <div className="flex-1">
                  <label className="block text-gray-700 font-semibold mb-2 text-sm">Semester Number</label>
                  <input
                    type="number"
                    value={newSemester.semester_number}
                    onChange={(e) => setNewSemester({ ...newSemester, semester_number: e.target.value })}
                    placeholder="e.g., 1"
                    required
                    min="1"
                    className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
                <div className="flex-1">
                  <label className="block text-gray-700 font-semibold mb-2 text-sm">Maximum Credits</label>
                  <input
                    type="number"
                    value={newSemester.max_credits}
                    onChange={(e) => setNewSemester({ ...newSemester, max_credits: e.target.value })}
                    placeholder="e.g., 24"
                    required
                    min="0"
                    className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
                <button
                  type="submit"
                  className="px-6 py-2.5 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-semibold rounded-lg hover:shadow-lg transition-all"
                >
                  Create Normal Card
                </button>
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="px-6 py-2.5 bg-gray-200 text-gray-700 font-semibold rounded-lg hover:bg-gray-300 transition-all"
                >
                  Cancel
                </button>
              </form>
            ) : (
              <form onSubmit={handleCreateHonourCard} className="flex gap-4 items-end">
                <div className="flex-1">
                  <label className="block text-gray-700 font-semibold mb-2 text-sm">Honour Card Title</label>
                  <input
                    type="text"
                    value={newHonourCard.title}
                    onChange={(e) => setNewHonourCard({ ...newHonourCard, title: e.target.value })}
                    placeholder="e.g., Honours Program"
                    required
                    className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                  />
                </div>
                <div className="flex-1">
                  <label className="block text-gray-700 font-semibold mb-2 text-sm">Semester Number</label>
                  <input
                    type="number"
                    value={newHonourCard.semester_number}
                    onChange={(e) => setNewHonourCard({ ...newHonourCard, semester_number: e.target.value })}
                    placeholder="e.g., 1"
                    required
                    min="1"
                    className="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                  />
                </div>
                <button
                  type="submit"
                  className="px-6 py-2.5 bg-gradient-to-r from-purple-500 to-pink-600 text-white font-semibold rounded-lg hover:shadow-lg transition-all"
                >
                  Create Honour Card
                </button>
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="px-6 py-2.5 bg-gray-200 text-gray-700 font-semibold rounded-lg hover:bg-gray-300 transition-all"
                >
                  Cancel
                </button>
              </form>
            )}
          </div>
        )}

        {/* Cards Grid */}
        {semesters.length === 0 && honourCards.length === 0 ? (
          <div className="text-center py-16 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl text-gray-600 text-lg">
            No cards found. Create one to get started!
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {/* Normal Cards (Semesters, Verticals, Electives, etc.) */}
            {semesters.map(sem => {
              const cardType = sem.card_type || 'semester'
              const getCardLabel = () => {
                if (cardType === 'semester') return `Semester ${sem.semester_number}`
                if (cardType === 'vertical') {
                  if (sem.semester_number && sem.name) return `Vertical ${sem.semester_number}: ${sem.name}`
                  if (sem.semester_number) return `Vertical ${sem.semester_number}`
                  return sem.name || 'Vertical'
                }
                if (cardType === 'elective') return sem.name || 'Elective'
                if (cardType === 'open_elective') return sem.name || 'Open Elective'
                if (cardType === 'one_credit') return sem.name || 'One Credit'
                return sem.name || 'Card'
              }
              
              const getCardColor = () => {
                if (cardType === 'vertical') return 'from-green-500 to-emerald-600'
                if (cardType === 'elective') return 'from-orange-500 to-amber-600'
                if (cardType === 'open_elective') return 'from-cyan-500 to-blue-600'
                if (cardType === 'one_credit') return 'from-pink-500 to-rose-600'
                return 'from-indigo-500 to-purple-600' // semester
              }
              
              return (
                <div
                  key={`sem-${sem.id}`}
                  onClick={() => navigate(`/regulation/${id}/curriculum/semester/${sem.id}`)}
                  className="group bg-white/95 backdrop-blur-md rounded-xl shadow-md hover:shadow-xl p-6 transition-all duration-300 hover:-translate-y-1 cursor-pointer border border-white/50 hover:border-indigo-300"
                >
                  <div className="text-center">
                    {sem.semester_number ? (
                      <div className={`text-5xl font-bold bg-gradient-to-br ${getCardColor()} bg-clip-text text-transparent mb-2`}>
                        {sem.semester_number}
                      </div>
                    ) : (
                      <div className="text-4xl mb-2">
                        {cardType === 'vertical' ? '📊' : cardType === 'elective' ? '📚' : cardType === 'open_elective' ? '🔓' : '⭐'}
                      </div>
                    )}
                    <h3 className="text-lg font-semibold text-gray-800">{getCardLabel()}</h3>
                    <p className="text-sm text-gray-500 mt-2">Max Credits: {sem.max_credits || 0}</p>
                    <span className="inline-block mt-2 px-3 py-1 bg-indigo-100 text-indigo-700 text-xs font-semibold rounded-full">
                      {cardType === 'semester' ? 'Semester' : cardType === 'vertical' ? 'Vertical' : cardType === 'elective' ? 'Elective' : cardType === 'open_elective' ? 'Open Elective' : 'One Credit'}
                    </span>
                    <p className="text-xs text-gray-400 mt-2">Click to manage courses</p>
                  </div>
                </div>
              )
            })}

            {/* Honour Cards */}
            {honourCards.map(card => (
              <div
                key={`honour-${card.id}`}
                onClick={() => navigate(`/regulation/${id}/curriculum/honour/${card.id}`)}
                className="group bg-gradient-to-br from-purple-50 to-pink-50 rounded-xl shadow-md hover:shadow-xl p-6 transition-all duration-300 hover:-translate-y-1 cursor-pointer border border-purple-200 hover:border-purple-400"
              >
                <div className="text-center">
                  <div className="text-3xl font-bold text-purple-600 mb-2">🎖️</div>
                  <h3 className="text-lg font-semibold text-gray-800">{card.title}</h3>
                  <p className="text-sm text-gray-500 mt-2">Semester: {card.semester_number}</p>
                  <span className="inline-block mt-2 px-3 py-1 bg-purple-100 text-purple-700 text-xs font-semibold rounded-full">Honour</span>
                  <p className="text-xs text-gray-400 mt-2">Click to manage verticals</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default CurriculumPage
