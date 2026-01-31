import { useEffect, useState, useRef } from "react";
import { useAuth } from "@/context/AuthContext";
import { Card, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { X, Send, User, Clock, Calendar, Sparkles, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { enhanceTask } from "@/lib/api";

const statusMap = {
    todo: "To Do",
    in_review: "In Review",
    in_progress: "In Progress",
    done: "Done",
};

const TaskDetailsModal = ({ task: initialTask, onClose }) => {
    const { user } = useAuth();
    const token = user?.token;
    const [task, setTask] = useState(initialTask);
    const [comments, setComments] = useState([]);
    const [newComment, setNewComment] = useState("");
    const [submitting, setSubmitting] = useState(false);
    const [enhancing, setEnhancing] = useState(false);
    const commentsEndRef = useRef(null);

    // Fetch comments
    const fetchComments = async () => {
         try {
            const response = await fetch(`/api/comments/${task.task_id}`, {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (response.ok) {
                const data = await response.json();
                setComments(data || []);
            }
         } catch (e) {
             console.error("Failed to fetch comments", e);
         }
    };

    useEffect(() => {
        if (task?.task_id) {
            fetchComments();
        }
    }, [task?.task_id, token]);

    const scrollToBottom = () => {
        commentsEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }

    useEffect(() => {
        scrollToBottom();
    }, [comments]);

    const handleSendComment = async (e) => {
        e.preventDefault();
        if (!newComment.trim()) return;

        setSubmitting(true);
        try {
            const response = await fetch("/api/comments", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({
                    task_id: task.task_id,
                    user_id: user.user_id,
                    content: newComment
                })
            });

            if (response.ok) {
                const createdComment = await response.json();
                
                // Optimistically update
                setComments(prev => [...prev, createdComment]);
                setNewComment("");
            }
        } catch (e) {
            console.error(e);
        } finally {
            setSubmitting(false);
        }
    };

    const handleEnhance = async () => {
        setEnhancing(true);
        try {
            const response = await enhanceTask(task.task_id, token);
            if (response.ok) {
                const enhancedTask = await response.json();
                setTask(enhancedTask);
            } else {
                console.error("Failed to enhance task");
            }
        } catch (e) {
            console.error("Error enhancing task:", e);
        } finally {
            setEnhancing(false);
        }
    }

    const getPriorityColor = (priority) => {
        switch(priority) {
            case 'high': return "text-red-500 bg-red-500/10 border-red-500/20";
            case 'medium': return "text-yellow-500 bg-yellow-500/10 border-yellow-500/20";
            case 'low': return "text-green-500 bg-green-500/10 border-green-500/20";
            default: return "text-slate-500 bg-slate-500/10 border-slate-500/20";
        }
    }

    if (!task) return null;

    const isCreator = user?.user_id === task.creator_id;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4 animate-in fade-in zoom-in-95 duration-200">
            <Card className="w-full max-w-4xl h-[85vh] flex flex-col shadow-2xl bg-card text-card-foreground border-border overflow-hidden">
                <CardHeader className="flex flex-row items-start justify-between p-6 border-b shrink-0 bg-muted/5">
                    <div className="space-y-1 pr-8 flex-1">
                         <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <CardTitle className="text-2xl font-bold tracking-tight">{task.title}</CardTitle>
                                <span className={cn("text-[10px] uppercase font-bold px-2 py-0.5 rounded border", getPriorityColor(task.priority))}>
                                    {task.priority}
                                </span>
                            </div>
                            {isCreator && (
                                <Button 
                                    size="sm" 
                                    variant="outline" 
                                    className="gap-2 text-xs h-8 border-primary/20 hover:bg-primary/10 hover:text-primary transition-all duration-300 shadow-sm"
                                    onClick={handleEnhance}
                                    disabled={enhancing}
                                >
                                    {enhancing ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Sparkles className="w-3.5 h-3.5" />}
                                    {enhancing ? "Enhancing..." : "Enhance with AI"}
                                </Button>
                            )}
                         </div>
                         <div className="flex items-center gap-4 text-xs text-muted-foreground mt-2">
                             <div className="flex items-center gap-1">
                                <Clock className="w-3.5 h-3.5" />
                                <span>{statusMap[task.status] || task.status}</span>
                             </div>
                             {task.created_at && (
                                 <div className="flex items-center gap-1">
                                    <Calendar className="w-3.5 h-3.5" />
                                    <span>Created {new Date(task.created_at).toLocaleDateString()}</span>
                                 </div>
                             )}
                         </div>
                    </div>
                    <Button variant="ghost" size="icon" onClick={onClose} className="shrink-0 -mr-2 -mt-2 ml-4">
                        <X className="w-5 h-5" />
                    </Button>
                </CardHeader>
                
                <div className="flex-1 overflow-hidden flex flex-col md:flex-row">
                    {/* Left Side: Task Details */}
                    <div className="p-6 md:w-5/12 border-b md:border-b-0 md:border-r overflow-y-auto bg-muted/5 space-y-6">
                         <div className="space-y-2">
                            <h4 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Description</h4>
                            <div className="text-sm leading-relaxed whitespace-pre-wrap text-foreground/90">
                                {task.description?.String || <span className="text-muted-foreground italic">No description provided.</span>}
                            </div>
                         </div>
                         
                         <div className="space-y-2">
                            <h4 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">Assignee</h4>
                            <div className="flex items-center gap-2">
                                <div className="w-8 h-8 rounded-full bg-secondary flex items-center justify-center">
                                    <User className="w-4 h-4 text-secondary-foreground" />
                                </div>
                                <span className="text-sm font-medium">
                                    {task.assignee_name || (task.assignee_id?.Valid ? `User ID: ${task.assignee_id.Int64}` : "Unassigned")}
                                </span>
                            </div>
                         </div>

                         {/* Can add more details here like Due Date if available */}
                    </div>

                    {/* Right Side: Comments */}
                    <div className="flex-1 flex flex-col min-h-0 bg-background">
                         <div className="p-4 border-b bg-muted/10 flex items-center justify-between">
                             <h4 className="text-sm font-semibold flex items-center gap-2">
                                 Comments
                                 <span className="bg-primary/10 text-primary text-[10px] px-2 py-0.5 rounded-full">
                                     {comments.length}
                                 </span>
                             </h4>
                         </div>
                         
                         <div className="flex-1 overflow-y-auto p-4 space-y-6">
                             {comments.length === 0 && (
                                 <div className="h-full flex flex-col items-center justify-center text-muted-foreground space-y-2">
                                     <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center">
                                         <User className="w-6 h-6 opacity-20" />
                                     </div>
                                     <p className="text-sm">No comments yet.</p>
                                     <p className="text-xs opacity-50">Be the first to share your thoughts.</p>
                                 </div>
                             )}
                             
                             {comments.map(comment => (
                                 <div key={comment.comment_id} className="flex gap-3 group">
                                     <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0 mt-0.5 border border-primary/20">
                                         <span className="text-xs font-bold text-primary">
                                             {(comment.user_name || "U")[0].toUpperCase()}
                                         </span>
                                     </div>
                                     <div className="space-y-1.5 flex-1">
                                         <div className="flex items-center justify-between">
                                             <div className="flex items-center gap-2">
                                                 <span className="text-sm font-semibold">{comment.user_name || `User ${comment.user_id}`}</span>
                                                 {comment.user_id === user?.user_id && (
                                                     <span className="text-[10px] bg-primary/10 text-primary px-1.5 py-0.5 rounded">You</span>
                                                 )}
                                             </div>
                                             <span className="text-[10px] text-muted-foreground tabular-nums">
                                                 {new Date(comment.created_at).toLocaleString(undefined, {
                                                     month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit'
                                                 })}
                                             </span>
                                         </div>
                                         <div className="text-sm text-foreground bg-muted/30 p-3 rounded-lg rounded-tl-none border border-border/50 group-hover:border-border/80 transition-colors">
                                             {comment.content}
                                         </div>
                                     </div>
                                 </div>
                             ))}
                             <div ref={commentsEndRef} />
                         </div>

                         <div className="p-4 border-t bg-background mt-auto">
                             <form onSubmit={handleSendComment} className="relative">
                                 <textarea 
                                    className="flex min-h-[80px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 resize-none pr-12"
                                    placeholder="Write a comment..."
                                    value={newComment}
                                    onChange={e => setNewComment(e.target.value)}
                                    onKeyDown={e => {
                                        if (e.key === 'Enter' && !e.shiftKey) {
                                            e.preventDefault();
                                            handleSendComment(e);
                                        }
                                    }}
                                 />
                                 <Button 
                                    size="icon" 
                                    type="submit" 
                                    disabled={!newComment.trim() || submitting}
                                    className="absolute right-2 bottom-2 h-8 w-8"
                                 >
                                     <Send className="w-4 h-4" />
                                 </Button>
                             </form>
                             <p className="text-[10px] text-muted-foreground mt-2 text-right">
                                 Press <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">Enter</kbd> to send, <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">Shift + Enter</kbd> for new line
                             </p>
                         </div>
                    </div>
                </div>
            </Card>
        </div>
    );
};

export default TaskDetailsModal;
