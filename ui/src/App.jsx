import React, { useState } from "react";
import Search from "./components/search";
import { useEffect } from "react";
import Grid from "./components/grid";
import { useLocation, useNavigate } from "react-router-dom";

const fetchLogs = async (query, page) => {
  const response = await fetch("http://localhost:8000/search", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ query, page }),
  });

  if (response.status == 400) {
    let text = await response.json();
    return { logs: null, err: text.detail };
  }

  const d = await response.json();
  return { logs: d, err: null };
};

export function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

function App() {
  const router = useNavigate();
  const q = useQuery();
  // const [v, setV] = useState("");

  const [err, setErr] = useState(null);
  const [logs, setLogs] = useState(null);

  const f = async () => {
    let page = q.get("p");
    let query = q.get("query");
    if (!query) {
      query = "";
    }

    if (!page) {
      page = 0;
    }

    const { logs, err } = await fetchLogs(q.get("query"), q.get("page"));
    setLogs(logs);
    setErr(err);
  };

  useEffect(() => {
    if (!q.get("query" || !q.get("p"))) {
      router("/?query=&page=0");
    }
  }, []);

  useEffect(() => {
    console.log("QUERY CHNSGED!!!!");
    console.log(q.get("query"));
    f();
  }, [q.get("query")]);

  // useEffect(() => {
  //   // const doWork = async () => {
  //   //   let page = q.get("p");

  //   //   if (!page) {
  //   //     page = 0;
  //   //   }

  //   //   const { logs, err } = await fetchLogs(query, page);
  //   //   setLogs(logs);
  //   //   setErr(err);
  //   // };

  //   // doWork();
  //   f();
  // }, [query]);

  // useEffect(() => {
  // }, [q.get("query")]);

  // async function handleSubmit(e) {
  //   e.preventDefault();

  //   if (e.key === "Enter") {
  //     setQuery(v);
  //   }

  //   // refreshLogs();
  //   // const { logs, err } = await fetchLogs(v);
  //   // setLogs(logs);
  //   // setErr(err);
  // }

  if (!logs) {
    return (
      <div>
        <Search
          // onSubmit={handleSubmit}
          // onChange={(e) => setV(e.target.value)}
          defaultValue={q.get("query")}
          onEnter={(e) => {
            console.log("enter!!");
            console.log(e);
            // setQuery(e);
            // q.set("query", e);
            router(`/?query=${e}&p=${0}`);
          }}
        />
        <p className="text-red-500 text-sm">{err}</p>
      </div>
    );
  }
  return (
    <div className="overflow-x-none">
      <Search
        // onSubmit={handleSubmit}
        // onChange={(e) => setV(e.target.value)}
        defaultValue={q.get("query")}
        onEnter={(e) => {
          // setQuery(e);
          console.log("enter!!");
          console.log(e);
          router(`/?query=${e}&p=${0}`);
          // q.set("query", e);
        }}
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
