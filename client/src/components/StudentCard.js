import React from "react";
import { API_BASE_URL } from "../config";

function StudentCard({ student, onEdit, onDelete }) {
  // Get base URL without /api for static files (if student images are supported in future)
  const baseUrl = API_BASE_URL.replace("/api", "");
  // const imageUrl = student.profile_img ? `${baseUrl}${student.profile_img}` : null;
  const imageUrl = null; // Placeholder until profile images are implemented

  return (
    <div className="p-6 border rounded-lg bg-white shadow-sm hover:shadow-md transition-shadow">
      {/* Avatar and Name */}
      <div className="flex items-center space-x-3 mb-4">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={student.student_name}
            className="w-16 h-16 rounded-full object-cover shadow-md flex-shrink-0 border-2 border-white"
            onError={(e) => {
              e.target.style.display = "none";
              e.target.nextElementSibling.style.display = "flex";
            }}
          />
        ) : null}
        <div
          className="w-16 h-16 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-semibold text-2xl shadow-md flex-shrink-0"
          style={{ display: imageUrl ? "none" : "flex" }}
        >
          {student.student_name ? student.student_name.charAt(0).toUpperCase() : "S"}
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="text-xl font-bold text-gray-900 truncate">
            {student.student_name || "—"}
          </h3>
          <p className="text-sm text-gray-600 mt-1">
            {student.enrollment_no || "—"}
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
              d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2"
            />
          </svg>
          <div className="flex-1 min-w-0">
            <p className="text-xs text-gray-500 uppercase font-medium">Student ID</p>
            <p className="text-sm text-gray-900 truncate" title={student.student_id}>
              {student.student_id || student.id || "—"}
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
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
            />
          </svg>
          <div className="flex-1">
            <p className="text-xs text-gray-500 uppercase font-medium">Age</p>
            <p className="text-sm text-gray-900">{student.age || "—"}</p>
          </div>
        </div>

        {/* Optional: Add Department or Batch if available in student object list view */}
        {student.department && (
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
              <p className="text-sm text-gray-900">{student.department}</p>
            </div>
          </div>
        )}
      </div>

      {/* Status Badge */}
      <div className="mb-4">
        {student.status === 1 && (
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
          onClick={() => onEdit(student.student_id || student.id)}
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
        {onDelete && (
        <button
          type="button"
          onClick={() => onDelete(student.student_id || student.id)}
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
        )}
      </div>
    </div>
  );
}

export default StudentCard;
