import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import MainLayout from '../components/MainLayout'

function SharingManagementPage() {
  const navigate = useNavigate()
  const [clusters, setClusters] = useState([])
  const [selectedCluster, setSelectedCluster] = useState(null)
  const [clusterDepartments, setClusterDepartments] = useState([])
  const [selectedRegulation, setSelectedRegulation] = useState(null)
  const [sharingInfo, setSharingInfo] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [clusterContent, setClusterContent] = useState(null)

  useEffect(() => {
    fetchClusters()
  }, [])

  const fetchClusters = async () => {
    try {
      setLoading(true)
      const response = await fetch('http://localhost:8080/api/clusters')
      if (!response.ok) {
        throw new Error('Failed to fetch clusters')
      }
      const data = await response.json()
      setClusters(data || [])
      setError('')
    } catch (err) {
      console.error('Error fetching clusters:', err)
      setError('Failed to load clusters')
    } finally {
      setLoading(false)
    }
  }

  const fetchClusterDepartments = async (clusterId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/cluster/${clusterId}/departments`)
      if (!response.ok) {
        throw new Error('Failed to fetch cluster departments')
      }
      const data = await response.json()
      setClusterDepartments(data || [])
    } catch (err) {
      console.error('Error fetching cluster departments:', err)
      setClusterDepartments([])
    }
  }

  const handleSelectCluster = (cluster) => {
    setSelectedCluster(cluster)
    setSelectedRegulation(null)
    setSharingInfo(null)
    fetchClusterDepartments(cluster.id)
    fetchClusterContent(cluster.id)
  }

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
      setError('Failed to load regulations')
    } finally {
      setLoading(false)
    }
  }

  const fetchSharingInfo = async (regulationId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/regulation/${regulationId}/sharing`)
      if (!response.ok) {
        throw new Error('Failed to fetch sharing info')
      }
      const data = await response.json()
      setSharingInfo(data)
    } catch (err) {
      console.error('Error fetching sharing info:', err)
      setError('Failed to load sharing information')
    }
  }

  const fetchClusterContent = async (clusterId) => {
    try {
      const response = await fetch(`http://localhost:8080/api/cluster/${clusterId}/shared-content`)
      if (!response.ok) {
        throw new Error('Failed to fetch cluster content')
      }
      const data = await response.json()
      setClusterContent(data)
    } catch (err) {
      console.error('Error fetching cluster content:', err)
    }
  }

  const handleSelectRegulation = (dept) => {
    setSelectedRegulation(dept)
    fetchSharingInfo(dept.regulation_id)
  }

  const handleToggleVisibility = async (itemType, itemId, currentVisibility) => {
    const newVisibility = currentVisibility === 'UNIQUE' ? 'CLUSTER' : 'UNIQUE'
    
    try {
      const response = await fetch('http://localhost:8080/api/sharing/visibility', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          item_type: itemType,
          item_id: itemId,
          visibility: newVisibility,
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to update visibility')
      }

      setSuccess(`Updated to ${newVisibility === 'CLUSTER' ? 'Shared' : 'Private'}`)
      setTimeout(() => setSuccess(''), 3000)
      
      // Refresh sharing info
      if (selectedRegulation) {
        fetchSharingInfo(selectedRegulation.id)
      }
    } catch (err) {
      console.error('Error updating visibility:', err)
      setError('Failed to update sharing settings')
    }
  }

  const renderItemList = (items, itemType, title) => {
    if (!items || items.length === 0) {
      return (
        <div className="text-center py-8 text-gray-500">
          <p className="text-sm">No {title.toLowerCase()} items</p>
        </div>
      )
    }

    return (
      <div className="space-y-2">
        {items.map((item, index) => (
          <div
            key={item.id}
            className="flex items-start justify-between p-3 bg-gray-50 rounded-lg border border-gray-200 hover:border-gray-300 transition-colors"
          >
            <div className="flex-1 flex items-start space-x-3">
              <span className="text-sm font-semibold text-gray-500 mt-0.5">{index + 1}.</span>
              <p className="text-sm text-gray-900 flex-1">{item.text}</p>
            </div>
            
            {sharingInfo?.in_cluster && (
              <button
                onClick={() => handleToggleVisibility(itemType, item.id, item.visibility)}
                className={`ml-3 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                  item.visibility === 'CLUSTER'
                    ? 'bg-green-100 text-green-700 hover:bg-green-200'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
                title={item.visibility === 'CLUSTER' ? 'Shared with cluster' : 'Private to this department'}
              >
                {item.visibility === 'CLUSTER' ? (
                  <div className="flex items-center space-x-1">
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    <span>Shared</span>
                  </div>
                ) : (
                  <div className="flex items-center space-x-1">
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    <span>Private</span>
                  </div>
                )}
              </button>
            )}
          </div>
        ))}
      </div>
    )
  }

  const renderClusterSharedContent = () => {
    if (!clusterContent || !clusterContent.departments) return null

    return (
      <div className="card-custom p-6 mt-6">
        <h2 className="text-lg font-bold text-gray-900 mb-5 flex items-center space-x-2">
          <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <span>Shared Content in Cluster</span>
        </h2>
        
        {clusterContent.departments.map(dept => {
          const hasSharedContent = 
            (dept.mission && dept.mission.length > 0) ||
            (dept.peos && dept.peos.length > 0) ||
            (dept.pos && dept.pos.length > 0) ||
            (dept.psos && dept.psos.length > 0)

          if (!hasSharedContent) return null

          return (
            <div key={dept.department_id} className="mb-6 p-4 bg-green-50 rounded-lg border border-green-200">
              <h3 className="font-semibold text-gray-900 mb-3">{dept.name}</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {dept.mission && dept.mission.length > 0 && (
                  <div>
                    <p className="text-xs font-semibold text-gray-600 mb-2">Mission</p>
                    <ul className="space-y-1">
                      {dept.mission.map(item => (
                        <li key={item.id} className="text-sm text-gray-700">• {item.text}</li>
                      ))}
                    </ul>
                  </div>
                )}
                {dept.peos && dept.peos.length > 0 && (
                  <div>
                    <p className="text-xs font-semibold text-gray-600 mb-2">PEOs</p>
                    <ul className="space-y-1">
                      {dept.peos.map(item => (
                        <li key={item.id} className="text-sm text-gray-700">• {item.text}</li>
                      ))}
                    </ul>
                  </div>
                )}
                {dept.pos && dept.pos.length > 0 && (
                  <div>
                    <p className="text-xs font-semibold text-gray-600 mb-2">POs</p>
                    <ul className="space-y-1">
                      {dept.pos.map(item => (
                        <li key={item.id} className="text-sm text-gray-700">• {item.text}</li>
                      ))}
                    </ul>
                  </div>
                )}
                {dept.psos && dept.psos.length > 0 && (
                  <div>
                    <p className="text-xs font-semibold text-gray-600 mb-2">PSOs</p>
                    <ul className="space-y-1">
                      {dept.psos.map(item => (
                        <li key={item.id} className="text-sm text-gray-700">• {item.text}</li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            </div>
          )
        })}
      </div>
    )
  }

  if (loading) {
    return (
      <MainLayout title="Sharing Management" subtitle="Loading...">
        <div className="flex justify-center items-center py-20">
          <div className="text-center">
            <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p className="text-gray-600">Loading...</p>
          </div>
        </div>
      </MainLayout>
    )
  }

  return (
    <MainLayout 
      title="Sharing Management" 
      subtitle="Manage content sharing between cluster departments"
      actions={
        <button
          onClick={() => navigate('/dashboard')}
          className="btn-secondary-custom flex items-center space-x-2"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          <span>Back to Dashboard</span>
        </button>
      }
    >
      <div className="space-y-6">
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

        {/* Two Column Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Department Selection */}
          <div className="card-custom p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-5">Select Department</h2>
            
            {regulations.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <p className="text-sm">No departments available</p>
              </div>
            ) : (
              <div className="space-y-2">
                {regulations.map(reg => (
                  <button
                    key={reg.id}
                    onClick={() => handleSelectRegulation(reg)}
                    className={`w-full text-left p-3 rounded-lg border-2 transition-all ${
                      selectedRegulation?.id === reg.id
                        ? 'border-blue-500 bg-blue-50'
                        : 'border-gray-200 hover:border-blue-300 hover:bg-gray-50'
                    }`}
                  >
                    <p className="font-semibold text-gray-900">{reg.name}</p>
                    <p className="text-xs text-gray-500 mt-1">{reg.academic_year}</p>
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Sharing Settings */}
          <div className="lg:col-span-2">
            {!selectedRegulation ? (
              <div className="card-custom p-12 text-center">
                <svg className="w-16 h-16 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p className="text-sm text-gray-500">Select a department to manage sharing</p>
              </div>
            ) : !sharingInfo ? (
              <div className="card-custom p-12 text-center">
                <svg className="animate-spin h-12 w-12 text-blue-600 mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p className="text-gray-600">Loading sharing information...</p>
              </div>
            ) : (
              <div className="space-y-6">
                {/* Cluster Status */}
                <div className={`card-custom p-4 ${sharingInfo.in_cluster ? 'bg-green-50 border-green-200' : 'bg-yellow-50 border-yellow-200'} border-2`}>
                  <div className="flex items-center space-x-3">
                    {sharingInfo.in_cluster ? (
                      <>
                        <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <div>
                          <p className="font-semibold text-green-900">In Cluster: {sharingInfo.cluster_name}</p>
                          <p className="text-xs text-green-700">You can share content with other departments in this cluster</p>
                        </div>
                      </>
                    ) : (
                      <>
                        <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                        </svg>
                        <div>
                          <p className="font-semibold text-yellow-900">Not in any cluster</p>
                          <p className="text-xs text-yellow-700">Add this department to a cluster to enable sharing</p>
                        </div>
                      </>
                    )}
                  </div>
                </div>

                {/* Content Sections */}
                <div className="card-custom p-6">
                  <h3 className="text-md font-bold text-gray-900 mb-4">Mission Statements</h3>
                  {renderItemList(sharingInfo.mission, 'mission', 'Mission')}
                </div>

                <div className="card-custom p-6">
                  <h3 className="text-md font-bold text-gray-900 mb-4">Program Educational Objectives (PEOs)</h3>
                  {renderItemList(sharingInfo.peos, 'peos', 'PEO')}
                </div>

                <div className="card-custom p-6">
                  <h3 className="text-md font-bold text-gray-900 mb-4">Program Outcomes (POs)</h3>
                  {renderItemList(sharingInfo.pos, 'pos', 'PO')}
                </div>

                <div className="card-custom p-6">
                  <h3 className="text-md font-bold text-gray-900 mb-4">Program Specific Outcomes (PSOs)</h3>
                  {renderItemList(sharingInfo.psos, 'psos', 'PSO')}
                </div>

                {/* Cluster Shared Content */}
                {sharingInfo.in_cluster && renderClusterSharedContent()}
              </div>
            )}
          </div>
        </div>
      </div>
    </MainLayout>
  )
}

export default SharingManagementPage
