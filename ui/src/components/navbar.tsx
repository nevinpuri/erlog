export default function Navbar() {
  return (
    <div className="">
      <div className="flex flex-row justify-between mx-2 my-1.5">
        <div className="my-auto">
          <h1>Logo</h1>
        </div>
        <div className="flex flex-row space-x-4">
          <h1 className="font-normal text-md text-gray-500 hover:text-gray-800 cursor-pointer">
            Logs
          </h1>
          <h1 className="font-normal text-md text-gray-500 hover:text-gray-800 cursor-pointer">
            Metrics
          </h1>
        </div>
      </div>
    </div>
  );
}
