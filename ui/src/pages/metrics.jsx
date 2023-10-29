import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";
import MetricView from "../components/metricview";

const data = [
  { amt: 10, hour: 1 },
  { amt: 25, hour: 2 },
  { amt: 35, hour: 3 },
];

export function Metrics() {
  return (
    <div className="grid grid-cols-3">
      <MetricView title="Active Users Per Hour" data={data} />
      <MetricView title="Errors Per Hour" data={data} />
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
