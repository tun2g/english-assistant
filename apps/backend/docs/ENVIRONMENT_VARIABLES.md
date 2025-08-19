# Environment Variables Documentation

This document describes all environment variables used by the English Learning Assistant backend application.

## Quick Setup

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Update the required API keys in `.env`:
   - `YOUTUBE_API_KEY` - **Required** for YouTube video processing
   - `GEMINI_API_KEY` - **Required** for translation services

**Note**: This configuration matches the existing `app.yaml` structure and only includes APIs that are actually implemented in the codebase.

## Core Application Configuration

### Database Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_HOST` | Yes | `localhost` | PostgreSQL database host |
| `DB_PORT` | Yes | `5434` | PostgreSQL database port |
| `DB_USER` | Yes | `postgres` | PostgreSQL username |
| `DB_PASSWORD` | Yes | `postgres` | PostgreSQL password |
| `DB_NAME` | Yes | `app_backend_dev` | PostgreSQL database name |

### Application Settings

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_ENVIRONMENT` | No | `development` | Application environment (development/production) |
| `APP_PORT` | No | `8000` | Port for the HTTP server |
| `APP_LOG_LEVEL` | No | `debug` | Logging level (debug/info/warn/error) |

### Security Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | Yes | - | JWT signing secret (minimum 32 characters) |
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000,http://localhost:3001` | Comma-separated list of allowed CORS origins |

## External API Services

### YouTube Data API v3 (Required)

The YouTube Data API is required for fetching video information, metadata, and transcripts.

| Variable | Required | Description |
|----------|----------|-------------|
| `YOUTUBE_API_KEY` | **Yes** | YouTube Data API v3 key |

**Setup Instructions:**
1. Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Create a new project or select existing one
3. Enable "YouTube Data API v3"
4. Create credentials (API Key)
5. Restrict the API key to YouTube Data API v3 for security

**API Usage:**
- Video metadata retrieval
- Transcript/caption fetching
- Available languages detection
- Video capabilities checking

### Google Gemini AI (Required)

Google Gemini AI is used for translation services and language detection.

| Variable | Required | Description |
|----------|----------|-------------|
| `GEMINI_API_KEY` | **Yes** | Google Gemini API key |

**Setup Instructions:**
1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Create a new API key
3. Copy the key to your environment file

### Encryption Key (Optional)

| Variable | Required | Description |
|----------|----------|-------------|
| `ENCRYPTION_KEY` | No | Encryption key for sensitive data storage |

**Note**: Only include external services that are actually implemented in the backend codebase. The configuration matches the existing `app.yaml` structure.

## Production Considerations

### Security

1. **JWT Secret**: Use a cryptographically secure random string (minimum 32 characters)
2. **API Keys**: Store API keys securely and rotate them regularly
3. **Database**: Use strong passwords and enable SSL in production
4. **CORS**: Restrict CORS origins to your actual domains

### Performance

1. **Rate Limits**: Adjust rate limits based on your API quotas and usage patterns
2. **Caching**: Enable caching to reduce API costs and improve response times
3. **Batch Processing**: Use batch translation for better performance with large transcripts

### Monitoring

1. **Log Level**: Set to `info` or `warn` in production
2. **API Usage**: Monitor API usage to avoid quota exceeded errors
3. **Error Handling**: Implement proper error handling and fallback mechanisms

## Environment-Specific Configurations

### Development
```bash
APP_ENVIRONMENT=development
APP_LOG_LEVEL=debug
ENABLE_API_CACHE=true
```

### Production
```bash
APP_ENVIRONMENT=production
APP_LOG_LEVEL=info
ENABLE_API_CACHE=true
DB_SSL_MODE=require
```

### Testing
```bash
APP_ENVIRONMENT=test
APP_LOG_LEVEL=warn
ENABLE_API_CACHE=false
```

## Troubleshooting

### Common Issues

1. **YouTube API Quota Exceeded**: Increase `YOUTUBE_API_RATE_LIMIT` or upgrade your quota
2. **Translation Timeouts**: Reduce `TRANSLATION_BATCH_SIZE` or increase timeouts
3. **Cache Issues**: Clear cache by restarting the application or disabling cache temporarily

### Required Minimum Configuration

For the application to function, you must set at least:
```bash
YOUTUBE_API_KEY=your-key-here
GEMINI_API_KEY=your-key-here
JWT_SECRET=your-secure-secret-here
```

### API Key Testing

Test your API keys with these endpoints:
- YouTube: `curl "https://www.googleapis.com/youtube/v3/videos?id=dQw4w9WgXcQ&key=YOUR_KEY&part=snippet"`
- Gemini: Use the Google AI Studio interface to test

## Configuration Mapping

This environment configuration maps to the existing `app.yaml` structure:

```yaml
external_apis:
  youtube:
    api_key: "${YOUTUBE_API_KEY}"
    api_url: "https://www.googleapis.com/youtube/v3"
    rate_limit: 100
  
  gemini:
    api_key: "${GEMINI_API_KEY}" 
    api_url: "https://generativelanguage.googleapis.com"
    rate_limit: 60

encryption:
  key: "${ENCRYPTION_KEY:-default-dev-key-change-in-production}"
```