-- +goose Up
-- +goose StatementBegin
CREATE TABLE requests(
  id INTEGER PRIMARY KEY,
  requested_path TEXT,
  requested_method TEXT,
  matched_path TEXT,
  response_body JSON,
  response_status INTEGER,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE requests;
-- +goose StatementEnd
