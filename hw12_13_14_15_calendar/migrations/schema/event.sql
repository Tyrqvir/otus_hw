-- TODO Refactor to go-migrations
CREATE TABLE IF NOT EXISTS events (
  id varchar(100) NOT NULL,
  title varchar(100) NOT NULL,
  start date NOT NULL,
  end date NOT NULL,
  description varchar(255),
  owner_id varchar(100) NOT NULL,
  notification_before SMALLINT,
  primary key (id)
);
