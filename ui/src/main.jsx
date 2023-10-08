import "./index.css";
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import { LogId, loader as logIdLoader } from "./pages/logid";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { Metrics } from "./pages/metrics";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
  },
  {
    path: "/:id",
    element: <LogId />,
    loader: logIdLoader,
  },
  {
    path: "/metrics",
    element: <Metrics />,
  },
]);

const selected = (route) => {
  console.log(window.location.href);
  if (window.location.pathname === route) {
    return "text-blue-600";
  }

  return "";
};

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <div className="w-full flex flex-row">
      <div className="w-16 flex flex-col items-center pt-4 space-y-6">
        <a
          href="/"
          className={`inline-flex flex-col items-center ${selected("/")}`}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
            />
          </svg>
          <span className="mt-1 text-sm">Logs</span>
        </a>

        <a
          href="/metrics"
          className={`inline-flex flex-col items-center ${selected(
            "/metrics"
          )}`}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
            />
          </svg>
          <span className="mt-1 text-sm">Metrics</span>
        </a>
      </div>
      <div className="w-full">
        <RouterProvider router={router} />
      </div>
    </div>

    {/* <App /> */}
  </React.StrictMode>
);
