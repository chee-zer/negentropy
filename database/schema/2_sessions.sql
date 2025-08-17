-- +goose Up
CREATE    TABLE sessions (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          start_time TIMESTAMP NOT NULL,
          duration_seconds INTEGER,
          task_id INTEGER NOT NULL,
          FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
          );

-- +goose Down
DROP      TABLE sessions;