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

export default function GridItem({ item }: Props) {
  const [expanded, setExpanded] = useState(false);
  const [childLogs, setChildLogs] = useState<Item[] | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [previewData, setPreviewData] = useState<LogPreview | null>(null);
  const [isLoadingPreview, setIsLoadingPreview] = useState(false);
  const log = JSON.parse(item.log);
  const bgColor = getColor(log.level);

  const fetchPreviewData = useCallback(async () => {
    setIsLoadingPreview(true);
    try {
      const response = await fetch(`http://localhost:8000/preview/${item.id}`);
      if (!response.ok) throw new Error('Failed to fetch preview');
      const data = await response.json();
      
      // Parse the log data into a more readable format
      const parsedLog = JSON.parse(data.log);
      const fields = Object.entries(parsedLog)
        .filter(([key]) => key !== 'timestamp' && key !== 'level')
        .map(([key, value]) => ({
          key,
          value: typeof value === 'object' ? JSON.stringify(value) : String(value)
        }));

      setPreviewData({
        timestamp: data.timestamp,
        childCount: data.child_logs,
        fields
      });
    } catch (error) {
      console.error('Failed to fetch preview:', error);
    } finally {
      setIsLoadingPreview(false);
    }
  }, [item.id]);

  const handleMouseEnter = () => {
    setShowPreview(true);
    if (!previewData) {
      fetchPreviewData();
    }
  };

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
    <div className="w-full relative group">
      <div 
        className={`flex justify-between items-center w-full px-2 ${bgColor} relative`}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={() => setShowPreview(false)}
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
        
        {showPreview && (
          <div className="absolute left-0 top-full mt-2 z-50 w-96 transform scale-95 opacity-0 group-hover:scale-100 group-hover:opacity-100 transition-all duration-200">
            <div className="bg-base-200 rounded-lg border-2 border-base-300 shadow-lg p-4 backdrop-blur-sm">
              {isLoadingPreview ? (
                <div className="flex items-center justify-center h-24">
                  <div className="loading loading-spinner loading-lg text-primary"></div>
                </div>
              ) : previewData && (
                <div className="space-y-2 animate-fadeIn">
                  <div className="flex justify-between items-center border-b border-base-300 pb-2">
                    <span className="text-sm font-medium">
                      {timeConverter(previewData.timestamp)}
                    </span>
                    {previewData.childCount > 0 && (
                      <span className="badge badge-primary badge-sm">
                        {previewData.childCount} children
                      </span>
                    )}
                  </div>
                  <div className="space-y-1">
                    {previewData.fields.map(({ key, value }, index) => (
                      <div key={index} className="grid grid-cols-[auto,1fr] gap-2 text-sm">
                        <span className="font-medium text-primary-content/70">{key}:</span>
                        <span className="truncate">{value}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
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
