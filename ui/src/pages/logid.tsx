import { useLoaderData } from "react-router-dom";
import BackButton from "../components/back";
import { timeConverter } from "../util";
import React from "react";

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
          <div key={c}>
            <a className="text-lg font-medium link" href={`/${c.id}`}>
              {timeConverter(c.timestamp)}
            </a>
            {/* <PrettyPrintJson data={JSON.parse(c.log)} /> */}
            <LogView data={JSON.parse(c.log)} />
          </div>
        ))}
      </div>
    </div>
  );
}

interface ILogViewProps {
  data: any;
}

export const LogView = ({ data }: ILogViewProps) => {
  return (
    <table className="overflow-x-scroll outline outline-neutral">
      <thead className="border-b border-b-black">
        <tr>
          <th className="bg-base-300">Field</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody className="">
        {Object.entries(data).map(([key, value]) => (
          <>
            {key !== "id" && key !== "timestamp" ? (
              <tr key={key} className="">
                <td className="font-bold bg-base-300 px-2 py-0.5">{key}</td>
                <td className="pl-2 pr-2">{value}</td>
              </tr>
            ) : (
              <></>
            )}
          </>
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
