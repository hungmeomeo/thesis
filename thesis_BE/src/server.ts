import express from "express";
import http from "http";
import WebSocket from "ws";
import { websocketHandler } from "./routes/ws";

export const app = express();
export const server = http.createServer(app);
export const wss = new WebSocket.Server({ server });

// Handle WebSocket connections
wss.on("connection", websocketHandler);
