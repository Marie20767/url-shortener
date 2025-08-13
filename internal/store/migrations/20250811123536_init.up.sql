CREATE TABLE users (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  name TEXT NOT NULL,
  age SMALLINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);