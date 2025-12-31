import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Button } from "./ui/button";

const NavBar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  return (
    <nav className="font-mono text-amber-50 text-2xl flex justify-between">
      <div>
        <NavLink to="/">TeamSync</NavLink>
      </div>
      {user ? (
        <div>
          <div className="text-[18px] cursor-pointer" onClick={logout}>
            Logout
          </div>
        </div>
      ) : (
        <div>
          <div
            className="text-[18px] cursor-pointer"
            onClick={() => navigate("/login")}
          >
            Login
          </div>
        </div>
      )}
    </nav>
  );
};

export default NavBar;
