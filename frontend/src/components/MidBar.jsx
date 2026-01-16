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
    if (!token) return;

    // Clear messages when connecting to a new team context (or just on change)
    setMessages([]);

    const websocket = new WebSocket(
      `ws://localhost:8080/api/ws?token=${token}`,
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
                    <div
                        className={cn(
                            "max-w-[80%] rounded-2xl px-4 py-2 text-sm shadow-sm",
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
