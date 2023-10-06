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
      <PrettyPrintJson data={JSON.parse(log.log)} />
    </div>
  );
}

const PrettyPrintJson = ({ data }) => {
  // (destructured) data could be a prop for example
  return (
    <div className="text-sm">
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};
