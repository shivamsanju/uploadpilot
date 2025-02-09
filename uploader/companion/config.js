import dotenv from "dotenv";
import fs from "fs";

// Load environment variables from .env file
dotenv.config({ path: "../.env" });

export const mkdirp = (dir) => {
  try {
    fs.mkdirSync(dir);
  } catch (err) {
    if (err.code !== "EEXIST") {
      throw err;
    }
  }
};

export const getEnv = (varName, defaultValue) => {
  return process.env[varName] || defaultValue;
};
