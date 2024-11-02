export enum OperationType {
  INSERT = "insert",
  DELETE = "delete",
}

export class Operation {
  type: OperationType;
  position: number;
  text: string;
  timestamp: number;

  constructor(type: OperationType, position: number, text: string) {
    this.type = type;
    this.position = position;
    this.text = text;
    this.timestamp = Date.now();
  }
}
