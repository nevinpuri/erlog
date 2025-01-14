import { LogView } from "../pages/logid";
import { cvtReadable, timeConverter } from "../util/index";
import React, { useEffect } from "react";
import GridItem from "./GridItem";

interface IProps {
  logs: any;
}

export default function Grid({ logs }: IProps) {
  return (
    <div className="w-full">
      {logs.map((log) => (
        <div key={log.id} className="w-full">
          <GridItem item={log} />
        </div>
      ))}
    </div>
  );
}
