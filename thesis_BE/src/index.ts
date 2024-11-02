import { app, server } from "./server";

const PORT = 8080;

server.listen(PORT, () => {
  console.log(`Server is running on`);
});
