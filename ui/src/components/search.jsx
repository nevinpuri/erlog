import { useState } from "react";
import TimeFilter from "./TimeFilter";

export default function Search({
  onSubmit,
  onChange,
  defaultValue,
  onEnter,
  onShowChildrenChange,
  showChildren,
  timeRange,
  onTimeRangeChange,
}) {
  const [v, setV] = useState("");

  const handleTimeRangeChange = (range) => {
    if (onTimeRangeChange) {
      onTimeRangeChange(range);
    }
  };

  const handleShowChildrenChange = (checked) => {
    if (onShowChildrenChange) {
      onShowChildrenChange(checked);
    }
  };

  return (
    <form className="grid grid-cols-12 gap-4" onSubmit={(e) => e.preventDefault()}>
      <label className="input input-bordered w-full col-span-6 flex items-center">
        <input
          type="text"
          onChange={(e) => {
            setV(e.target.value);
            if (onChange) {
              onChange(e);
            }
          }}
          className="grow"
          placeholder="search"
          defaultValue={defaultValue}
          onKeyDown={(e) => {
            if (e.key === "Enter" && onEnter) {
              onEnter(v);
            }
          }}
        />
        <div className="tooltip tooltip-bottom" data-tip="Query help">
          <a
            href="https://github.com/nevinpuri/erlog?tab=readme-ov-file#querying"
            target="_blank"
            rel="noreferrer"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="h-4 w-4 opacity-70"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z"
              />
            </svg>
          </a>
        </div>
      </label>

      <div className="col-span-2">
        <TimeFilter
          value={timeRange}
          onChange={handleTimeRangeChange}
        />
      </div>

      <div className="col-span-4 flex items-center">
        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">Show child logs</span>
            <input
              type="checkbox"
              className="checkbox ml-2"
              onChange={(e) => handleShowChildrenChange(e.target.checked)}
              checked={showChildren}
            />
          </label>
        </div>
      </div>
    </form>
  );
}
