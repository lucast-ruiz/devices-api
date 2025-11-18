CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE devices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  brand TEXT NOT NULL,
  state TEXT NOT NULL CHECK (state IN ('available','in-use','inactive')),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
