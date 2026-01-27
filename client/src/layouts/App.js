import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "../pages/curriculum/loginPage";
import Dashboard from "../pages/curriculum/dashboard";
import CurriculumMainPage from "../pages/curriculum/curriculumMainPage";
import DepartmentOverviewPage from "../pages/curriculum/departmentOverviewPage";
import ManageCurriculumPage from "../pages/curriculum/manageCurriculumPage";
import SemesterDetailPage from "../pages/curriculum/semesterDetailPage";
import HonourCardPage from "../pages/curriculum/honourCardPage";
import SyllabusPage from "../pages/curriculum/syllabusPage";
import MappingPage from "../pages/curriculum/mappingPage";
import PEOPOMappingPage from "../pages/curriculum/peoPOMappingPage";
import ClusterManagementPage from "../pages/curriculum/clusterManagementPage";
import SharingManagementPage from "../pages/curriculum/sharingManagementPage";
import RegulationPage from "../pages/regulation/regulationPage";
import RegulationEditorPage from "../pages/regulation/regulationEditorPage";
import UsersPage from "../pages/curriculum/usersPage";
import StudentDetailsPage from "../pages/student-teacher_entry/studentDetailsPage";
import TeacherStudentDashboard from "../pages/student-teacher_entry/TeacherStudentDashboard";
import TeacherDetailsPage from "../pages/student-teacher_entry/TeacherDetailsPage";
import TeacherStudentMappingPage from "../pages/student-teacher_entry/TeacherStudentMappingPage";
import CourseAllocationPage from "../pages/curriculum/CourseAllocationPage";
import PrivateRoute from "../components/PrivateRoute";

function App() {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
      <Route path="/users" element={<PrivateRoute><UsersPage /></PrivateRoute>} />
      <Route path="/Student_details" element={<PrivateRoute><StudentDetailsPage /></PrivateRoute>} />
      <Route path="/student-teacher-dashboard" element={<PrivateRoute><TeacherStudentDashboard /></PrivateRoute>} />
      <Route path="/teacher-details" element={<PrivateRoute><TeacherDetailsPage /></PrivateRoute>} />
      <Route path="/teacher-student-mapping" element={<PrivateRoute><TeacherStudentMappingPage /></PrivateRoute>} />
      <Route path="/course-allocation" element={<PrivateRoute><CourseAllocationPage /></PrivateRoute>} />
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
