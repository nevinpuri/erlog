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
              <details className="dropdown px-2 w-full">
                <summary className="">
                  <a
                    href={`/${log.id}`}
                    key={log.id}
                    className="link link-hover"
                  >
                    {/* <span className="float-left">{JSON.parse(log.log).event}</span> */}
                    <span>{cvtReadable(JSON.parse(log.log))}</span>
                    <span className="ml-8">{timeConverter(log.timestamp)}</span>
                  </a>
                </summary>
                <>
                  {/* {JSON.stringify(log)} */}
                  <LogView data={JSON.parse(log.log)} />
                  <a href="/" className="link link-info">
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
