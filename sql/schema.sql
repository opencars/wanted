

--
CREATE TABLE resource_revisions(
    "id"            SERIAL,
    "name"          TEXT NOT NULL,
    "url"           TEXT NOT NULL,
    "file_hash_sum" VARCHAR(32) NOT NULL,
);

--
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--
CREATE TYPE STATE_T AS ENUM ('STOLEN', 'REMOVED');

--
CREATE TABLE wanted_vehicles(
    "file_revision_id" SERIAL NOT NULL,
    "id"             TEXT NOT NULL,
    "brand"          TEXT,
    "color"          TEXT,
    "number"         TEXT,
    "body_number"    TEXT,
    "chassis_number" TEXT,
    "engine_number"  TEXT,
    "ovd"            TEXT      NOT NULL,
    "kind"           TEXT      NOT NULL,
    "state"          STATE_T   NOT NULL DEFAULT 'STOLEN'
    "theft_date"     TIMESTAMP NOT NULL,
    "insert_date"    TIMESTAMP NOT NULL,
    "created_at"     TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"     TIMESTAMP NOT NULL DEFAULT NOW(),
);

CREATE UNIQUE INDEX wanted_vehicles_id_idx             ON wanted_vehicles("id");
CREATE        INDEX wanted_vehicles_number_idx         ON wanted_vehicles("number");
CREATE        INDEX wanted_vehicles_body_number_idx    ON wanted_vehicles("body_number");
CREATE        INDEX wanted_vehicles_chassis_number_idx ON wanted_vehicles("chassis_number");
CREATE        INDEX wanted_vehicles_engine_number_idx  ON wanted_vehicles("engine_number");
