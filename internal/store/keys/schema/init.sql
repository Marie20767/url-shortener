CREATE TABLE keys (
  key_value VARCHAR(8) CHECK (LENGTH("key_value") = 8) PRIMARY KEY,
  used BOOLEAN DEFAULT false
);