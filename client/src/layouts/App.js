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
import UsersPage from "../pages/usersPage";
import PrivateRoute from "../components/PrivateRoute";

function App() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
      <Route path="/users" element={<PrivateRoute><UsersPage /></PrivateRoute>} />
      <Route path="/regulations" element={<PrivateRoute><RegulationPage /></PrivateRoute>} />
      <Route path="/curriculum/:id/editor" element={<PrivateRoute><RegulationEditorPage /></PrivateRoute>} />
      <Route path="/curriculum" element={<PrivateRoute><CurriculumMainPage /></PrivateRoute>} />
      <Route path="/clusters" element={<PrivateRoute><ClusterManagementPage /></PrivateRoute>} />
      <Route path="/sharing" element={<PrivateRoute><SharingManagementPage /></PrivateRoute>} />
      <Route
        path="/curriculum/:id/overview"
        element={<PrivateRoute><DepartmentOverviewPage /></PrivateRoute>}
      />
      <Route
        path="/curriculum/:id/curriculum"
        element={<PrivateRoute><ManageCurriculumPage /></PrivateRoute>}
      />
      <Route
        path="/curriculum/:id/curriculum/semester/:semId"
        element={<PrivateRoute><SemesterDetailPage /></PrivateRoute>}
      />
      <Route
        path="/curriculum/:id/curriculum/honour/:cardId"
        element={<PrivateRoute><HonourCardPage /></PrivateRoute>}
      />
      <Route path="/course/:courseId/syllabus" element={<PrivateRoute><SyllabusPage /></PrivateRoute>} />
      <Route path="/course/:courseId/mapping" element={<PrivateRoute><MappingPage /></PrivateRoute>} />
      <Route
        path="/curriculum/:id/peo-po-mapping"
        element={<PrivateRoute><PEOPOMappingPage /></PrivateRoute>}
      />
    </Routes>
  );
}

export default App;
