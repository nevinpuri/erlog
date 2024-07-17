import { LogView } from "../pages/logid";
import { cvtReadable, timeConverter } from "../util/index";
import React from "react";

interface IProps {
  logs: any;
}

export default function Grid({ logs }: IProps) {
  return (
    <>
      <div className="table table-zebra table-lg">
        <tbody>
          {logs.map((log) => (
            <tr className="px-2">
              <details className="px-2 w-full">
                <summary className="flex justify-between">
                  <a
                    href={`/${log.id}`}
                    key={log.id}
                    className="link link-hover"
                  >
                    {cvtReadable(JSON.parse(log.log))}
                  </a>

                  <span className="mr-2">{timeConverter(log.timestamp)}</span>
                </summary>
                <>
                  <LogView data={JSON.parse(log.log)} />
                  <a href={`/${log.id}`} className="link link-info">
                    View Children
                  </a>
                </>
              </details>
            </tr>
          ))}
        </tbody>
      </div>
    </>
  );
}
