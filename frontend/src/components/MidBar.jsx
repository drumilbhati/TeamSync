import { useEffect, useState, useRef } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";
import { Send, MessageSquare } from "lucide-react";
import { cn } from "@/lib/utils";

const MidBar = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const wsRef = useRef(null);
  const { user } = useAuth();
  const token = user?.token;
  const { selectedTeam } = useTeam();
  const scrollRef = useRef(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  useEffect(() => {
    if (selectedTeam && token) {
        setMessages([]); // Clear existing messages immediately on switch
        const fetchHistory = async () => {
             try {
                 const res = await fetch(`/api/messages?team_id=${selectedTeam.team_id}`, {
                     headers: { Authorization: `Bearer ${token}` }
                 });
                 if (res.ok) {
                     const data = await res.json();
                     setMessages(data || []);
                 }
             } catch (e) {
                 console.error(e);
             }
        };
        fetchHistory();
    }
  }, [selectedTeam, token]);

  useEffect(() => {
    if (!token) return;

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const websocket = new WebSocket(
      `${protocol}//${window.location.host}/api/ws?token=${token}`,
    );
    wsRef.current = websocket;

    websocket.onopen = () => console.log("Connected to Websocket server");
    websocket.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (selectedTeam && msg.team_id === selectedTeam.team_id) {
            setMessages((prevMsg) => [...prevMsg, msg]);
        }
      } catch (error) {
        console.error("Failed to parse message", error);
        setMessages((prevMsg) => [...prevMsg, { content: event.data, isSystem: true }]);
      }
    };
    websocket.onclose = () => console.log("Disconnected from Websocket server");

    return () => {
        websocket.close();
        wsRef.current = null;
    };
  }, [token, selectedTeam]);

  const sendMessage = (e) => {
    e?.preventDefault();
    if (
      wsRef.current &&
      wsRef.current.readyState == WebSocket.OPEN &&
      input.trim() !== "" &&
      selectedTeam
    ) {
      const msg = {
        team_id: selectedTeam.team_id,
        content: input,
        user_id: user?.user_id 
      };
      wsRef.current.send(JSON.stringify(msg));
      setInput("");
    }
  };

  if (!selectedTeam) {
      return (
          <div className="flex h-full items-center justify-center p-8 text-muted-foreground text-center">
              <div>
                  <MessageSquare className="w-12 h-12 mx-auto mb-4 opacity-20" />
                  <p>Select a team to start chatting</p>
              </div>
          </div>
      );
  }

  return (
    <div className="flex flex-col h-full bg-background/50">
      {/* Header */}
      <div className="p-4 border-b border-border bg-muted/20 backdrop-blur-sm sticky top-0 z-10">
        <h2 className="text-xl font-semibold tracking-tight flex items-center gap-2">
            <MessageSquare className="w-5 h-5" />
            {selectedTeam.team_name}
        </h2>
      </div>

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4" ref={scrollRef}>
        {messages.map((message, index) => {
            const isMe = message.user_id === user?.user_id;
            const isSystem = message.isSystem;
            
            return (
                <div
                    key={index}
                    className={cn(
                        "flex w-full",
                        isMe ? "justify-end" : "justify-start"
                    )}
                >
                    <div className={cn("flex flex-col max-w-[80%]", isMe ? "items-end" : "items-start")}>
                        {!isMe && !isSystem && (
                            <span className="text-[10px] text-muted-foreground mb-1 px-1">
                                {message.user_name}
                            </span>
                        )}
                        <div
                            className={cn(
                                "rounded-2xl px-4 py-2 text-sm shadow-sm break-words",
                                isMe 
                                    ? "bg-primary text-primary-foreground rounded-br-none" 
                                    : isSystem 
                                        ? "bg-muted text-muted-foreground text-xs mx-auto" 
                                        : "bg-muted text-foreground rounded-bl-none"
                            )}
                        >
                            <p>{message.content}</p>
                        </div>
                    </div>
                </div>
            );
        })}
      </div>

      {/* Input Area */}
      <div className="p-4 border-t border-border bg-background mt-auto">
        <form onSubmit={sendMessage} className="flex gap-2">
            <Input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder={`Message ${selectedTeam.team_name}...`}
                className="flex-1 bg-muted/50 focus-visible:ring-primary"
            />
            <Button type="submit" size="icon" disabled={!input.trim()}>
                <Send className="w-4 h-4" />
                <span className="sr-only">Send</span>
            </Button>
        </form>
      </div>
    </div>
  );
};

export default MidBar;
