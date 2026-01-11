import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";
import { useEffect, useState, useCallback } from "react";
import { Card, CardDescription, CardHeader } from "./ui/card";
import { Button } from "./ui/button";
import CreateTaskForm from "./CreateTaskForm";

const RightBar = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { selectedTeam } = useTeam();
  const [tasks, setTasks] = useState([]);
  const { user } = useAuth();
  const token = user?.token;

  const statusMap = {
    todo: "TODO",
    in_review: "IN REVIEW",
    in_progress: "IN PROGRESS",
    done: "DONE",
  };

  const statusColourMap = {
    todo: "text-red-500!",
    in_review: "text-yellow-700!",
    in_progress: "text-blue-500!",
    done: "text-green-500!",
  };

  const priorityOrder = {
    high: 3,
    medium: 2,
    low: 1,
  };

  const priorityMap = {
    high: "HIGH",
    medium: "MEDIUM",
    low: "LOW",
  };
  const priorityColourMap = {
    high: "text-red-500!",
    medium: "text-yellow-500!",
    low: "text-green-500!",
  };

  const fetchTasks = useCallback(async () => {
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
        const sortedData = (data || []).sort((a, b) => {
          const priorityA = priorityOrder[a.priority] || 0;
          const priorityB = priorityOrder[b.priority] || 0;
          return priorityB - priorityA;
        });
        setTasks(sortedData);
        console.table(sortedData);
      }
    } catch (error) {
      console.log(error);
    }
  }, [selectedTeam, token]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

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

      <Button
        onClick={() => {
          setIsModalOpen(true);
        }}
      >
        Create new task
      </Button>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
          <CreateTaskForm 
            onClose={() => setIsModalOpen(false)} 
            onTaskCreated={fetchTasks}
          />
        </div>
      )}

      {tasks && tasks.length > 0 ? (
        tasks.map((task) => (
          <Card key={task.task_id} className="bg-purple-200">
            <div>
              <CardHeader>
                <h3 className="text-teal-950!">{task.title}</h3>
              </CardHeader>
              <CardDescription>
                <p className="text-teal-950!">{task.description?.String}</p>
                <p className={statusColourMap[task.status] || "text-teal-950!"}>
                  {statusMap[task.status]}
                </p>
                <p
                  className={
                    priorityColourMap[task.priority] || "text-teal-950!"
                  }
                >
                  {priorityMap[task.priority]}
                </p>
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