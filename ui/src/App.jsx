import { useState } from "react";
import Search from "./components/search";
import { useEffect } from "react";
import Grid from "./components/grid";

const fetchLogs = async (query) => {
  const response = await fetch("http://localhost:8000/search", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ query }),
  });

  if (response.status == 400) {
    let text = await response.json();
    return { logs: null, err: text.detail };
  }

  const d = await response.json();
  return { logs: d, err: null };
};

function App() {
  const [v, setV] = useState("");

  const [err, setErr] = useState(null);
  const [logs, setLogs] = useState(null);

  useEffect(() => {
    const doWork = async () => {
      const { logs, err } = await fetchLogs("");
      setLogs(logs);
      setErr(err);
    };

    doWork();
  }, []);

  async function handleSubmit(e) {
    e.preventDefault();
    const { logs, err } = await fetchLogs(v);
    setLogs(logs);
    setErr(err);
  }

  if (!logs) {
    return (
      <div>
        <Search
          onSubmit={handleSubmit}
          onChange={(e) => setV(e.target.value)}
          value={v}
        />
        <p className="text-red-500 text-sm">{err}</p>
      </div>
    );
  }
  return (
    <div className="overflow-x-none">
      <Search
        onSubmit={handleSubmit}
        onChange={(e) => setV(e.target.value)}
        value={v}
      />
      <p className="text-red-500 text-sm">{err}</p>
      <Grid logs={logs} />
      {/* <div className="flex flex-col">
        {logs.map((log) => (
          <a href={`/${log.id}`} key={log.id}>
            {log.log}
          </a>
        ))}
      </div> */}
    </div>
  );
}

export default App;
