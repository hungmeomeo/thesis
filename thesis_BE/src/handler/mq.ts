import { Operation } from "./../models/operation";

export class MessageQueueClient {
  constructor() {}

  sendMessage(docId: string, operation: Operation): void {
    console.log(`Sending operation for docId ${docId}:`, operation);
  }

  onMessage(docId: string, callback: (message: Operation) => void): void {
    console.log(`Listening for operations on docId ${docId}`);
  }
}
