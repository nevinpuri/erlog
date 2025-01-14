import React from 'react';

export type TimeRange = 'all' | '1h' | '24h' | '7d' | '30d';

interface TimeFilterProps {
  value: TimeRange;
  onChange: (range: TimeRange) => void;
}

export default function TimeFilter({ value, onChange }: TimeFilterProps) {
  return (
    <select 
      className="select select-bordered select-sm w-40"
      value={value}
      onChange={(e) => onChange(e.target.value as TimeRange)}
    >
      <option value="all">All time</option>
      <option value="1h">Last hour</option>
      <option value="24h">Last 24 hours</option>
      <option value="7d">Last 7 days</option>
      <option value="30d">Last 30 days</option>
    </select>
  );
} 