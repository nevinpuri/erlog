import { useEffect, useRef, useState } from "react";
import useSWR, { useSWRConfig } from "swr";
import { API_URL, ErLog, getData, postData } from "./types";
import { toFormattedDate } from "./utils";
import { v4 as uuid } from "uuid";
import axios, { AxiosResponse } from "axios";

function App() {
  let listener: any;
  const inputRef = useRef<HTMLInputElement>(null);
  const [query, setQuery] = useState<string>("");
  const { mutate } = useSWRConfig();

  useEffect(() => {
    listener = window.addEventListener("keydown", (event) => {
      if (event.key !== "/") {
        return;
      }

      if (!inputRef.current) {
        return;
      }

      event.preventDefault();
      inputRef.current.focus();
    });

    return function cleanup() {
      window.removeEventListener(listener, () => {});
    };
  }, []);

  function getLevelBg(level: string) {
    switch (level) {
      case "debug":
        return "bg-blue-100";
      case "warning":
        return "bg-yellow-100";
      case "error":
        return "bg-red-100";
      case "info":
        return "bg-gray-100";
      default:
        return "";
    }
  }
  const fetcher = async (url: string) => {
    const response: AxiosResponse<ErLog[]> = await axios.post(
      new URL(url, API_URL).href,
      {
        search: query,
      }
    );

    return response.data;
  };

  useEffect(() => {
    mutate("/search/logs");
  }, [query]);

  const { data, error } = useSWR("/search/logs", fetcher, {
    refreshInterval: 1000,
  });

  if (error) return <div>An error has occured</div>;

  if (!data) return <div>Loading</div>;

  return (
    <>
      <div className="">
        <div className="">
          {data.map((log) => (
            <div
              key={log.id}
              className={`flex space-x-4 ${getLevelBg(log.data.level)}`}
            >
              <h1>{toFormattedDate(log.createdAt)}</h1>
              {Object.entries(log.data).map((data) => (
                <h1 key={uuid()}>
                  <span>{data[0]}</span>
                  <span>=</span>
                  <span>{data[1] as any}</span>
                </h1>
              ))}
            </div>
          ))}
        </div>
        <div className="">
          <div className="py-2 px-4">
            <input
              ref={inputRef}
              onChange={(e) => setQuery(e.target.value)}
              type="text"
              placeholder="Search"
              className="w-full border-2 rounded-md border-gray-400 focus:border-gray-800 focus:ring-gray-800 px-2 py-1.5"
            />
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
