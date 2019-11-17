CREATE TABLE revisions(
    "id"            VARCHAR(11) NOT NULL,
    "name"          TEXT NOT NULL,
    "url"           TEXT NOT NULL,
    "file_hash_sum" VARCHAR(32),
    "removed"       INT NOT NULL,
    "added"         INT NOT NULL,
    "created_at"    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX revisions_id_idx ON revisions("id");

CREATE TYPE VEHICLE_STATUS_T AS ENUM ('stolen', 'removed', 'failed');

CREATE TABLE vehicles(
    "revision_id"    VARCHAR(11) NOT NULL REFERENCES revisions(id),
    "id"             TEXT NOT NULL,
    "brand"          TEXT,
    "color"          TEXT,
    "number"         TEXT,
    "body_number"    TEXT,
    "chassis_number" TEXT,
    "engine_number"  TEXT,
    "ovd"            TEXT             NOT NULL,
    "kind"           TEXT             NOT NULL,
    "status"         VEHICLE_STATUS_T NOT NULL DEFAULT 'stolen',
    "theft_date"     TIMESTAMP NOT NULL,
    "insert_date"    TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX vehicles_id_idx             ON vehicles("id");
CREATE        INDEX vehicles_number_idx         ON vehicles("number");
CREATE        INDEX vehicles_body_number_idx    ON vehicles("body_number");
CREATE        INDEX vehicles_chassis_number_idx ON vehicles("chassis_number");
CREATE        INDEX vehicles_engine_number_idx  ON vehicles("engine_number");
