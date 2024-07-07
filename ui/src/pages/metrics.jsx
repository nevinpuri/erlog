import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";
import MetricView from "../components/metricview";
import { useEffect, useState } from "react";

export function Metrics() {
  const [filter, setFilter] = useState("hour");
  const [data, setData] = useState();

  useEffect(() => {
    fetch("http://localhost:8000/search", {
      method: "POST",
      body: JSON.stringify({
        per: filter,
      }),
    }).then(async (e) => {
      const d = await e.json();
      console.log(d);
      setData(d);
    });
  }, [filter]);

  // fuck me, this whole add shit to that object shit is really really good
  return (
    <div>
      <div className="mx-20 pt-2">
        <h1>Metrics</h1>
        <h1>
          Follow my{" "}
          <a
            href="https://x.com/nevinpuri"
            target="_blank"
            rel="noreferrer"
            className="link link-info"
          >
            Twitter
          </a>{" "}
          to stay updated
        </h1>
        {/* <div>
          <label htmlFor="filter">Filter: </label>
          <select
            name="filter"
            id="filter"
            onChange={(e) => setFilter(e.target.value)}
          >
            <option value="hour">hour</option>
            <option value="day">day</option>
          </select>
        </div> */}
      </div>
      {/* <MetricView title="Active Users Per Hour" data={data} per="hour" /> */}
      {/* <MetricView title="Active users per day" data={dDay} per="day" />
      <MetricView title="Errors Per Hour" data={data} />
      <MetricView title="Errors Per Hour" data={data} /> */}
      {/* <a
        href="/metrics/new"
        className="px-3.5 py-2 from-green-500 via-emerald-500 to-green-500 bg-gradient-to-tr hover:from-green-600 hover:via-emerald-600 hover:to-green-600 transition-all border-emerald-600 border font-semibold shadow-sm rounded-md text-white"
      >
        + New Metric
      </a> */}
    </div>
    // <div>
    //   <h1>Active Users Per Hour</h1>
    //   <LineChart
    //     width={400}
    //     height={200}
    //     data={data}
    //     margin={{ top: 5, right: 5, bottom: 5, left: -35 }}
    //   >
    //     <Line type="monotone" dataKey="amt" stroke="#8884d8" />
    //     <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
    //     <XAxis dataKey="hour" />
    //     <YAxis />
    //     <Tooltip />
    //   </LineChart>
    // </div>
  );
}

/* <h2>Send a POST requests with</h2>
      <p>/metrics "id": "cpu_usage", value: "ifjasdiojf"</p> */
