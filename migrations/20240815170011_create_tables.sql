-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS banners (
  id SERIAL PRIMARY KEY,
  descr TEXT
);


CREATE TABLE IF NOT EXISTS banners (
  id SERIAL PRIMARY KEY,
  descr TEXT
);

--INSERT INTO events (title, event_date, user_id) VALUES ( 'Some event', '2024-08-14 16:50:36', '1');
--select * from events;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

--DROP TABLE events;
-- +goose StatementEnd
