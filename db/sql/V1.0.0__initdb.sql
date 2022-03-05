--
-- Initial setup
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- User Account
--
CREATE TABLE IF NOT EXISTS user_account
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username      varchar(24) unique       not null,
    password      varchar                  not null,
    email         varchar unique,
    role_id       UUID, -- TODO("Create foreign key when create role table")
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL
);
--
-- Skill Config
--
CREATE TABLE IF NOT EXISTS skill_config
(
    id          CHAR(3) PRIMARY KEY,
    name        VARCHAR UNIQUE NOT NULL,
    description VARCHAR
);

--
-- Skill
--
CREATE TABLE IF NOT EXISTS skill
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    skill_id      CHAR(3)                    NOT NULL,
    user_id       UUID, -- TODO("Create foreign key when create user table")
    exp           INTEGER          DEFAULT 0 NOT NULL,
    txp           INTEGER          DEFAULT 0 NOT NULL,
    level         SMALLINT         DEFAULT 1,
    date_created  TIMESTAMP WITH TIME ZONE   NOT NULL,
    date_modified TIMESTAMP WITH TIME ZONE   NOT NULL,

    CONSTRAINT fk_skill_id
        FOREIGN KEY (skill_id)
            REFERENCES skill_config (id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_account (id)
);


