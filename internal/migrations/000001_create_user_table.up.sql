CREATE TABLE IF NOT EXISTS tasks(
    task_id TEXT PRIMARY KEY, 
    user_id TEXT NOT NULL,
    task_title TEXT NOT NULL, 
    done_status BOOLEAN NOT NULL CHECK (done_status IN (0, 1)),
    added_at DATETIME NOT NULL, 
    modified_at DATETIME);