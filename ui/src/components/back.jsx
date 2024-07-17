import { useNavigate } from "react-router-dom";

export default function BackButton() {
  const router = useNavigate();
  return (
    <button
      href="/"
      className="inline-flex hover:underline mt-2"
      onClick={() => router(-1)}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth={1.5}
        stroke="currentColor"
        className="w-6 h-6"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18"
        />
      </svg>
      <span className="ml-1.5">Back</span>
    </button>
  );
}
