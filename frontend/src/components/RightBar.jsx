import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";
import { useEffect, useState, useCallback } from "react";
import { Card, CardDescription, CardHeader, CardTitle, CardContent } from "./ui/card";
import { Button } from "./ui/button";
import CreateTaskForm from "./CreateTaskForm";
import EditTaskForm from "./EditTaskForm";
import TaskDetailsModal from "./TaskDetailsModal";
import { Plus, CheckCircle2, Clock, Circle, AlertCircle, Pencil } from "lucide-react";
import { cn } from "@/lib/utils";

const statusMap = {
    todo: "To Do",
    in_review: "In Review",
    in_progress: "In Progress",
    done: "Done",
};

const priorityOrder = {
    high: 3,
    medium: 2,
    low: 1,
};

const RightBar = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState(null);
  const [viewingTask, setViewingTask] = useState(null);
  const { selectedTeam } = useTeam();
  const [tasks, setTasks] = useState([]);
  const { user } = useAuth();
  const token = user?.token;

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
      }
    } catch (error) {
      console.log(error);
    }
  }, [selectedTeam, token]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks, editingTask]);

  if (!selectedTeam) {
    return (
        <div className="flex h-full items-center justify-center p-8 text-muted-foreground text-center">
            <div>
                <Circle className="w-12 h-12 mx-auto mb-4 opacity-20" />
                <p>Select a team to view tasks</p>
            </div>
        </div>
    );
  }

  const getStatusIcon = (status) => {
      switch(status) {
          case 'done': return <CheckCircle2 className="w-4 h-4 text-green-500" />;
          case 'in_progress': return <Clock className="w-4 h-4 text-blue-500" />;
          case 'in_review': return <AlertCircle className="w-4 h-4 text-yellow-500" />;
          default: return <Circle className="w-4 h-4 text-slate-500" />;
      }
  };

  const getPriorityColor = (priority) => {
      switch(priority) {
          case 'high': return "text-red-500 bg-red-500/10 border-red-500/20";
          case 'medium': return "text-yellow-500 bg-yellow-500/10 border-yellow-500/20";
          case 'low': return "text-green-500 bg-green-500/10 border-green-500/20";
          default: return "text-slate-500 bg-slate-500/10 border-slate-500/20";
      }
  }

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b border-border flex items-center justify-between bg-muted/20">
         <div>
            <h2 className="text-xl font-semibold tracking-tight">Tasks</h2>
            <p className="text-xs text-muted-foreground">{selectedTeam.team_name}</p>
         </div>
         <Button onClick={() => setIsModalOpen(true)} size="sm" className="gap-2">
            <Plus className="w-4 h-4" /> New Task
         </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {tasks.length === 0 && (
             <div className="text-center text-muted-foreground text-sm py-12">
                No tasks found for this team.
            </div>
        )}
        
        {tasks.map((task) => (
          <Card 
            key={task.task_id} 
            className="group hover:border-primary/50 transition-colors relative cursor-pointer"
            onClick={() => setViewingTask(task)}
          >
            <CardHeader className="p-4 pb-2">
              <div className="flex justify-between items-start gap-2">
                  <div className="pr-6">
                      <CardTitle className="text-base font-medium leading-none">
                        {task.title}
                      </CardTitle>
                      {task.assigned_to && (
                          <p className="text-[10px] text-muted-foreground mt-1">
                              Assigned to ID: {task.assigned_to}
                          </p>
                      )}
                  </div>
                  <span className={cn("text-[10px] uppercase font-bold px-2 py-0.5 rounded border", getPriorityColor(task.priority))}>
                    {task.priority}
                  </span>
                  
                  <Button 
                    variant="ghost" 
                    size="icon" 
                    className="absolute right-2 top-2 h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                    onClick={(e) => {
                        e.stopPropagation();
                        setEditingTask(task);
                    }}
                  >
                      <Pencil className="w-3 h-3" />
                  </Button>
              </div>
            </CardHeader>
            <CardContent className="p-4 pt-2">
              <CardDescription className="mb-3 line-clamp-2">
                {task.description?.String || "No description"}
              </CardDescription>
              
              <div className="flex items-center gap-2 text-xs text-muted-foreground bg-muted/30 p-2 rounded-md w-fit">
                {getStatusIcon(task.status)}
                <span className="font-medium uppercase tracking-wider">{statusMap[task.status] || task.status}</span>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <CreateTaskForm 
            onClose={() => setIsModalOpen(false)} 
            onTaskCreated={fetchTasks}
          />
        </div>
      )}

      {editingTask && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
          <EditTaskForm 
            task={editingTask}
            onClose={() => setEditingTask(null)} 
            onTaskUpdated={() => {
                setEditingTask(null);
                fetchTasks(); // Ensure refetch happens
            }}
          />
        </div>
      )}

      {viewingTask && (
        <TaskDetailsModal 
            task={viewingTask}
            onClose={() => setViewingTask(null)}
        />
      )}
    </div>
  );
};

export default RightBar;
