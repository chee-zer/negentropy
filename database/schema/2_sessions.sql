-- +goose Up
CREATE    TABLE sessions (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          start_time TEXT NOT NULL,
          end_time TEXT,
          task_id INTEGER NOT NULL,
          FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
          );

-- +goose Down
DROP      TABLE sessions;   