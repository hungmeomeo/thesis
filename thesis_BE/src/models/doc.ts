import { Operation, OperationType } from "./operation";
import { User } from "./user";

export class Doc {
  docId: string;
  documentText: string;
  operationHistory: Operation[];
  users: User[];
  operationQueue: Operation[]; // Queue to store operations

  constructor(docId: string) {
    this.docId = docId;
    this.documentText = "";
    this.operationHistory = [];
    this.users = [];
    this.operationQueue = [];
  }

  addUser(user: User) {
    this.users.push(user);
  }

  removeUser(userId: string) {
    this.users = this.users.filter((user) => user.userId !== userId);
  }

  // Add an operation to the queue
  queueOperation(operation: Operation): void {
    this.operationQueue.push(operation);
    this.processQueue();
  }

  // Process the operation queue
  private processQueue(): void {
    while (this.operationQueue.length > 0) {
      const operation = this.operationQueue.shift()!;

      // Transform operation if necessary
      const transformedOperation = this.transformOperation(operation);

      // Apply the transformed operation to the document
      this.applyOperation(transformedOperation);

      // Store the operation in the history
      this.operationHistory.push(transformedOperation);
    }
  }

  // Apply an individual operation (already transformed)
  private applyOperation(operation: Operation): void {
    const { type, position, text } = operation;

    switch (type) {
      case OperationType.INSERT:
        this.documentText =
          this.documentText.slice(0, position) +
          text +
          this.documentText.slice(position);
        break;
      case OperationType.DELETE:
        this.documentText =
          this.documentText.slice(0, position) +
          this.documentText.slice(position + text.length);
        break;
      default:
        throw new Error("Unknown operation type");
    }
  }

  private transformOperation(newOp: Operation): Operation {
    // For each operation in the history, transform the new operation to avoid conflicts
    for (let i = 0; i < this.operationHistory.length; i++) {
      const historyOp = this.operationHistory[i];

      if (
        historyOp.type === OperationType.INSERT &&
        newOp.type === OperationType.INSERT
      ) {
        if (historyOp.position <= newOp.position) {
          newOp.position += historyOp.text.length;
        }
      } else if (
        historyOp.type === OperationType.DELETE &&
        newOp.type === OperationType.INSERT
      ) {
        if (historyOp.position < newOp.position) {
          newOp.position -= historyOp.text.length;
        }
      } else if (
        historyOp.type === OperationType.INSERT &&
        newOp.type === OperationType.DELETE
      ) {
        if (historyOp.position <= newOp.position) {
          newOp.position += historyOp.text.length;
        }
      } else if (
        historyOp.type === OperationType.DELETE &&
        newOp.type === OperationType.DELETE
      ) {
        if (historyOp.position < newOp.position) {
          newOp.position -= historyOp.text.length;
        }
      }
    }

    return newOp;
  }

  // Get connected users
  getConnectedUsers() {
    return [...this.users];
  }
}
