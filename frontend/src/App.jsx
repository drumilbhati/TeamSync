import "./App.css";
import { Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import NavBar from "./components/Navbar";

const App = () => {
  return (
    <div>
      <NavBar />
      <main className="content-main">
        <Routes>
          {/* add paths here */}
          <Route path="/" element={<Home />} />
        </Routes>
      </main>
    </div>
  );
};

export default App;
