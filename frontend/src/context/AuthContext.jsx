import { useState, createContext, useContext } from "react";
import { parseJwt } from "@/lib/utils";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(() => {
    try {
      const token = localStorage.getItem("token");
      if (token) {
        const decoded = parseJwt(token);
        return { token, ...decoded };
      }
    } catch (e) {
      console.error("Failed to restore session:", e);
      localStorage.removeItem("token");
    }
    return null;
  });
  
  // Loading is false because initialization is synchronous
  const [loading] = useState(false);

  const login = (token) => {
    localStorage.setItem("token", token);
    const decoded = parseJwt(token);
    setUser({ token, ...decoded });
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, loading }}>
      {!loading && children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
