import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Card, CardHeader, CardTitle, CardDescription } from "./ui/card";
import { useTeam } from "@/context/TeamContext";
import { Button } from "./ui/button";
import CreateTeamForm from "./CreateTeamForm";
import EditTeamForm from "./EditTeamForm";
import AddMemberForm from "./AddMemberForm";
import { Plus, Users, Pencil, UserPlus } from "lucide-react";
import { cn } from "@/lib/utils";

const LeftBar = () => {
  const [teams, setTeams] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTeam, setEditingTeam] = useState(null);
  const [addingMemberTeam, setAddingMemberTeam] = useState(null);
  const { user } = useAuth();
  const { setSelectedTeam, selectedTeam } = useTeam();
  const token = user?.token;

  useEffect(() => {
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
            setTeams(data || []);
        }
        } catch (error) {
        console.log(error);
        }
    };
    fetchTeams();
  }, [token, isModalOpen, editingTeam]); // Refetch on modal close (after create or edit)

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b border-border flex items-center justify-between bg-muted/20">
        <h2 className="text-xl font-semibold tracking-tight flex items-center gap-2">
            <Users className="w-5 h-5" /> Teams
        </h2>
        <Button
            size="icon"
            variant="ghost"
            onClick={() => setIsModalOpen(true)}
            title="Create new team"
        >
            <Plus className="w-5 h-5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {teams.length === 0 && (
            <div className="text-center text-muted-foreground text-sm py-8">
                No teams found. Create one to get started.
            </div>
        )}
        {teams.map((team) => (
          <Card
            key={team.team_id}
            onClick={() => setSelectedTeam(team)}
            className={cn(
                "cursor-pointer transition-all hover:bg-accent hover:text-accent-foreground border-transparent hover:border-border group relative",
                selectedTeam?.team_id === team.team_id ? "bg-accent border-primary/50 ring-1 ring-primary/20" : "bg-card"
            )}
          >
            <CardHeader className="p-4">
              <div className="flex justify-between items-start">
                  <CardTitle className="text-base font-medium leading-none">
                    {team.team_name}
                  </CardTitle>
                  <div className="absolute right-2 top-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      {user?.user_id === team.team_leader_id && (
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6"
                            title="Add Member"
                            onClick={(e) => {
                                e.stopPropagation();
                                setAddingMemberTeam(team);
                            }}
                        >
                            <UserPlus className="w-3 h-3" />
                        </Button>
                      )}
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6"
                        title="Edit Team"
                        onClick={(e) => {
                            e.stopPropagation();
                            setEditingTeam(team);
                        }}
                      >
                          <Pencil className="w-3 h-3" />
                      </Button>
                  </div>
              </div>
              <CardDescription className="text-xs mt-1.5 truncate">
                Leader: {team.team_leader_name || team.team_leader_id}
              </CardDescription>
            </CardHeader>
          </Card>
        ))}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <CreateTeamForm
            onClose={() => setIsModalOpen(false)}
            onTeamCreated={() => {
              setIsModalOpen(false);
            }}
          />
        </div>
      )}

      {editingTeam && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <EditTeamForm
            team={editingTeam}
            onClose={() => setEditingTeam(null)}
            onTeamUpdated={() => {
              setEditingTeam(null);
            }}
          />
        </div>
      )}

      {addingMemberTeam && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <AddMemberForm
            team={addingMemberTeam}
            onClose={() => setAddingMemberTeam(null)}
          />
        </div>
      )}
    </div>
  );
};

export default LeftBar;