-- +goose Up
-- +goose StatementBegin
CREATE TABLE NETEX_StopPlace (
    id VARCHAR(255) PRIMARY KEY,
    modification VARCHAR(50),
    created_timestamp TIMESTAMPTZ,
    changed_timestamp TIMESTAMPTZ,
    valid_from_date TIMESTAMPTZ,
    valid_to_date TIMESTAMPTZ,
    name VARCHAR(255),
    name_lang VARCHAR(10),
    longitude DECIMAL(9,6),
    latitude DECIMAL(8,6),
    transport_mode VARCHAR(50),
    other_transport_modes TEXT,
    stop_place_type VARCHAR(50),
    weighting VARCHAR(50),
    UNIQUE (id)
);

CREATE TABLE NETEX_Quay (
    id VARCHAR(255) PRIMARY KEY,
    changed_timestamp TIMESTAMPTZ,
    name VARCHAR(255),
    name_lang VARCHAR(10),
    longitude DECIMAL(9,6),
    latitude DECIMAL(8,6),
    postal_region VARCHAR(50),
    site_ref_stopplace_id VARCHAR(255) REFERENCES NETEX_StopPlace(id) ON DELETE SET NULL,
    transport_mode VARCHAR(50),
    UNIQUE (id)
);

CREATE TABLE NETEX_StopPlace_QuayRef (
    stop_place_id VARCHAR(255) REFERENCES NETEX_StopPlace(id) ON DELETE CASCADE,
    quay_id VARCHAR(255) REFERENCES NETEX_Quay(id) ON DELETE CASCADE,
    quay_ref_version VARCHAR(50),
    PRIMARY KEY (stop_place_id, quay_id),
    UNIQUE (stop_place_id, quay_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE NETEX_StopPlace;
DROP TABLE NETEX_Quay;
DROP TABLE NETEX_StopPlace_QuayRef;
-- +goose StatementEnd
