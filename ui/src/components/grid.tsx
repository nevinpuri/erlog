import { LogView } from "../pages/logid";
import { cvtReadable, timeConverter } from "../util/index";
import React, { useState } from "react";
import GridItem from "./GridItem";

interface LogPreview {
  timestamp: string;
  childCount: number;
  fields: { key: string; value: string }[];
}

interface IProps {
  logs: any;
}

export default function Grid({ logs }: IProps) {
  const [showPreview, setShowPreview] = useState(false);
  const [previewData, setPreviewData] = useState<LogPreview | null>(null);
  const [isLoadingPreview, setIsLoadingPreview] = useState(false);
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const [activeItemId, setActiveItemId] = useState<string | null>(null);

  const handleMouseMove = (e: React.MouseEvent) => {
    setMousePosition({ x: e.clientX, y: e.clientY });
  };

  const handleMouseLeave = () => {
    setShowPreview(false);
    setActiveItemId(null);
  };

  const fetchPreviewData = async (id: string) => {
    if (id === activeItemId) return;
    setActiveItemId(id);
    setIsLoadingPreview(true);
    try {
      const response = await fetch(`http://localhost:8000/preview/${id}`);
      if (!response.ok) throw new Error('Failed to fetch preview');
      const data = await response.json();
      
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
      setShowPreview(true);
    } catch (error) {
      console.error('Failed to fetch preview:', error);
    } finally {
      setIsLoadingPreview(false);
    }
  };

  return (
    <div 
      className="w-full relative" 
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
    >
      {logs.map((log) => (
        <div key={log.id} className="w-full">
          <GridItem 
            item={log} 
            onHover={() => fetchPreviewData(log.id)}
          />
        </div>
      ))}

      {showPreview && (
        <div 
          className="fixed z-50 w-96 pointer-events-none"
          style={{ 
            left: `${mousePosition.x + 12}px`,
            top: `${mousePosition.y + 8}px`,
            transform: 'translate(0, 0)'
          }}
        >
          <div className="bg-gray-800/95 rounded-lg border-2 border-gray-600/50 shadow-lg p-4 backdrop-blur-[2px] before:absolute before:inset-0 before:bg-noise before:opacity-[0.15] before:mix-blend-overlay before:contrast-150 relative overflow-hidden">
            {isLoadingPreview ? (
              <div className="flex items-center justify-center h-24">
                <div className="loading loading-spinner loading-lg text-white"></div>
              </div>
            ) : previewData && (
              <div className="space-y-3 animate-fadeIn text-white">
                <div className="flex justify-between items-center border-b border-gray-600/30 pb-2">
                  <span className="text-base font-semibold text-white">
                    {timeConverter(previewData.timestamp)}
                  </span>
                  {previewData.childCount > 0 && (
                    <span className="badge badge-md bg-gray-700/80 text-white border-gray-600">
                      {previewData.childCount} children
                    </span>
                  )}
                </div>
                <div className="space-y-2">
                  {previewData.fields.map(({ key, value }, index) => (
                    <div key={index} className="grid grid-cols-[auto,1fr] gap-3 text-base">
                      <span className="font-bold text-white">{key}:</span>
                      <span className="truncate text-white/95">{value}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
