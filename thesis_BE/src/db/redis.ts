import { createClient } from "redis";

const redisClient = createClient({
  password: "T9j3YkngwEWotDSm6HQlm33kvWbXCBoW",
  socket: {
    host: "redis-15725.c325.us-east-1-4.ec2.redns.redis-cloud.com",
    port: 15725,
  },
});

redisClient.on("error", (err: string) => {
  console.error("Redis Client Error", err);
});

const connectRedis = async () => {
  try {
    await redisClient.connect();
    console.log("Connected to Redis Cloud");
  } catch (error) {
    console.error("Could not connect to Redis Cloud", error);
  }
};

const setDocumentText = async (docId: string, text: string) => {
  try {
    await redisClient.set(docId, text);
    console.log(`Document text for ${docId} set successfully.`);
  } catch (error) {
    console.error("Error setting document text:", error);
  }
};

const getDocumentText = async (docId: string): Promise<string | null> => {
  try {
    const text = await redisClient.get(docId);
    return text;
  } catch (error) {
    console.error("Error getting document text:", error);
    return null;
  }
};

const deleteDocumentText = async (docId: string) => {
  try {
    await redisClient.del(docId);
    console.log(`Document text for ${docId} deleted successfully.`);
  } catch (error) {
    console.error("Error deleting document text:", error);
  }
};

export {
  connectRedis,
  setDocumentText,
  getDocumentText,
  deleteDocumentText,
  redisClient,
};
