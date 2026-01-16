import { Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import NavBar from "./components/Navbar";
import Login from "./pages/Login";
import Register from "./pages/Register";

const App = () => {
  return (
    <div>
      <NavBar />
      <main className="content-main">
        <Routes>
          {/* add paths here */}
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </main>
    </div>
  );
};

export default App;
