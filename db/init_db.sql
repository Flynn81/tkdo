CREATE USER IF NOT EXISTS tk;

CREATE DATABASE IF NOT EXISTS tkdo;

GRANT ALL ON DATABASE tkdo TO tk;

set database = tkdo;

CREATE TABLE IF NOT EXISTS task (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(256) NOT NULL,
  type varchar(64) NOT NULL,
  status varchar(256) NOT NULL
);
