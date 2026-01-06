CREATE TABLE keys (
  id VARCHAR(8) CHECK (LENGTH(id) = 8) PRIMARY KEY,
  used BOOLEAN DEFAULT false
);

CREATE TABLE urls (
  short VARCHAR(8) PRIMARY KEY,
  long string NOT NULL,
  expiry TIMESTAMPTZ,
  user_id PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,

  CONSTRAINT fk_short_url FOREIGN KEY(short_url) REFERENCES keys(id)
)