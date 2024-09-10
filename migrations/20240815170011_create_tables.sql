-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner (
  id SERIAL PRIMARY KEY,
  descr TEXT
);


CREATE TABLE IF NOT EXISTS slot (
  id SERIAL PRIMARY KEY,
  descr TEXT
);


CREATE TABLE IF NOT EXISTS social_group (
  id SERIAL PRIMARY KEY,
  descr TEXT
);


CREATE TABLE IF NOT EXISTS rotation (
  banner INT,
  slot INT,

  CONSTRAINT fk_banner
      FOREIGN KEY(banner) 
        REFERENCES banner(id),
  CONSTRAINT fk_slot
      FOREIGN KEY(slot) 
        REFERENCES slot(id)
);


CREATE UNIQUE INDEX rotation_index on rotation(banner, slot);


CREATE TABLE IF NOT EXISTS statistic (
  banner INT,
  slot INT,
  clicks INT DEFAULT 0,
  shows INT DEFAULT 0,
  s_group INT,

  CONSTRAINT fk_group_stat
      FOREIGN KEY(s_group) 
        REFERENCES social_group(id),
  CONSTRAINT fk_banner_stat
      FOREIGN KEY(banner) 
        REFERENCES banner(id),
  CONSTRAINT fk_slot_stat
      FOREIGN KEY(slot) 
        REFERENCES slot(id)

);

CREATE UNIQUE INDEX stat_index on statistic(banner, slot, s_group);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

--DROP TABLE events;
-- +goose StatementEnd
