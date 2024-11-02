import { createClient } from "redis";

const REDIS_URL = "your_redis_cloud_url"; // Replace with your Redis Cloud URL

const redisClient = createClient({
  url: REDIS_URL,
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
