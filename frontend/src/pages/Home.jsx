import LeftBar from "@/components/LeftBar";
import RightBar from "@/components/RightBar";

const Home = () => {
  return (
    <div>
      <div className="max-w-dvw font-mono grid gap-2 sm:grid-cols-12 min-h-dvh">
        <div className="sm:col-span-3 rounded-lg bg-amber-500 shadow sm:block hidden">
          <LeftBar />
        </div>
        <div className="sm:col-span-7 rounded-lg bg-blue-500 shadow sm:block"></div>
        <div className="sm:col-span-2 rounded-lg bg-purple-500 shadow sm:block hidden">
          <RightBar></RightBar>
        </div>
      </div>
    </div>
  );
};

export default Home;
