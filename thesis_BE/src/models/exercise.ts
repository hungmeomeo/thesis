import { Doc } from "./doc";
import { User } from "./user";
import { getDocumentText } from "../db/redis";

export class Exercise {
  exerId: string;
  files: Doc[];
  users: User[];
  lang: string;
  online: boolean;

  constructor(files: Doc[], lang: string, exerId: string) {
    this.exerId = exerId;
    this.files = files;
    this.lang = lang;
    this.online = false;
    this.users = [];
  }

  removeContent() {
    for (let i = 0; i < this.files.length; i++) {
      this.files[i].documentText = "";
    }
  }

  async loadContent() {
    for (let i = 0; i < this.files.length; i++) {
      const text = await getDocumentText(this.files[i].docId);
      console.log("Document text loaded from cache:", text);
      this.files[i].documentText = text !== null ? text : "";
    }
  }

  addUser(user: User) {
    if (this.users.length === 0) {
      this.loadContent();
    }
    this.users.push(user);
  }

  removeUser(userId: string) {
    this.users = this.users.filter((user) => user.userId !== userId);
    if (this.users.length === 0) {
      this.removeContent();
    }
  }
}
