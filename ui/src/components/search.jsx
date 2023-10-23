export default function Search({ onSubmit, onChange, value, onEnter }) {
  return (
    <form onSubmit={(e) => e.preventDefault()}>
      <input
        type="text"
        onChange={onChange}
        className="ring-1 ring-gray-300 w-full"
        placeholder="search"
        value={value}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onEnter(value);
          }
        }}
      ></input>
    </form>
  );
}
