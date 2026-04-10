-- +goose Up
CREATE TABLE rides (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  driver_id   UUID NOT NULL REFERENCES users(id),
  origin      TEXT NOT NULL,
  destination TEXT NOT NULL,
  departure   TIMESTAMP NOT NULL,
  seats       INT NOT NULL,
  created_at  TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS rides;
