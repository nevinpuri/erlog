import { useEffect, useState } from "react";
import useSWR from "swr";
import { getData } from "./types";
import { toFormattedDate } from "./utils";

function App() {
  const { data, error } = useSWR("/logs", getData);

  if (error) return <div>An error has occured</div>;

  if (!data) return <div>Loading</div>;

  return (
    <div className="App">
      <h1 className="text-2xl">hello</h1>
      {data.map((log) => (
        <div key={log.id} className="flex space-x-4">
          <h1>{toFormattedDate(log.createdAt)}</h1>
          <h1>{log.data.level}</h1>
        </div>
      ))}
    </div>
  );
}

export default App;
