import WebSocket from "ws";
import { User } from "../models/user";
import { Doc } from "../models/doc";
import { Operation } from "../models/operation";
import { wss } from "../server";

const docs: { [key: string]: Doc } = {};

export function websocketHandler(ws: WebSocket, req: any): void {
  const params = new URLSearchParams(req.url!.split("?")[1]);
  const docId = params.get("docId")!;
  const userId = params.get("userId")!;
  const exerId = params.get("exerId")!;
  if (docId && userId) {
    const userName = `User${userId}`;

    const user = new User(userId, userName);

    (ws as any).docId = docId; // Type assertion

    if (!docs[docId]) {
      docs[docId] = new Doc(docId);
    }

    const doc = docs[docId];
    (async () => {
      await doc.addUser(user);

      console.log(`User ${userId} connected to document ${docId}`);

      wss.clients.forEach((client) => {
        if (
          client.readyState === WebSocket.OPEN &&
          (client as any).docId === docId
        ) {
          client.send(doc.users ? JSON.stringify(doc.users) : "[]");
        }
      });

      ws.send(JSON.stringify({ content: doc.documentText }));
    })();

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
      wss.clients.forEach((client) => {
        if (
          client.readyState === WebSocket.OPEN &&
          (client as any).docId === docId
        ) {
          client.send(doc.users ? JSON.stringify(doc.users) : "[]");
        }
      });
    });
  }
  if (exerId && userId) {
    console.log(`User ${userId} connected to exercise ${exerId}`);
    ws.on("close", () => {
      console.log(`User ${userId} disconnected from exercise ${exerId}`);
    });
  }
}
