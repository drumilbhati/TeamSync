import { useState } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { useAuth } from "@/context/AuthContext";
import { Card, CardHeader, CardTitle, CardContent, CardFooter, CardDescription } from "./ui/card";
import { Label } from "./ui/label";
import { Loader2 } from "lucide-react";

const AddMemberForm = ({ team, onClose }) => {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);
  const { user } = useAuth();
  const token = user?.token;

  const addMember = async (e) => {
    e.preventDefault();
    if (!email.trim()) return;

    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      const response = await fetch("/api/members", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ 
            team_id: parseInt(team.team_id), 
            email: email 
        }),
      });

      if (response.ok) {
        setSuccess(true);
        setEmail("");
        setTimeout(() => {
            if (onClose) onClose();
        }, 1500);
      } else {
          const data = await response.json().catch(() => ({}));
          throw new Error(data.message || "Failed to add member");
      }
    } catch (error) {
      console.error(error);
      setError(error.message);
    } finally {
        setLoading(false);
    }
  };

  return (
    <Card className="w-[350px] sm:w-[450px] shadow-lg border-border bg-card text-card-foreground">
      <CardHeader>
        <CardTitle>Add Team Member</CardTitle>
        <CardDescription>Add a new member to <strong>{team.team_name}</strong></CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={addMember} className="grid w-full items-center gap-4">
            {error && (
                <div className="text-destructive text-sm bg-destructive/10 p-2 rounded">
                    {error}
                </div>
            )}
            {success && (
                <div className="text-green-500 text-sm bg-green-500/10 p-2 rounded">
                    Member added successfully!
                </div>
            )}
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="email">Email Address</Label>
                <Input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Enter email to add"
                    disabled={loading || success}
                />
            </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between">
          <Button variant="outline" onClick={onClose} type="button">Close</Button>
          <Button onClick={addMember} disabled={loading || success || !email.trim()}>
              {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Add Member"}
          </Button>
      </CardFooter>
    </Card>
  );
};

export default AddMemberForm;
