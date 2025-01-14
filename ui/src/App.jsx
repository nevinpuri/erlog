import React, { useState } from "react";
import Search from "./components/search";
import { useEffect } from "react";
import Grid from "./components/grid";
import { useLocation, useNavigate } from "react-router-dom";
import TimeFilter from "./components/TimeFilter";

const fetchLogs = async (query, page, showChildren, timeRange) => {
  console.log(showChildren);
  console.log("SHOW CHILDREN");
  let q = query.replace(" and ", " AND ");
  q = q.replace(" or ", " OR ");
  const response = await fetch("http://localhost:8000/search", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ query: q, page, showChildren, timeRange }),
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

  const [err, setErr] = useState(null);
  const [logs, setLogs] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [showChildren, setShowChildren] = useState(() => {
    const c = q.get("children");
    return c === "true";
  });
  const [timeRange, setTimeRange] = useState(() => {
    return q.get("time") || "all";
  });

  // Handle URL updates when filters change
  useEffect(() => {
    const currentQuery = q.get("query") || "";
    const currentPage = q.get("page") || "0";
    router(`/?query=${currentQuery}&page=${currentPage}&children=${showChildren}&time=${timeRange}`);
  }, [showChildren, timeRange]);

  // Fetch logs when URL params change
  useEffect(() => {
    console.log("Fetching logs with params:", {
      query: q.get("query"),
      children: q.get("children"),
      page: q.get("page"),
      time: q.get("time")
    });
    f();
  }, [q.get("query"), q.get("children"), q.get("page"), q.get("time")]);

  const f = async () => {
    setIsLoading(true);
    try {
      let page = q.get("page");
      let query = q.get("query");
      let showChildren = q.get("children");
      let timeRange = q.get("time");

      if (!query) query = "";
      if (!page) page = 0;
      if (!showChildren) showChildren = false;
      if (!timeRange) timeRange = "all";

      const { logs, err } = await fetchLogs(
        query,
        page,
        showChildren,
        timeRange
      );
      setLogs(logs);
      setErr(err);
    } catch (error) {
      setErr("An error occurred while fetching logs");
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="overflow-x-none">
      <div>
        <Search
          defaultValue={q.get("query")}
          onEnter={(e) => {
            router(`/?query=${e}&page=0&children=${showChildren}&time=${timeRange}`);
          }}
          showChildren={showChildren}
          onShowChildrenChange={(e) => {
            setShowChildren(e);
          }}
          timeRange={timeRange}
          onTimeRangeChange={(range) => {
            setTimeRange(range);
          }}
        />
      </div>
      <p className="text-red-500 text-sm">{err}</p>
      {isLoading ? (
        <div className="flex justify-center items-center p-8">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      ) : (
        <Grid logs={logs || []} />
      )}
    </div>
  );
}

export default App;
