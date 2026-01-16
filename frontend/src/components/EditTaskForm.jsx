import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";
import { useState, useEffect } from "react";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardFooter,
} from "./ui/card";
import { Loader2 } from "lucide-react";

const EditTaskForm = ({ task, onClose, onTaskUpdated }) => {
  const { selectedTeam } = useTeam();
  const { user } = useAuth();
  const token = user?.token;

  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description?.String || "");
  const [status, setStatus] = useState(task.status);
  const [priority, setPriority] = useState(task.priority);
  const [assignedTo, setAssignedTo] = useState(task.assigned_to || "");
  const [loading, setLoading] = useState(false);
  const [members, setMembers] = useState([]);

  // Fetch team members to populate the assignment dropdown
  useEffect(() => {
      const fetchMembers = async () => {
          if (!selectedTeam) return;
          try {
              // Assuming this endpoint exists based on the requirement
              const response = await fetch(`/api/teams/members?team_id=${selectedTeam.team_id}`, {
                  headers: {
                      Authorization: `Bearer ${token}`
                  }
              });
              if (response.ok) {
                  const data = await response.json();
                  setMembers(data || []);
              }
          } catch (error) {
              console.error("Failed to fetch members", error);
          }
      };
      fetchMembers();
  }, [selectedTeam, token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title || !selectedTeam) return;

    setLoading(true);
    try {
      const response = await fetch("/api/tasks", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          task_id: task.task_id,
          title,
          description: { String: description, Valid: true },
          status,
          priority,
          team_id: parseInt(selectedTeam.team_id),
          assigned_to: assignedTo ? parseInt(assignedTo) : null
        }),
      });

      if (response.ok) {
        if (onTaskUpdated) onTaskUpdated();
        if (onClose) onClose();
      } else {
        console.error("Failed to update task");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-[350px] sm:w-[450px] shadow-lg bg-card text-card-foreground border-border">
      <CardHeader>
        <CardTitle>Edit Task</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={handleSubmit}
          className="grid w-full items-center gap-4"
        >
          <div className="flex flex-col space-y-1.5">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              placeholder="Task title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              className="bg-background"
            />
          </div>
          <div className="flex flex-col space-y-1.5">
            <Label htmlFor="description">Description</Label>
            <Input
              id="description"
              placeholder="Task description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="bg-background"
            />
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="status">Status</Label>
                <select
                id="status"
                className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                >
                <option value="todo" className="bg-popover text-popover-foreground">To Do</option>
                <option value="in_progress" className="bg-popover text-popover-foreground">In Progress</option>
                <option value="in_review" className="bg-popover text-popover-foreground">In Review</option>
                <option value="done" className="bg-popover text-popover-foreground">Done</option>
                </select>
            </div>
            <div className="flex flex-col space-y-1.5">
                <Label htmlFor="priority">Priority</Label>
                <select
                id="priority"
                className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                value={priority}
                onChange={(e) => setPriority(e.target.value)}
                >
                <option value="low" className="bg-popover text-popover-foreground">Low</option>
                <option value="medium" className="bg-popover text-popover-foreground">Medium</option>
                <option value="high" className="bg-popover text-popover-foreground">High</option>
                </select>
            </div>
          </div>

          <div className="flex flex-col space-y-1.5">
            <Label htmlFor="assignedTo">Assign To</Label>
            <select
              id="assignedTo"
              className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              value={assignedTo}
              onChange={(e) => setAssignedTo(e.target.value)}
            >
              <option value="" className="bg-popover text-popover-foreground">Unassigned</option>
              {members.map(member => (
                  <option key={member.user_id} value={member.user_id} className="bg-popover text-popover-foreground">
                      {member.username || member.email}
                  </option>
              ))}
            </select>
          </div>

        </form>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button
          variant="outline"
          onClick={onClose}
          type="button"
        >
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Save Changes"}
        </Button>
      </CardFooter>
    </Card>
  );
};

export default EditTaskForm;
