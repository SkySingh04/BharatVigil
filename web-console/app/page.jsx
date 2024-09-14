import Image from "next/image";
import DashBoard from "./Components/Dashboard";



export default function Home() {
  return (
    <div className="h-screen overflow-y-hiddden">
     <DashBoard />
    </div>
  );
}
