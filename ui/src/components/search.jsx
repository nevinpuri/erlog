import { useState } from "react";

export default function Search({ onSubmit, onChange, defaultValue, onEnter }) {
  const [v, setV] = useState("");
  return (
    <form className="grid grid-cols-8" onSubmit={(e) => e.preventDefault()}>
      <input
        type="text"
        onChange={(e) => {
          setV(e.target.value);
          if (!onChange) {
            return;
          }

          onChange(e);
        }}
        className="ring-1 ring-gray-300 w-full col-span-5"
        placeholder="search"
        defaultValue={defaultValue}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onEnter(v);
          }
        }}
      ></input>
    </form>
  );
}
