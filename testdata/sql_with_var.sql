-- This is a comment, I think this is probally getting lost
SELECT -- WE are going to get the following columns
    sl.source_id
    , sl.layer_id
    , sl.srid -- this is the SRID of the geomentry, import to have it be accurate
    , sl.description
    , sl.current_data_checksum
    , sl.last_data_check
    , sl.last_data_update
    , sl.created
    , sl.updated
    , $1
FROM -- this is the FROM block
    gis.source_layers sl -- nice little short cut
WHERE -- This is the WHERE BLOCK
    sl.source_id = $1 -- First Condition
    AND sl.layer_id = $2 -- Second Condition
;
--   END of the SQL file

