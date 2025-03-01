import express from "express";
import bodyParser from "body-parser";
import session from "express-session";
import companion from "@uppy/companion";
import { mkdirp, getEnv } from "./config.js";
import dotenv from "dotenv";
import fs from "fs";

// Load environment variables from .env file
dotenv.config({ path: "../.env" });

const app = express();

// Companion requires body-parser and express-session middleware.
// You can add it like this if you use those throughout your app.
// If you are using something else in your app, you can add these
// middlewares in the same subpath as Companion instead.

app.use(bodyParser.json());
app.use(session({ secret: getEnv("COMPANION_SECRET", crypto.randomUUID()) }));

const options = {
  filePath: getEnv("COMPANION_DATADIR", "./tmp"),
  secret: getEnv("COMPANION_SECRET", crypto.randomUUID()),
  preAuthSecret: getEnv("COMPANION_PREAUTH_SECRET"),
  redisUrl: getEnv("COMPANION_REDIS_URL"),
  redisOptions: {},
  redisPubSubScope: getEnv("COMPANION_REDIS_PUBSUB_SCOPE", "_companion"),
  corsOrigins:
    getEnv("COMPANION_CLIENT_ORIGINS", "true") === "true"
      ? true
      : getEnv("COMPANION_CLIENT_ORIGINS").split(","),
  uploadUrls: getEnv(
    "COMPANION_UPLOAD_URLS",
    "http://localhost:8081/upload"
  ).split(","),
  server: {
    host: getEnv("COMPANION_DOMAIN", "localhost"),
    protocol: getEnv("COMPANION_PROTOCOL", "http"),
    oauthDomain: getEnv("COMPANION_OAUTH_DOMAIN", "localhost:8081/remote"),
    path: getEnv("COMPANION_PATH", "/companion"),
    implicitPath: getEnv("COMPANION_IMPLICIT_PATH", "/companion"),
    validHosts: getEnv("COMPANION_DOMAINS", "localhost,127.0.0.1").split(","),
  },
  sendSelfEndpoint: getEnv("COMPANION_SELF_ENDPOINT", "localhost:8081/remote"),
  providerOptions: {
    drive: {
      key: getEnv("COMPANION_GOOGLE_KEY"),
      secret: getEnv("COMPANION_GOOGLE_SECRET"),
    },
    dropbox: {
      key: getEnv("COMPANION_DROPBOX_KEY"),
      secret: getEnv("COMPANION_DROPBOX_SECRET"),
    },
    instagram: {
      key: getEnv("COMPANION_INSTAGRAM_KEY"),
      secret: getEnv("COMPANION_INSTAGRAM_SECRET"),
    },
    facebook: {
      key: getEnv("COMPANION_FACEBOOK_KEY"),
      secret: getEnv("COMPANION_FACEBOOK_SECRET"),
    },
    onedrive: {
      key: getEnv("COMPANION_ONEDRIVE_KEY"),
      secret: getEnv("COMPANION_ONEDRIVE_SECRET"),
    },
    googlephotos: {
      key: getEnv("COMPANION_GOOGLE_KEY"),
      secret: getEnv("COMPANION_GOOGLE_SECRET"),
    },
    zoom: {
      key: getEnv("COMPANION_ZOOM_KEY"),
      secret: getEnv("COMPANION_ZOOM_SECRET"),
    },
  },
  enableGooglePickerEndpoint: getEnv(
    "COMPANION_ENABLE_GOOGLE_PICKER_ENDPOINT",
    true
  ),
  s3: {
    key: getEnv("COMPANION_AWS_KEY"),
    secret: getEnv("COMPANION_AWS_SECRET"),
    endpoint: getEnv("COMPANION_AWS_ENDPOINT"),
    bucket: getEnv("COMPANION_AWS_BUCKET", "uploadpilot"),
    forcePathStyle: getEnv("COMPANION_AWS_FORCE_PATH_STYLE", false),
    region: getEnv("COMPANION_AWS_REGION"),
    useAccelerateEndpoint: getEnv(
      "COMPANION_AWS_USE_ACCELERATE_ENDPOINT",
      false
    ),
    expires: getEnv("COMPANION_AWS_EXPIRES", 3600),
    acl: getEnv("COMPANION_AWS_ACL", "private"),
  },
  customProviders: {},
  logClientVersion: true,
  streamingUpload: getEnv("COMPANION_STREAMING_UPLOAD", true),
  metrics: getEnv("COMPANION_HIDE_METRICS", true),
  debug: true,
  maxFileSize: getEnv("COMPANION_MAX_FILE_SIZE", 1024 * 1024 * 1024), // 1GB
  periodicPingUrls: getEnv("COMPANION_PERIODIC_PING_URLS"),
  periodicPingInterval: getEnv("COMPANION_PERIODIC_PING_INTERVAL", 0),
  periodicPingStaticPayload: JSON.parse(
    getEnv("COMPANION_PERIODIC_PING_STATIC_PAYLOAD", '{"static": "payload"}')
  ),
  allowLocalUrls: getEnv("COMPANION_ALLOW_LOCAL_URLS", false),
  chunkSize: getEnv("COMPANION_CHUNK_SIZE", 1024 * 1024 * 5), // 5MB
};

const { app: companionApp } = companion.app(options);

app.get("/health", (req, res) => {
  res.send("OK");
});

app.use("/remote", companionApp);

const server = app.listen(getEnv("COMPANION_PORT", 8082), () => {
  console.log(`Listening companion on port ${getEnv("COMPANION_PORT", 8082)}`);
});

companion.socket(server);
