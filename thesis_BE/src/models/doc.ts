import { Operation, OperationType } from "./operation";
import { User } from "./user";
import { getDocumentText, setDocumentText } from "../db/redis";

export class Doc {
  docName: string;
  docId: string;
  documentText: string | null; // Allow documentText to be null
  users: User[];
  operationQueue: Operation[];

  constructor(docId: string) {
    this.docName = "";
    this.docId = docId;
    this.documentText = null; // Initialize as null
    this.users = [];
    this.operationQueue = [];
  }

  async addUser(user: User) {
    this.users.push(user);
  }

  removeUser(userId: string) {
    this.users = this.users.filter((user) => user.userId !== userId);
    if (this.users.length === 0) {
      setDocumentText(this.docId, this.documentText ?? "");
    }
  }

  queueOperation(operation: Operation): void {
    this.operationQueue.push(operation);
    this.processQueue();
  }

  private processQueue(): void {
    const operationsToProcess = [...this.operationQueue];
    this.operationQueue = []; // Clear the queue for new operations

    operationsToProcess.forEach((operation) => {
      const transformedOperation = this.transformOperation(operation);
      this.applyOperation(transformedOperation);
    });
  }

  private applyOperation(operation: Operation): void {
    const { type, position, text } = operation;

    switch (type) {
      case OperationType.INSERT:
        this.documentText =
          (this.documentText ?? "").slice(0, position) +
          text +
          (this.documentText ?? "").slice(position);
        break;
      case OperationType.DELETE:
        this.documentText =
          (this.documentText ?? "").slice(0, position) +
          (this.documentText ?? "").slice(position + text.length);
        break;
      default:
        throw new Error("Unknown operation type");
    }

    console.log(this.operationQueue);
  }

  private transformOperation(newOp: Operation): Operation {
    return newOp;
  }

  clearQueue(): void {
    this.operationQueue = [];
  }

  getConnectedUsers(): User[] {
    return [...this.users];
  }
}
