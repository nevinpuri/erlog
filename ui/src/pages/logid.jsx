import { useLoaderData } from "react-router-dom";
import BackButton from "../components/back";
import { timeConverter } from "../util";

const fetchLog = async (id) => {
  const response = await fetch("http://localhost:8000/get", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id }),
  });

  const d = await response.json();
  return { log: d };
};

export async function loader({ params }) {
  const { log } = await fetchLog(params.id);
  return { log };
}

export function LogId() {
  const { log } = useLoaderData();

  return (
    <div>
      <BackButton />
      <h1 className="text-lg font-medium">{timeConverter(log.timestamp)}</h1>
      <br />
      <LogView data={JSON.parse(log.log)} />
      {/* <PrettyPrintJson data={JSON.parse(log.log)} /> */}
      <h1 className="font-normal text-lg my-2">Children</h1>
      <div className="space-y-4">
        {log.children.map((c) => (
          <div>
            <h1 className="text-lg font-medium">
              {timeConverter(c.timestamp)}
            </h1>
            {/* <PrettyPrintJson data={JSON.parse(c.log)} /> */}
            <LogView data={JSON.parse(c.log)} />
          </div>
        ))}
      </div>
    </div>
  );
}

const LogView = ({ data }) => {
  return (
    <table className="">
      <tbody className="align-baseline">
        {Object.entries(data).map(([key, value]) => (
          <tr>
            <td className="whitespace-nowrap pr-2 bg-gray-100 rounded-sm">
              {key}
            </td>
            <td className="pl-2">{value}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

const PrettyPrintJson = ({ data }) => {
  // (destructured) data could be a prop for example
  return (
    <div className="text-sm">
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};
