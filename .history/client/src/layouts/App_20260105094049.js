import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "../pages/loginPage";
import Dashboard from "../pages/dashboard";
import CurriculumMainPage from "../pages/curriculumMainPage";
import DepartmentOverviewPage from "../pages/departmentOverviewPage";
import ManageCurriculumPage from "../pages/manageCurriculumPage";
import SemesterDetailPage from "../pages/semesterDetailPage";
import HonourCardPage from "../pages/honourCardPage";
import SyllabusPage from "../pages/syllabusPage";
import MappingPage from "../pages/mappingPage";
import PEOPOMappingPage from "../pages/peoPOMappingPage";
import ClusterManagementPage from "../pages/clusterManagementPage";
import SharingManagementPage from "../pages/sharingManagementPage";
import RegulationPage from "../pages/regulation/regulationPage";
import RegulationEditorPage from "../pages/regulation/regulationEditorPage";

function App() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/regulations" element={<RegulationPage />} />
      <Route path="/regulation/:id/editor" element={<RegulationEditorPage />} />
      <Route path="/curriculum" element={<CurriculumMainPage />} />
      <Route path="/clusters" element={<ClusterManagementPage />} />
      <Route path="/sharing" element={<SharingManagementPage />} />
      <Route
        path="/regulation/:id/overview"
        element={<DepartmentOverviewPage />}
      />
      <Route
        path="/regulation/:id/curriculum"
        element={<ManageCurriculumPage />}
      />
      <Route
        path="/regulation/:id/curriculum/semester/:semId"
        element={<SemesterDetailPage />}
      />
      <Route path="/course/:courseId/syllabus" element={<SyllabusPage />} />
      <Route path="/course/:courseId/mapping" element={<MappingPage />} />
      <Route
        path="/regulation/:id/peo-po-mapping"
        element={<PEOPOMappingPage />}
      />
    </Routes>
  );
}

export default App;
