-- Add up migration script here
-- User Table
CREATE TABLE IF NOT EXISTS Student (
    student_id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Subject Table
CREATE TABLE IF NOT EXISTS Subject (
    subject_id SERIAL PRIMARY KEY,
    student_id TEXT NOT NULL REFERENCES Student(student_id) ON DELETE CASCADE ,
    subject_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sync_status TEXT DEFAULT 'pending'
);

-- Test Table
CREATE TABLE IF NOT EXISTS Test (
    test_id SERIAL PRIMARY KEY,
    subject_id INT NOT NULL REFERENCES Subject(subject_id) ON DELETE CASCADE,
    test_name TEXT NOT NULL,
    test_date DATE NOT NULL,
    max_marks INT NOT NULL,
    obtained_marks INT NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sync_status TEXT DEFAULT 'pending',
    UNIQUE(subject_id, test_name, test_date, max_marks)
);

-- -- SyncMetadata
-- CREATE TABLE IF NOT EXISTS SyncMetadata(
--     sync_id SERIAL PRIMARY KEY,
--     entity_id INTEGER NOT NULL,
--     entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('subject', 'test')),
--     student_id TEXT REFERENCES Student(student_id) ON DELETE CASCADE,
--     last_sync_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     conflict_resolved BOOLEAN DEFAULT FALSE,
--     -- FOREIGN KEY (entity_id, entity_type) REFERENCES (
--     --     SELECT subject_id AS entity_id, 'subject' AS entity_type FROM Subject
--     --     UNION ALL
--     --     SELECT test_id AS entity_id, 'test' AS entity_type FROM Test
--     -- ) ON DELETE CASCADE
-- );