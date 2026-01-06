CREATE TABLE keys (
  id VARCHAR(8) CHECK (LENGTH(id) = 8) PRIMARY KEY,
  used BOOLEAN DEFAULT false
);

CREATE TABLE urls (
  short VARCHAR(8) PRIMARY KEY,
  long TEXT NOT NULL,
  expiry TIMESTAMPTZ NULL,
  user_id UUID DEFAULT gen_random_uuid() NOT NULL,

  CONSTRAINT fk_short FOREIGN KEY(short) REFERENCES keys(id)
);
