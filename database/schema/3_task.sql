-- +goose Up
INSERT INTO tasks (
    id,
    name,
    color_hex,
    completed,
    daily_target
) VALUES (
    0,
    "ENTROPY",
    "#FF0000",
    false,
    16
);
