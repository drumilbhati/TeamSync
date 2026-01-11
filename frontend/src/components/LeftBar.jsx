import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Card, CardDescription, CardHeader } from "./ui/card";
import { useTeam } from "@/context/TeamContext";
import { Button } from "./ui/button";
import CreateTeamForm from "./CreateTeamForm";

const LeftBar = () => {
  const [teams, setTeams] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { user } = useAuth();
  const { setSelectedTeam } = useTeam();
  const token = user?.token;

  const fetchTeams = async () => {
    if (!token) return;
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

  useEffect(() => {
    fetchTeams();
  }, [token]);

  return (
    <div>
      {isModalOpen ? (
        <div>
          <CreateTeamForm
            onTeamCreated={() => {
              fetchTeams();
              setIsModalOpen(false);
            }}
          />
        </div>
      ) : (
        <>
          <div
            className="
        flex flex-col gap-4"
          >
            <Button
              onClick={() => {
                setIsModalOpen((prev) => !prev);
              }}
            >
              Create new team
            </Button>
          </div>
        </>
      )}
      <div className="flex flex-col gap-4">
        {teams.map((team) => (
          <Card
            key={team.team_id}
            onClick={() => {
              setSelectedTeam(team);
            }}
            className="cursor-pointer"
          >
            <div>
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
    </div>
  );
};

export default LeftBar;
