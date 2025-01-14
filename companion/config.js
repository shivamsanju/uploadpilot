import dotenv from 'dotenv';

// Load environment variables from .env file
dotenv.config();

const config = {
    mandatory: {
        companionSecret: process.env.COMPANION_SECRET,
        companionDomain: process.env.COMPANION_DOMAIN,
        companionDataDir: process.env.COMPANION_DATADIR,
        uploadPilotUrl: process.env.UPLOADPILOT_URL || 'http://localhost:8081',
    },
    optional: {
        companionProtocol: process.env.COMPANION_PROTOCOL || 'http',
        companionPort: process.env.COMPANION_PORT || 8080,
        companionPath: process.env.COMPANION_PATH || '',
        companionHideWelcome: process.env.COMPANION_HIDE_WELCOME === 'false',
        companionHideMetrics: process.env.COMPANION_HIDE_METRICS === 'false',
        companionLoggerProcessName: process.env.COMPANION_LOGGER_PROCESS_NAME || 'companion',
        companionImplicitPath: process.env.COMPANION_IMPLICIT_PATH || '',
        companionClientOrigins: process.env.COMPANION_CLIENT_ORIGINS?.split(',') || [],
        companionClientOriginsRegex: process.env.COMPANION_CLIENT_ORIGINS_REGEX || null,
        companionRedisUrl: process.env.COMPANION_REDIS_URL || null,
        integrations: {
            dropbox: {
                key: process.env.COMPANION_DROPBOX_KEY,
                secret: process.env.COMPANION_DROPBOX_SECRET,
            },
            box: {
                key: process.env.COMPANION_BOX_KEY,
                secret: process.env.COMPANION_BOX_SECRET,
            },
            google: {
                key: process.env.COMPANION_GOOGLE_KEY,
                secret: process.env.COMPANION_GOOGLE_SECRET,
            },
            instagram: {
                key: process.env.COMPANION_INSTAGRAM_KEY,
                secret: process.env.COMPANION_INSTAGRAM_SECRET,
            },
            facebook: {
                key: process.env.COMPANION_FACEBOOK_KEY,
                secret: process.env.COMPANION_FACEBOOK_SECRET,
            },
            onedrive: {
                key: process.env.COMPANION_ONEDRIVE_KEY,
                secret: process.env.COMPANION_ONEDRIVE_SECRET,
            },
            zoom: {
                key: process.env.COMPANION_ZOOM_KEY,
                secret: process.env.COMPANION_ZOOM_SECRET,
            },
            s3: {
                key: process.env.COMPANION_AWS_KEY,
                secret: process.env.COMPANION_AWS_SECRET,
                bucket: process.env.COMPANION_AWS_BUCKET,
                region: process.env.COMPANION_AWS_REGION,
                useAccelerateEndpoint: process.env.COMPANION_AWS_USE_ACCELERATE_ENDPOINT === 'true',
                expires: parseInt(process.env.COMPANION_AWS_EXPIRES || '800', 10),
                acl: process.env.COMPANION_AWS_ACL || 'private',
                prefix: process.env.COMPANION_AWS_PREFIX || '',
            },
        },
        oauthDomain: process.env.COMPANION_OAUTH_DOMAIN || '',
        domains: process.env.COMPANION_DOMAINS?.split(',') || [],
        selfEndpoint: process.env.COMPANION_SELF_ENDPOINT || '',
        uploadUrls: process.env.COMPANION_UPLOAD_URLS?.split(',') || [],
        streamingUpload: process.env.COMPANION_STREAMING_UPLOAD === 'true',
        allowLocalUrls: process.env.COMPANION_ALLOW_LOCAL_URLS === 'true',
        maxFileSize: parseInt(process.env.COMPANION_MAX_FILE_SIZE || '0', 10),
        chunkSize: parseInt(process.env.COMPANION_CHUNK_SIZE || '0', 10),
        periodicPingUrls: process.env.COMPANION_PERIODIC_PING_URLS?.split(',') || [],
        periodicPingInterval: parseInt(process.env.COMPANION_PERIODIC_PING_INTERVAL || '0', 10),
        periodicPingStaticPayload: process.env.COMPANION_PERIODIC_PING_STATIC_JSON_PAYLOAD || null,
        redisExpressSessionPrefix: process.env.COMPANION_REDIS_EXPRESS_SESSION_PREFIX || 'sess:',
        preauthSecret: process.env.COMPANION_PREAUTH_SECRET || '',
    },
};

export { config };
