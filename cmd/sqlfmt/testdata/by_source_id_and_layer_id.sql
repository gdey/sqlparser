SELECT
    sl.source_id
    , sl.layer_id
    , sl.srid
    , sl.description
    , sl.current_data_checksum
    , sl.last_data_check
    , sl.last_data_update
    , sl.created
    , sl.updated
FROM
    gis.source_layers sl
WHERE
    sl.source_id = $1
    AND sl.layer_id = $2;

