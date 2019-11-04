CREATE TYPE STATE_T AS ENUM ('STOLEN', 'REMOVED');

CREATE TABLE wanted_vehicles(
    "id"             TEXT NOT NULL,
    "brand"          TEXT,
    "color"          TEXT,
    "number"         TEXT,
    "body_number"    TEXT,
    "chassis_number" TEXT,
    "engine_number"  TEXT,
    "ovd"            TEXT      NOT NULL,
    "kind"           TEXT      NOT NULL,
    "theft_date"     TIMESTAMP NOT NULL,
    "insert_date"    TIMESTAMP NOT NULL,
    "state"          STATE_T   NOT NULL DEFAULT 'STOLEN'
);

CREATE UNIQUE INDEX wanted_vehicles_id_idx             ON wanted_vehicles("id");
CREATE        INDEX wanted_vehicles_number_idx         ON wanted_vehicles("number");
CREATE        INDEX wanted_vehicles_body_number_idx    ON wanted_vehicles("body_number");
CREATE        INDEX wanted_vehicles_chassis_number_idx ON wanted_vehicles("chassis_number");
CREATE        INDEX wanted_vehicles_engine_number_idx  ON wanted_vehicles("engine_number");
