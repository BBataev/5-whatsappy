let socket;

window.onload = async () => {
  socket = new WebSocket("ws://localhost:8080/api/ws");

  socket.onopen = () => {
    console.log("Connected to WebSocket");
  };

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    appendMessage(`From ${msg.sender_id}: ${msg.message}`);
  };

  socket.onclose = () => {
    console.log("Disconnected from WebSocket");
  };

  document.getElementById("sendButton").onclick = () => {
    const recipientId = document.getElementById("recipientId").value.trim();
    const message = document.getElementById("messageInput").value.trim();

    if (!recipientId || !message) return;

    const msg = {
      recipient_id: parseInt(recipientId),
      message: message,
    };

    socket.send(JSON.stringify(msg));
    appendMessage(`To ${recipientId}: ${message}`);
    document.getElementById("messageInput").value = "";
  };
};

function appendMessage(text) {
  const div = document.getElementById("messages");
  const p = document.createElement("p");
  p.textContent = text;
  div.appendChild(p);
}
