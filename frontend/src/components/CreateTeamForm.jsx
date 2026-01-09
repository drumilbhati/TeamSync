import { Input } from "./ui/input";
import { useState } from "react";
import { Button } from "./ui/button";
import { useAuth } from "@/context/AuthContext";

const CreateTeamForm = ({ onTeamCreated }) => {
  const [teamName, setTeamName] = useState("");
  const { user } = useAuth();
  const token = user?.token;

  const createTeam = async () => {
    try {
      const response = await fetch("/api/teams", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ team_name: teamName }),
      });

      if (response.ok) {
        const data = await response.json();
        setTeamName("");
        if (onTeamCreated) {
          onTeamCreated();
        }
      }
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div className="flex flex-col gap-4">
      Team Name
      <Input
        type="text"
        value={teamName}
        onChange={(e) => setTeamName(e.target.value)}
        placeholder="Team Name"
      ></Input>
      <Button onClick={createTeam}>Create Team</Button>
    </div>
  );
};

export default CreateTeamForm;
