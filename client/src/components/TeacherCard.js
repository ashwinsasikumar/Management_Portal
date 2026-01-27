

import React, { useState, useEffect } from "react";
import { API_BASE_URL } from "../config";

function TeacherCard({ teacher, onEdit, onDelete }) {
  const [courses, setCourses] = useState([]);
  const [showCourses, setShowCourses] = useState(false);
  const [loading, setLoading] = useState(false);

  // Get base URL without /api for static files
  const baseUrl = API_BASE_URL.replace("/api", "");
  const imageUrl = teacher.profile_img
    ? `${baseUrl}${teacher.profile_img}`
    : null;

  const fetchTeacherCourses = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/teachers/${teacher.id}/courses?academic_year=2025-2026`);
      if (response.ok) {
        const data = await response.json();
        setCourses(data || []);
      }
    } catch (error) {
      console.error("Error fetching teacher courses:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleViewCourses = () => {
    if (!showCourses && courses.length === 0) {
      fetchTeacherCourses();
    }
    setShowCourses(!showCourses);
  };

  return (
    <div className="p-6 border rounded-lg bg-white shadow-sm hover:shadow-md transition-shadow">
      {/* Avatar and Name */}
      <div className="flex items-center space-x-3 mb-4">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={teacher.name}
            className="w-16 h-16 rounded-full object-cover shadow-md flex-shrink-0 border-2 border-white"
            onError={(e) => {
              e.target.style.display = "none";
              e.target.nextElementSibling.style.display = "flex";
            }}
          />
        ) : null}
        <div
          className="w-16 h-16 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-semibold text-2xl shadow-md flex-shrink-0"
          style={{ display: imageUrl ? "none" : "flex" }}
        >
          {teacher.name ? teacher.name.charAt(0).toUpperCase() : "T"}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-xl font-bold text-gray-900 truncate">
            {teacher.name || "—"}
          </h3>
          <p className="text-sm text-gray-600 mt-1">
            {teacher.designation || teacher.desg || "—"}
          </p>
        </div>
      </div>

      {/* Detailed Information */}
      <div className="space-y-3 mb-4">
        <div className="flex items-start">
          <svg
            className="w-5 h-5 text-gray-400 mr-2 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
            />
          </svg>
          <div className="flex-1 min-w-0">
            <p className="text-xs text-gray-500 uppercase font-medium">Email</p>
            <p className="text-sm text-gray-900 truncate" title={teacher.email}>
              {teacher.email || "—"}
            </p>
          </div>
        </div>

        <div className="flex items-start">
          <svg
            className="w-5 h-5 text-gray-400 mr-2 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"
            />
          </svg>
          <div className="flex-1">
            <p className="text-xs text-gray-500 uppercase font-medium">Phone</p>
            <p className="text-sm text-gray-900">{teacher.phone || "—"}</p>
          </div>
        </div>

        <div className="flex items-start">
          <svg
            className="w-5 h-5 text-gray-400 mr-2 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
            />
          </svg>
          <div className="flex-1">
            <p className="text-xs text-gray-500 uppercase font-medium">
              Department
            </p>
            <p className="text-sm text-gray-900">{teacher.department || "—"}</p>
          </div>
        </div>
      </div>

      {/* Status Badge */}
      <div className="mb-4">
        {teacher.status === 1 && (
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
            <svg
              className="w-3 h-3 mr-1"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
            Active
          </span>
        )}
      </div>

      {/* Action Buttons */}
      <div className="flex gap-2 pt-4 border-t border-gray-100">
        <button
          type="button"
          onClick={handleViewCourses}
          className="flex-1 px-4 py-2.5 text-sm font-medium bg-purple-50 text-purple-700 border border-purple-200 rounded-lg hover:bg-purple-100 transition-colors flex items-center justify-center"
        >
          <svg
            className="w-4 h-4 mr-1.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
            />
          </svg>
          Courses
        </button>
        <button
          type="button"
          onClick={() => onEdit(teacher)}
          className="flex-1 px-4 py-2.5 text-sm font-medium bg-blue-50 text-blue-700 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors flex items-center justify-center"
        >
          <svg
            className="w-4 h-4 mr-1.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
            />
          </svg>
          Edit
        </button>
        <button
          type="button"
          onClick={() => onDelete(teacher)}
          className="flex-1 px-4 py-2.5 text-sm font-medium bg-red-50 text-red-700 border border-red-200 rounded-lg hover:bg-red-100 transition-colors flex items-center justify-center"
        >
          <svg
            className="w-4 h-4 mr-1.5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
            />
          </svg>
          Delete
        </button>
      </div>

      {/* Courses Section (Expandable) */}
      {showCourses && (
        <div className="mt-4 pt-4 border-t border-gray-100">
          <h4 className="text-sm font-semibold text-gray-700 mb-3">Assigned Courses (2025-2026)</h4>
          {loading ? (
            <p className="text-sm text-gray-500">Loading courses...</p>
          ) : courses.length === 0 ? (
            <p className="text-sm text-gray-500 italic">No courses assigned</p>
          ) : (
            <div className="space-y-2">
              {courses.map(course => (
                <div key={course.id} className="bg-gray-50 rounded-lg p-3 border border-gray-200">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <p className="font-semibold text-sm text-gray-900">{course.course_code}</p>
                      <p className="text-xs text-gray-600 mt-0.5">{course.course_name}</p>
                    </div>
                    <div className="flex flex-col items-end gap-1 ml-2">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${course.role === 'Primary' ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'}`}>
                        {course.role}
                      </span>
                      <span className="text-xs text-gray-500">Sec: {course.section}</span>
                    </div>
                  </div>
                  <div className="mt-2 flex items-center text-xs text-gray-500">
                    <span className="mr-3">{course.course_type}</span>
                    <span>{course.credit} Credits</span>
                    <span className="mx-2">•</span>
                    <span>Sem {course.semester}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default TeacherCard;
