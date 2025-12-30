import { NavLink } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";

const NavBar = () => {
  const { user, logout } = useAuth();

  return (
    <nav className="font-mono text-amber-50 text-2xl flex justify-center">
      <div>
        <NavLink to="/">TeamSync</NavLink>
      </div>
      <div>
        <NavLink to="/">Home</NavLink>
      </div>
      {user ? (
        <div className="cursor-pointer" onClick={logout}>
          Logout
        </div>
      ) : (
        <div>
          <NavLink to="/login">Login</NavLink>
        </div>
      )}
    </nav>
  );
};

export default NavBar;
