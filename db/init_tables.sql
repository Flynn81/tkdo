GRANT ALL ON DATABASE tkdo TO tk;

CREATE TABLE IF NOT EXISTS task (
  id UUID PRIMARY KEY,
  name varchar(256) NOT NULL,
  type varchar(64) NOT NULL,
  status varchar(256) NOT NULL,
  user_id UUID
);

CREATE TABLE IF NOT EXISTS task_user (
  id UUID PRIMARY KEY,
  name varchar(256) NOT NULL,
  email varchar(256) NOT NULL UNIQUE,
  hash bytea,
  status varchar(64) NOT NULL
);
