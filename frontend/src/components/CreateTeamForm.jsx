import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { useAuth } from "@/context/AuthContext";
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from "./ui/card";
import { Label } from "./ui/label";

const CreateTeamForm = ({ onTeamCreated, onClose }) => {
  const [teamName, setTeamName] = useState("");
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();
  const token = user?.token;

  const createTeam = async (e) => {
    e.preventDefault();
    if (!teamName.trim()) return;

    setLoading(true);
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
        setTeamName("");
        if (onTeamCreated) {
          onTeamCreated();
        }
      }
    } catch (error) {
      console.log(error);
    } finally {
        setLoading(false);
    }
  };

  return (
    <Card className="w-[350px] sm:w-[450px] shadow-lg border-border bg-card text-card-foreground">
      <CardHeader>
        <CardTitle>Create New Team</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={createTeam} className="grid w-full items-center gap-4">
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="teamName">Team Name</Label>
                <Input
                    id="teamName"
                    type="text"
                    value={teamName}
                    onChange={(e) => setTeamName(e.target.value)}
                    placeholder="Enter team name"
                    disabled={loading}
                />
            </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between">
          <Button variant="outline" onClick={onClose} type="button">Cancel</Button>
          <Button onClick={createTeam} disabled={loading || !teamName.trim()}>
              {loading ? "Creating..." : "Create Team"}
          </Button>
      </CardFooter>
    </Card>
  );
};

export default CreateTeamForm;