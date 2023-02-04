export default function Navbar() {
  return (
    <div className="flex space-x-2">
      <h1 className="font-medium text-xl text-gray-500 hover:text-gray-800 cursor-pointer">
        Logs
      </h1>
      <h1 className="font-medium text-xl text-gray-500 hover:text-gray-800 cursor-pointer">
        Metrics
      </h1>
    </div>
  );
}
