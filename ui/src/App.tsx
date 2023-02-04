import { useEffect, useRef } from "react";
import useSWR from "swr";
import { getData } from "./types";
import { toFormattedDate } from "./utils";
import { v4 as uuid } from "uuid";

function App() {
  let listener: any;
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    listener = window.addEventListener("keydown", (event) => {
      if (event.key !== "/") {
        return;
      }

      if (!inputRef.current) {
        return;
      }

      event.preventDefault();
      inputRef.current.focus();
    });

    return function cleanup() {
      window.removeEventListener(listener, () => {});
    };
  }, []);

  function getLevelBg(level: string) {
    switch (level) {
      case "debug":
        return "bg-blue-100";
      case "warning":
        return "bg-yellow-100";
      case "error":
        return "bg-red-100";
      case "info":
        return "bg-gray-100";
      default:
        return "";
    }
  }
  const { data, error } = useSWR("/logs", getData);

  if (error) return <div>An error has occured</div>;

  if (!data) return <div>Loading</div>;

  return (
    <div>
      <div className="grid grid-rows-6 h-screen">
        <div className="row-span-9">
          {data.map((log) => (
            <div
              key={log.id}
              className={`flex space-x-4 ${getLevelBg(log.data.level)}`}
            >
              <h1>{toFormattedDate(log.createdAt)}</h1>
              {Object.entries(log.data).map((data) => (
                <h1 key={uuid()}>
                  <span>{data[0]}</span>
                  <span>=</span>
                  <span>{data[1] as any}</span>
                </h1>
              ))}
            </div>
          ))}
        </div>
        <div className="py-2 px-4">
          <input
            ref={inputRef}
            type="text"
            placeholder="Search"
            className="w-full border-2 rounded-md border-gray-400 focus:border-gray-800 focus:ring-gray-800 px-2 py-1.5"
          />
        </div>
      </div>
      {/* // <h1 className="text-2xl mb-2">Logs</h1> */}
    </div>
  );
}

export default App;
