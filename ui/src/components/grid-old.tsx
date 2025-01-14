import { LogView } from "../pages/logid";
import { cvtReadable, timeConverter } from "../util/index";
import React from "react";

interface IProps {
  logs: any;
}

export default function GridOld({ logs }: IProps) {
  return (
    <>
      <div className="table table-lg">
        <tbody>
          {logs.map((log) => (
            <tr className="px-2">
              <details className="px-2 w-full">
                <summary className="flex justify-between">
                  <a
                    href={`/${log.id}`}
                    key={log.id}
                    className="link link-hover inline-flex bg-gray-100 m-2 my-1 p-2 px-4 rounded-md justify-center"
                  >
                    <span className="text-4xl">🎓</span>
                    <div className="flex justify-center flex-col ml-4">
                      <p>
                        {cvtReadable(JSON.parse(log.log))}{" "}
                        <span className="text-sm badge ml-1.5 badge-sm border-gray-400">
                          +1
                        </span>
                      </p>
                      <span className="mr-2">
                        {timeConverter(log.timestamp)}
                      </span>
                    </div>
                  </a>
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
