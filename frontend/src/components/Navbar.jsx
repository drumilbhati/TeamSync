import { NavLink } from "react-router-dom";

const NavBar = () => {
  return (
    <nav className="font-mono text-amber-50 text-2xl flex justify-center">
      <div>
        <NavLink to="/">TeamSync</NavLink>
      </div>
      <div>
        <NavLink to="/">Home</NavLink>
      </div>
      <div>
        <NavLink to="/login">Login</NavLink>
      </div>
    </nav>
  );
};

export default NavBar;
