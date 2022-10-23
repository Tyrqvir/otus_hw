-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id varchar(100) NOT NULL,
    title varchar(100) NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    description varchar(255),
    owner_id varchar(100) NOT NULL,
    notification_date timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NULL,
    is_notified smallint DEFAULT 0,
    primary key (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
