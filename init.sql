CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY UNIQUE,
  username VARCHAR(15) NOT NULL UNIQUE,
  user_password VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  email_verified BOOLEAN DEFAULT FALSE,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  user_bio VARCHAR(255),
  birthday VARCHAR(10),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  verified BOOLEAN DEFAULT FALSE,
  profile_picture_url VARCHAR(255),
  notifications_enabled BOOLEAN DEFAULT TRUE,
  flagged INTEGER DEFAULT 0,
  rank INTEGER DEFAULT 0,
  creator BOOLEAN DEFAULT FALSE,
  salt BYTEA,
  latitude DECIMAL(9,6),
  longitude DECIMAL(9,6),
  session_token VARCHAR(255), 
);

CREATE TABLE IF NOT EXISTS followings (
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(user_id),
    FOREIGN KEY (following_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS blocked_users (
    blocker_id UUID NOT NULL,
    blocked_user_id UUID NOT NULL,
    PRIMARY KEY (blocker_id, blocked_user_id),
    FOREIGN KEY (blocker_id) REFERENCES users(user_id),
    FOREIGN KEY (blocked_user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    content TEXT,
    photo_url VARCHAR(255),
    caption TEXT,
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    location_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS likes (
    like_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id)
);

CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id)
);

CREATE TABLE IF NOT EXISTS notifications  ( -- (e.g., "like", "comment", "follow", etc.).
    notification_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- CREATE TABLE IF NOT EXISTS post_tags (
--     post_id UUID NOT NULL,
--     tag_id UUID NOT NULL,
--     PRIMARY KEY (post_id, tag_id),
--     FOREIGN KEY (post_id) REFERENCES posts(post_id),
--     FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
-- );

CREATE TABLE IF NOT EXISTS flagged_accounts (
    user_id UUID PRIMARY KEY,
    flag_count INTEGER DEFAULT 0,
    is_suspended BOOLEAN DEFAULT false,
    reason TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    content_type VARCHAR(10) NOT NULL,
    content TEXT,
    media_url VARCHAR(255), -- For pictures and videos
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sender_id) REFERENCES users(user_id),
    FOREIGN KEY (receiver_id) REFERENCES users(user_id)
);
