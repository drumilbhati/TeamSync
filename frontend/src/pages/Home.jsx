import LeftBar from "@/components/LeftBar";
import MidBar from "@/components/MidBar";
import RightBar from "@/components/RightBar";

const Home = () => {
  return (
    <div className="bg-background min-h-screen w-full">
      <div className="grid gap-4 p-4 sm:grid-cols-12 h-[calc(100vh-2rem)]">
        <div className="sm:col-span-3 rounded-xl border bg-card shadow-sm overflow-hidden flex flex-col">
          <LeftBar />
        </div>
        <div className="sm:col-span-6 rounded-xl border bg-card shadow-sm overflow-hidden flex flex-col">
          <MidBar />
        </div>
        <div className="sm:col-span-3 rounded-xl border bg-card shadow-sm overflow-hidden flex flex-col">
          <RightBar />
        </div>
      </div>
    </div>
  );
};

export default Home;
