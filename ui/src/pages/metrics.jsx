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

const data = [
  { amt: 10, hour: 1 },
  { amt: 25, hour: 2 },
  { amt: 35, hour: 3 },
];

export function Metrics() {
  const [dHour, setDHour] = useState();
  const [dDay, setDDay] = useState();

  useEffect(() => {
    const doWork = () => {
      fetch("http://localhost:8000/search", {
        method: "POST",
        body: JSON.stringify({
          per: "hour",
        }),
      }).then(async (e) => {
        const d = await e.json();
        console.log(d);
        setDHour(d);
      });

      fetch("http://localhost:8000/search", {
        method: "POST",
        body: JSON.stringify({
          per: "day",
        }),
      }).then(async (e) => {
        const d = await e.json();
        console.log(d);
        setDDay(d);
      });
    };

    doWork();
  }, []);
  return (
    <div className="grid grid-cols-3">
      <MetricView title="Active Users Per Hour" data={dHour} per="hour" />
      <MetricView title="Active users per day" data={dDay} per="day" />
      <MetricView title="Errors Per Hour" data={data} />
      <MetricView title="Errors Per Hour" data={data} />
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
