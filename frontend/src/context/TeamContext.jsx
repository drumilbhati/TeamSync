import { useContext, useState, createContext } from "react";

const TeamContext = createContext();

export const TeamProvider = ({ children }) => {
  const [selectedTeam, setSelectedTeam] = useState();

  return (
    <TeamContext.Provider value={{ selectedTeam, setSelectedTeam }}>
      {children}
    </TeamContext.Provider>
  );
};

export const useTeam = () => useContext(TeamContext);
