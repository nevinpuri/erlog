import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

export default function MetricView({ title, data, per }) {
  return (
    <div>
      <h1>{title}</h1>
      <LineChart
        width={400}
        height={200}
        data={data}
        // margin={{ top: 5, right: 5, bottom: 5, left: 5 }}
      >
        <Line type="monotone" dataKey="count" stroke="#8884d8" />
        <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
        <XAxis dataKey="dateTime" />
        <YAxis width={20} dataKey="count" />
        <Tooltip />
      </LineChart>
      {/* <h2>Send a POST requests with</h2>
    <p>/metrics "id": "cpu_usage", value: "ifjasdiojf"</p> */}
    </div>
  );
}
