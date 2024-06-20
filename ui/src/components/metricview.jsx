import {
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { LineChart } from "@tremor/react";

export default function MetricView({ title, data, per }) {
  return (
    <div>
      <LineChart
        data={data}
        index="dateTime"
        categories={["count"]}
        colors={["emerald"]}
        yAxisWidth={60}
        className="h-96 w-full"
      ></LineChart>
      {/* <ResponsiveContainer width="95%" height={400}>
        <LineChart
          // width={400}
          // height={200}
          data={data}
          // margin={{ top: 5, right: 5, bottom: 5, left: 5 }}
        >
          <Line type="monotone" dataKey="count" stroke="#8884d8" />
          <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
          <XAxis dataKey="dateTime" />
          <YAxis width={20} dataKey="count" />
          <Tooltip />
        </LineChart>
      </ResponsiveContainer> */}
      {/* <h2>Send a POST requests with</h2>
    <p>/metrics "id": "cpu_usage", value: "ifjasdiojf"</p> */}
    </div>
  );
}
