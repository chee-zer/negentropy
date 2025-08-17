-- +goose Up
CREATE    TABLE tasks (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name Text NOT NULL UNIQUE,
          color_hex text,
          completed BOOLEAN,
          daily_target INTEGER
          );

-- +goose Down
DROP      TABLE tasks;