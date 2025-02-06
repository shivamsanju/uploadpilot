import express from 'express'
import bodyParser from 'body-parser'
import session from 'express-session'
import companion from '@uppy/companion'
import { getConfig } from "./config.js"

const config = getConfig()
const app = express()

// Companion requires body-parser and express-session middleware.
// You can add it like this if you use those throughout your app.
// If you are using something else in your app, you can add these
// middlewares in the same subpath as Companion instead.

app.use(bodyParser.json())
app.use(session({ secret: config.mandatory.companionSecret }))


const options = {
    providerOptions: {
        drive: {
            key: config.optional.integrations.google.key,
            secret: config.optional.integrations.google.secret,
        },
        dropbox: {
            key: config.optional.integrations.dropbox.key,
            secret: config.optional.integrations.dropbox.secret,
        },
        instagram: {
            key: config.optional.integrations.instagram.key,
            secret: config.optional.integrations.instagram.secret,
        },
        facebook: {
            key: config.optional.integrations.facebook.key,
            secret: config.optional.integrations.facebook.secret,
        },
        onedrive: {
            key: config.optional.integrations.onedrive.key,
            secret: config.optional.integrations.onedrive.secret,
        },
    },
    s3: {
        getKey: (req, filename, metadata) => `${crypto.randomUUID()}-${filename}`,
        key: config.optional.integrations.s3.key,
        secret: config.optional.integrations.s3.secret,
        bucket: config.optional.integrations.s3.bucket || 'bucket-name',
        region: config.optional.integrations.s3.region || 'us-east-1',
        useAccelerateEndpoint: config.optional.integrations.s3.useAccelerateEndpoint || false,
        expires: config.optional.integrations.s3.expires || 3600,
        acl: config.optional.integrations.s3.acl || 'private',
    },
    server: {
        host: `${config.mandatory.companionDomain}:${config.optional.companionPort}`,
        protocol: config.optional.companionProtocol || 'http',
        path: '/remote',
        oauthDomain: config.optional.oauthDomain,
        validHosts: config.optional.validHosts || ["localhost:8082"],
    },
    filePath: config.mandatory.companionDataDir || '/tmp',
    sendSelfEndpoint: config.optional.selfEndpoint || 'localhost:3020',
    secret: config.mandatory.companionSecret || 'mysecret',
    uploadUrls: config.optional.uploadUrls || [],
    debug: true,
    metrics: !config.optional.companionHideMetrics,
    streamingUpload: config.optional.streamingUpload,
    allowLocalUrls: config.optional.allowLocalUrls,
    maxFileSize: config.optional.maxFileSize || 100000000,
    chunkSize: config.optional.chunkSize || 10000000,
    periodicPingUrls: config.optional.periodicPingUrls || [],
    periodicPingInterval: config.optional.periodicPingInterval || 60000,
    periodicPingStaticPayload: JSON.parse(config.optional.periodicPingStaticPayload || '{"static": "payload"}'),
    corsOrigins: config.optional.companionClientOrigins || true,
};

console.log(options)
const { app: companionApp } = companion.app(options)

app.get("/health", (req, res) => {
    res.send("OK")
})

app.use("/remote", companionApp)

const server = app.listen(config.optional.companionPort, () => {
    console.log(`Listening on port ${config.optional.companionPort}`)
})
companion.socket(server)
