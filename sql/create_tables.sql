CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SELECT uuid_generate_v1();

CREATE TABLE IF NOT EXISTS flowchart (
    id           uuid DEFAULT uuid_generate_v4 (),
    created_at   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title        varchar NOT NULL,
    key        varchar(50) UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS node (
    id           uuid DEFAULT uuid_generate_v4 (),
    created_at   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data         JSONB NOT NULL,
    position     JSONB NOT NULL,
    width        int NOT NULL,
    height       int NOT NULL,
    position_absolute JSONB,
    selected     boolean DEFAULT false,
    dragging     boolean DEFAULT false,
    internal_id  int NOT NULL,
    parent_id    int DEFAULT 0,
    flowchart_id uuid NOT NULL,
    type         varchar(30) NOT NULL,
    CONSTRAINT   flowchart_pk FOREIGN KEY (flowchart_id) REFERENCES flowchart(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT   node_pk PRIMARY KEY (id)
);

