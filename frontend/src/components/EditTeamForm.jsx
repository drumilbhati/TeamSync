import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { useAuth } from "@/context/AuthContext";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "./ui/card";
import { Label } from "./ui/label";
import { Loader2, Trash2 } from "lucide-react";

const EditTeamForm = ({ team, onTeamUpdated, onClose }) => {
  const [teamName, setTeamName] = useState(team.team_name);
  const [loading, setLoading] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const { user } = useAuth();
  const token = user?.token;

  const updateTeam = async (e) => {
    e.preventDefault();
    if (!teamName.trim()) return;

    setLoading(true);
    try {
      const response = await fetch(`/api/teams/${team.team_id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ team_id: team.team_id, team_name: teamName }),
      });

      if (response.ok) {
        if (onTeamUpdated) {
          onTeamUpdated();
        }
        if (onClose) onClose();
      } else {
          console.error("Failed to update team");
      }
    } catch (error) {
      console.log(error);
    } finally {
        setLoading(false);
    }
  };

  const deleteTeam = async () => {
    if (!window.confirm("Are you sure you want to delete this team? This action cannot be undone.")) {
      return;
    }

    setDeleteLoading(true);
    try {
      const response = await fetch(`/api/teams/${team.team_id}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.ok || response.status === 204) {
        if (onTeamUpdated) {
          onTeamUpdated();
        }
        if (onClose) onClose();
      } else {
        console.error("Failed to delete team");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setDeleteLoading(false);
    }
  };

  const isLeader = user?.user_id === team.team_leader_id;

  return (
    <Card className="w-[350px] sm:w-[450px] shadow-lg border-border bg-card text-card-foreground">
      <CardHeader>
        <CardTitle>Edit Team</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={updateTeam} className="grid w-full items-center gap-4">
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="teamName">Team Name</Label>
                <Input
                    id="teamName"
                    type="text"
                    value={teamName}
                    onChange={(e) => setTeamName(e.target.value)}
                    placeholder="Enter team name"
                    disabled={loading || deleteLoading}
                />
            </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between items-center">
          <div className="flex gap-2">
            <Button variant="outline" onClick={onClose} type="button" disabled={loading || deleteLoading}>Cancel</Button>
            {isLeader && (
              <Button 
                variant="destructive" 
                onClick={deleteTeam} 
                disabled={loading || deleteLoading}
                type="button"
              >
                {deleteLoading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Trash2 className="w-4 h-4" />}
              </Button>
            )}
          </div>
          <Button onClick={updateTeam} disabled={loading || deleteLoading || !teamName.trim()}>
              {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Save Changes"}
          </Button>
      </CardFooter>
    </Card>
  );
};

export default EditTeamForm;
