--
-- Initial setup
--
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- Role
--
CREATE TABLE IF NOT EXISTS role
(
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    name          varchar(24) unique       NOT NULL,
    permissions   varchar[]                NOT NULL,
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--
-- Role
--
CREATE TABLE IF NOT EXISTS permission
(
    id           VARCHAR                  NOT NULL PRIMARY KEY,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--
-- User Account
--
CREATE TABLE IF NOT EXISTS user_account
(
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    username      varchar(24) unique       NOT NULL,
    password      varchar                  NOT NULL,
    email         varchar unique,
    role_id       UUID,
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_id
        FOREIGN KEY (role_id)
            REFERENCES role (id)
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
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    skill_id      CHAR(3)                  NOT NULL,
    user_id       UUID,
    exp           INTEGER                  NOT NULL DEFAULT 0,
    txp           INTEGER                  NOT NULL DEFAULT 0,
    level         SMALLINT                          DEFAULT 1,
    date_created  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_skill_id
        FOREIGN KEY (skill_id)
            REFERENCES skill_config (id),
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES user_account (id)
);


