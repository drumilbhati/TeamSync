import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";
import { useEffect, useState } from "react";
import { Card, CardDescription, CardHeader } from "./ui/card";

const RightBar = () => {
  const { selectedTeam } = useTeam();
  const [tasks, setTasks] = useState([]);
  const { user } = useAuth();
  const token = user?.token;

  useEffect(() => {
    const fetchTasks = async () => {
      if (!selectedTeam) return;
      try {
        const query = `team_id=${selectedTeam.team_id}`;
        const response = await fetch(`/api/tasks?${query}`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          setTasks(data || []);
        }
      } catch (error) {
        console.log(error);
      }
    };
    fetchTasks();
  }, [selectedTeam, token]);

  if (!selectedTeam) {
    return <div className="p-4">Please select a team</div>;
  }

  return (
    <div className="flex flex-col gap-4">
      <Card key={selectedTeam.team_id}>
        <div>
          <CardHeader>
            <h3 className="text-teal-950!">
              Selected Team:{selectedTeam.team_name}
            </h3>
          </CardHeader>
        </div>
      </Card>
      {tasks && tasks.length > 0 ? (
        tasks.map((task) => (
          <Card key={task.task_id}>
            <div>
              <CardHeader>
                <h3 className="text-teal-950!">{task.title}</h3>
              </CardHeader>
              <CardDescription>
                <p className="text-teal-950!">{task.description?.String}</p>
              </CardDescription>
            </div>
          </Card>
        ))
      ) : (
        <p>No tasks found.</p>
      )}
    </div>
  );
};

export default RightBar;
