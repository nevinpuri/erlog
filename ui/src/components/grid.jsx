import { timeConverter } from "../util/index";

export default function Grid({ logs }) {
  return (
    <div className="flex flex-col">
      {logs.map((log) => (
        <a
          href={`/${log.id}`}
          key={log.id}
          className="even:bg-gray-100 hover:bg-gray-50 px-2"
        >
          <span className="float-left">{log.log}</span>
          <span className="float-right">{timeConverter(log.timestamp)}</span>
        </a>
      ))}
    </div>
  );
}
