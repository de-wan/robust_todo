-- 000001_add_todo_table.up.sql
CREATE TABLE todo(
    uuid VARCHAR(36) PRIMARY KEY,
    value VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT NOW() NOT NULL,
    done_at  DATETIME,
    deadline DATETIME,
    archived_at DATETIME
)