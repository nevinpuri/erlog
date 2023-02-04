import { PropsWithChildren } from "react";
import App from "./App";
import Navbar from "./components/navbar";

export default function Layout({ children }: PropsWithChildren) {
  return (
    <div className="h-screen flex flex-col">
      <Navbar />
      <App />
    </div>
  );
}
