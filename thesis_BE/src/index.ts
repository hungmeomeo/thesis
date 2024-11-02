import express from "express";
import http from "http";
import WebSocket from "ws";
import { User } from "./models/user"; // Import User class
import { Doc } from "./models/doc"; // Import Doc class
import { Operation } from "./models/operation"; // Import Operation class

// Initialize express and WebSocket server
const app = express();
const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

const docs: { [key: string]: Doc } = {};

// WebSocket connection handling
wss.on("connection", (ws: WebSocket, req) => {
  const params = new URLSearchParams(req.url!.split("?")[1]);
  const docId = params.get("docId")!;
  const userId = params.get("userId")!;
  const userName = `User${userId}`; // For demonstration purposes

  const user = new User(userId, userName);

  // Store the document ID in the WebSocket instance
  (ws as any).docId = docId; // Type assertion for TypeScript

  // Create or get the document
  if (!docs[docId]) {
    docs[docId] = new Doc(docId);
  }

  const doc = docs[docId];
  doc.addUser(user);

  console.log(`User ${userId} connected to document ${docId}`);

  ws.send(JSON.stringify({ content: doc.documentText }));

  // Handle incoming messages from clients
  ws.on("message", (message: string) => {
    console.log(`Received message from ${userId}: ${message}`);
    const operationData = JSON.parse(message);

    const operation = new Operation(
      operationData.type,
      operationData.position,
      operationData.text
    );

    const broadcast = JSON.stringify({
      type: operation.type,
      position: operation.position,
      text: operationData.text,
    });

    doc.queueOperation(operation);

    // Send the broadcast message to clients connected to the same document
    wss.clients.forEach((client) => {
      if (
        client.readyState === WebSocket.OPEN &&
        (client as any).docId === docId
      ) {
        client.send(broadcast);
      }
    });
  });

  ws.on("close", () => {
    doc.removeUser(userId);
    console.log(`User ${userId} disconnected from document ${docId}`);
  });
});

// Start the server
const PORT = 8080;
server.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
