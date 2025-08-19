-- Create video transcript cache table
CREATE TABLE video_transcript_cache (
    id SERIAL PRIMARY KEY,
    video_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    language VARCHAR(10) NOT NULL,
    content TEXT NOT NULL,
    source VARCHAR(50) DEFAULT 'manual',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for video transcript cache
CREATE INDEX idx_video_transcript_cache_video_id ON video_transcript_cache(video_id);
CREATE INDEX idx_video_transcript_cache_expires_at ON video_transcript_cache(expires_at);
CREATE UNIQUE INDEX idx_video_transcript_cache_unique ON video_transcript_cache(video_id, provider, language);

-- Create video translation cache table
CREATE TABLE video_translation_cache (
    id SERIAL PRIMARY KEY,
    video_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    content TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for video translation cache
CREATE INDEX idx_video_translation_cache_video_id ON video_translation_cache(video_id);
CREATE INDEX idx_video_translation_cache_expires_at ON video_translation_cache(expires_at);
CREATE UNIQUE INDEX idx_video_translation_cache_unique ON video_translation_cache(video_id, provider, source_lang, target_lang);

-- Create user API keys table
CREATE TABLE user_api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_name VARCHAR(50) NOT NULL,
    encrypted_key TEXT NOT NULL,
    key_hash VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for user API keys
CREATE INDEX idx_user_api_keys_user_id ON user_api_keys(user_id);
CREATE INDEX idx_user_api_keys_key_hash ON user_api_keys(key_hash);
CREATE UNIQUE INDEX idx_user_api_keys_unique ON user_api_keys(user_id, service_name);

-- Create video analytics table
CREATE TABLE video_analytics (
    id SERIAL PRIMARY KEY,
    video_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    source_language VARCHAR(10),
    target_language VARCHAR(10),
    processing_time_ms BIGINT,
    success BOOLEAN DEFAULT true,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for video analytics
CREATE INDEX idx_video_analytics_video_id ON video_analytics(video_id);
CREATE INDEX idx_video_analytics_user_id ON video_analytics(user_id);
CREATE INDEX idx_video_analytics_action ON video_analytics(action);
CREATE INDEX idx_video_analytics_created_at ON video_analytics(created_at);