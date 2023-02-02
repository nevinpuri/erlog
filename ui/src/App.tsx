import { useEffect, useState } from "react";
import useSWR from "swr";
import { getData } from "./types";
import { toFormattedDate } from "./utils";

function App() {
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
    <div className="App">
      <h1 className="text-2xl mb-2">Logs</h1>
      {data.map((log) => (
        <div
          key={log.id}
          className={`flex space-x-4 ${getLevelBg(log.data.level)}`}
        >
          <h1>{toFormattedDate(log.createdAt)}</h1>
          {Object.entries(log.data).map((data) => (
            <h1>
              <span>{data[0]}</span>
              <span>=</span>
              <span>{data[1] as any}</span>
            </h1>
          ))}
        </div>
      ))}
    </div>
  );
}

export default App;
