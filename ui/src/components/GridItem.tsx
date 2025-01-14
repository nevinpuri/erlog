import React, { useState, useCallback } from "react";
import { cvtReadable, timeConverter } from "../util";

function getColor(logLevel: any) {
  console.log('Log level:', logLevel, typeof logLevel);
  
  if (!logLevel) {
    return '';
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
    default:
      console.log('No match for log level:', logLevel);
      return "";
  }
}

function getLink(item) {
  console.log(item.parentId);
  if (item.parentId) {
    return `/${item.parentId}#${item.id}`;
  }

  return `/${item.id}`;
}

interface Props {
  item: Item;
  onHover: () => void;
}

interface Item {
  id: string;
  child_logs: number;
  timestamp: string;
  log: string;
}

interface LogPreview {
  timestamp: string;
  childCount: number;
  fields: { key: string; value: string }[];
}

export default function GridItem({ item, onHover }: Props) {
  const [expanded, setExpanded] = useState(false);
  const [childLogs, setChildLogs] = useState<Item[] | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const log = JSON.parse(item.log);
  const bgColor = getColor(log.level);

  const fetchChildLogs = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`http://localhost:8000/children/${item.id}`);
      if (!response.ok) {
        throw new Error('Failed to fetch children');
      }
      const data = await response.json();
      setChildLogs(data);
    } catch (error) {
      console.error('Failed to fetch child logs:', error);
      setChildLogs([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleExpand = async (e) => {
    e.preventDefault();
    if (!expanded && !childLogs) {
      await fetchChildLogs();
    }
    setExpanded(!expanded);
  };

  return (
    <div className="w-full relative">
      <div 
        className={`flex justify-between items-center w-full px-2 ${bgColor}`}
        onMouseEnter={onHover}
      >
        <div className="flex items-center flex-1 min-w-0">
          {item.child_logs > 0 && (
            <div className="flex items-center">
              <span className="text-xs badge badge-sm badge-neutral mr-2">
                {item.child_logs}
              </span>
              <button
                onClick={handleExpand}
                className="btn btn-ghost btn-xs btn-square"
                disabled={isLoading}
              >
                {isLoading ? (
                  <span className="loading loading-spinner loading-xs"></span>
                ) : (
                  <svg
                    className={`w-4 h-4 transition-transform duration-200 ${expanded ? 'rotate-90' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                )}
              </button>
            </div>
          )}
          <a href={getLink(log)} key={item.id} className="link link-hover flex-1 min-w-0 ml-2">
            <span className="block truncate">{cvtReadable(log)}</span>
          </a>
        </div>
        <span className="ml-4 whitespace-nowrap">{timeConverter(item.timestamp)}</span>
      </div>
      
      {expanded && (
        <div className="ml-6 border-l-2 border-base-200 pl-2">
          {childLogs && childLogs.map(childLog => (
            <GridItem key={childLog.id} item={childLog} />
          ))}
        </div>
      )}
    </div>
  );
}
