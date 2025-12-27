import React from 'react'
import {Routes,Route  } from 'react-router-dom'
import LoginPage from '../pages/loginPage'
import Dashboard from '../pages/dashboard'
import CurriculumMainPage from '../pages/curriculumMainPage'
import DepartmentOverviewPage from '../pages/departmentOverviewPage'
import ManageCurriculumPage from '../pages/manageCurriculumPage'
import SemesterDetailPage from '../pages/semesterDetailPage'
import SyllabusPage from '../pages/syllabusPage'
import MappingPage from '../pages/mappingPage'
import PEOPOMappingPage from '../pages/peoPOMappingPage'
import ClusterManagementPage from '../pages/clusterManagementPage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/curriculum" element={<CurriculumMainPage />} />
      <Route path="/clusters" element={<ClusterManagementPage />} />
      <Route path="/regulation/:id/overview" element={<DepartmentOverviewPage />} />
      <Route path="/regulation/:id/curriculum" element={<ManageCurriculumPage />} />
      <Route path="/regulation/:id/curriculum/semester/:semId" element={<SemesterDetailPage />} />
      <Route path="/course/:courseId/syllabus" element={<SyllabusPage />} />
      <Route path="/course/:courseId/mapping" element={<MappingPage />} />
      <Route path="/regulation/:id/peo-po-mapping" element={<PEOPOMappingPage />} />
    </Routes>
  )
}

export default App