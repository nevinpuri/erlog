import React from "react";
import { cvtReadable, timeConverter } from "../util";

function getColor(logLevel: any) {
  if (!logLevel) {
    return;
  }

  switch (logLevel.toString().toLowerCase()) {
    case "info":
      return "bg-green-400/20";
    case "warning":
      return "bg-yellow-400/20";
    case "error":
      return "bg-red-400/20";
    case "critical":
      return "bg-red-400/40";
    case "debug":
      return "bg-purple-400/20";
  }
  return "";
  console.log(logLevel);
}

function getLink(item) {
  console.log(item.parentId);
  if (item.parentId) {
    return `/${item.parentId}#${item.id}`;
  }

  return `/${item.id}`;
}

export default function GridItem({ item }) {
  const log = JSON.parse(item.log);
  return (
    <div className={`flex justify-between items-center w-full px-2 ${getColor(log.level)}`}>
      <a href={getLink(log)} key={item.id} className="link link-hover flex-1 min-w-0">
        <span className="block truncate">{cvtReadable(log)}</span>

        {item.child_logs > 0 && (
          <span className="ml-1.5 badge badge-neutral badge-xs">
            {item.child_logs}
          </span>
        )}
      </a>

      <span className="ml-4 whitespace-nowrap">{timeConverter(item.timestamp)}</span>
    </div>
  );
}
