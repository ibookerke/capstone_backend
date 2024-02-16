CREATE SCHEMA IF NOT EXISTS "har";

-- Users Table
CREATE TABLE har.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    surname VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Camera Feeds Table
CREATE TABLE har.cameras (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    camera_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- User Cameras Table
CREATE TABLE har.user_cameras (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES har.users(id) ON DELETE CASCADE,
    camera_id INTEGER REFERENCES har.cameras(id) ON DELETE CASCADE
);

-- Notifications Table
CREATE TABLE har.notifications (
    id SERIAL PRIMARY KEY,
    camera_id INTEGER REFERENCES har.cameras(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES har.users(id) ON DELETE CASCADE,
    notification_message TEXT NOT NULL,
    arrival_time TIMESTAMP WITH TIME ZONE NOT NULL,
    type INTEGER NOT NULL
);

-- User Notification Settings Table
CREATE TABLE har.user_notification_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES har.users(id) ON DELETE CASCADE,
    notification_option VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE
);
