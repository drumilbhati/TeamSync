import { useEffect, useState } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useAuth } from "@/context/AuthContext";
import { useTeam } from "@/context/TeamContext";

const MidBar = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [ws, setWs] = useState(null);
  const { user } = useAuth();
  const token = user?.token;
  const { selectedTeam } = useTeam();

  useEffect(() => {
    if (!token) return;

    const websocket = new WebSocket(
      `ws://localhost:8080/api/ws?token=${token}`,
    );
    setWs(websocket);

    websocket.onopen = () => console.log("Connected to Websocket server");
    websocket.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (selectedTeam && msg.team_id === selectedTeam.team_id) {
          setMessages((prevMsg) => [...prevMsg, msg.content]);
        }
      } catch (error) {
        console.error("Failed to parse message", error);
        setMessages((prevMsg) => [...prevMsg, event.data]);
      }
    };
    websocket.onclose = () => console.log("Disconnected from Websocket server");

    return () => websocket.close();
  }, [token, selectedTeam]);

  const sendMessage = () => {
    if (
      ws &&
      ws.readyState == WebSocket.OPEN &&
      input.trim() !== "" &&
      selectedTeam
    ) {
      const msg = {
        team_id: selectedTeam.team_id,
        content: input,
      };
      ws.send(JSON.stringify(msg));
      setInput("");
    }
  };

  return (
    <div>
      <h2>
        {selectedTeam ? `Chat for: ${selectedTeam.team_name}` : "select a team"}
      </h2>
      <div>
        {messages.map((message, index) => (
          <p key={index}>{message}</p>
        ))}
      </div>
      <Input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Type a message"
      ></Input>
      <Button onClick={sendMessage}>Send</Button>
    </div>
  );
};

export default MidBar;
