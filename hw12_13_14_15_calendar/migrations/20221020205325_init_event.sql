-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
                      id serial NOT NULL,
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

CREATE UNIQUE INDEX events_owner_id_with_start_date_idx ON events USING btree (owner_id, start_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
