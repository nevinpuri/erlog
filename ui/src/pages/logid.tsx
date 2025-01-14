import { useLoaderData, useLocation, useParams } from "react-router-dom";
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
  const location = useLocation();
  const f = location.hash.substring(1);

  return (
    <div className="mx-4">
      <BackButton />
      <h1 className="text-lg font-medium">{timeConverter(log.timestamp)}</h1>
      <br />
      <LogView data={JSON.parse(log.log)} />
      {/* <PrettyPrintJson data={JSON.parse(log.log)} /> */}
      <h1 className="font-normal text-lg my-2">Children</h1>
      <div className="space-y-4">
        {log.children.map((c) => (
          <div key={c}>
            <a className={`text-lg font-medium`} href={`/${log.id}#${c.id}`}>
              <span className="link">{timeConverter(c.timestamp)}</span>
              {c.id === f ? (
                <span className="ml-2 no-underline">⬅</span>
              ) : (
                <></>
              )}
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
    <div id={data.id}>
      <table className="overflow-x-scroll border border-black">
        <thead className="border-b border-b-black">
          <tr>
            <th className="bg-base-300">Field</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody className="">
          {Object.entries(data).map(([key, value]) => (
            <>
              {key !== "timestamp" ? (
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
    </div>
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
