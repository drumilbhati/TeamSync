import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { Button } from "./ui/button";
import { useTheme } from "./theme-provider";
import { Sun, Moon, LogOut } from "lucide-react";

const NavBar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();

  return (
    <nav className="border-b border-border bg-background px-6 py-4 flex justify-between items-center sticky top-0 z-50">
      <div className="font-mono text-xl font-bold tracking-tighter">
        <NavLink to="/" className="text-foreground hover:text-primary transition-colors">TeamSync</NavLink>
      </div>
      
      <div className="flex items-center gap-4">
        <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className="rounded-full"
            title="Toggle Theme"
        >
            {theme === "dark" ? (
                <Sun className="h-5 w-5 transition-all" />
            ) : (
                <Moon className="h-5 w-5 transition-all" />
            )}
            <span className="sr-only">Toggle theme</span>
        </Button>

        {user ? (
          <Button variant="ghost" className="gap-2" onClick={logout}>
             <LogOut className="w-4 h-4" />
             Logout
          </Button>
        ) : (
          <Button variant="default" onClick={() => navigate("/login")}>
            Login
          </Button>
        )}
      </div>
    </nav>
  );
};

export default NavBar;
