import express from "express";
import http from "http";
import WebSocket from "ws";
import { websocketHandler } from "./routes/ws";
import { connectRedis } from "./db/redis";

export const app = express();
export const server = http.createServer(app);
export const wss = new WebSocket.Server({ server });

//connectRedis();
wss.on("connection", websocketHandler);
