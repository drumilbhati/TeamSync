import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Card, CardDescription, CardHeader } from "./ui/card";

const LeftBar = () => {
  const [teams, setTeams] = useState([]);
  const { user } = useAuth();
  const token = user?.token;
  useEffect(() => {
    const fetchTeams = async () => {
      try {
        const response = await fetch("/api/teams", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          const data = await response.json();
          setTeams(data);
        }
      } catch (error) {
        console.log(error);
      }
    };

    fetchTeams();
  }, [token]);

  return (
    <div className="flex flex-col gap-4">
      {teams.map((team) => (
        <Card>
          <div key={team.team_id}>
            <CardHeader>
              <h3 className="text-teal-950!">{team.team_name}</h3>
            </CardHeader>
            <CardDescription>
              <p className="text-teal-950!">
                Team Leader ID: {team.team_leader_id}
              </p>
            </CardDescription>
          </div>
        </Card>
      ))}
    </div>
  );
};

export default LeftBar;
