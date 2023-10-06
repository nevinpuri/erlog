export default function Search({ onSubmit, onChange, value }) {
  return (
    <form onSubmit={onSubmit}>
      <input
        type="text"
        onChange={onChange}
        className="ring-1 ring-gray-300 w-full"
        placeholder="search"
        value={value}
      ></input>
    </form>
  );
}
