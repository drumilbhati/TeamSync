import { NavLink } from "react-router-dom";

const NavBar = () => {
  return (
    <nav className="font-mono text-amber-50 text-4xl flex  items-center p-4 m-4 gap-8">
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
