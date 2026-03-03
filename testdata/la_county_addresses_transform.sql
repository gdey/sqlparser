/*
	Description: Los Angeles County Addresses Transform Script
	Table names:
		- parcels_data: parcels geometry data from shapefile
		- parcels_csv: parcels data from csv files
		- addresses: final combined data to be loaded into addresses layer
*/


-- Create some indexes to speed up merging the tables
CREATE INDEX parcels_data_ain_idx ON parcels_data(ain);
CREATE INDEX parcels_data_parcel_typ_idx ON parcels_data(parcel_typ);
CREATE INDEX parcels_csv_ain_idx ON parcels_csv(ain);

-- Drop rows that don't have a valid ain, geom, or parcel type
DELETE FROM parcels_data WHERE ain IS NULL OR cast(ain as int) = 0 OR length(ain) > 10;-- OR parcel_typ = 0;
DELETE FROM parcels_data WHERE geometry IS NULL;
DELETE FROM parcels_csv WHERE ain IS NULL OR cast(ain as int) = 0 OR length(ain) > 10;

-- Before we merge tables we should eliminate any duplicate ain rows from parcels_csv
-- The parcels_csv table should never contain dupes so this is only here as a safeguard. Note
-- that duplicate ain rows in parcels_data are expected and will be dealt with later by
-- creating geometry unions.
DELETE FROM parcels_csv WHERE rowid NOT IN
	(SELECT min(rowid) FROM parcels_csv GROUP BY ain);

-- delete "Toxic Geometries" which are geometries
-- that contain too few points. At the time of adding this
-- condition the dataset has 2 hits with this condition.
-- these geometries are causing the ST_MakeValid fail
DELETE FROM parcels_data WHERE instr(IsValidReason(geometry), 'Toxic Geometry') > 0;

-- Sanitize all geometries
UPDATE parcels_data
	SET geometry = ST_MakeValid(geometry);
DELETE FROM parcels_data WHERE geometry IS NULL;

-- Create the addresses table by unioning geometries, same as parcels
CREATE TABLE addresses AS
	SELECT
		trim(data.ain) reference_id,
		'06' state_fips_code,
		'037' county_fips_code,
		trim(data.ain) parcel_apn,
		ST_Multi(ST_Union(data.geometry)) geom_src,
		nullif(cast(csv.situs_house_no as int), 0) street_number, 
		nullif(trim(csv.fraction), '') street_number_fraction,
		nullif(trim(csv.direction), '') street_pre_direction,
		nullif(trim(csv.street_name), '') street_name,
		cast(null as text) street_suffix,
		cast(null as text) street_post_direction,
		nullif(trim(csv.unit), '') unit,
		nullif(trim(csv.zip), '') postal,
		nullif(trim(rtrim(upper(trim(csv.city_state)), 'CA')),'') jurisdiction,
		nullif(trim(rtrim(upper(trim(csv.city_state)), 'CA')),'') postal_city,
		'CA' state,
		'USA' country
	FROM
		parcels_data AS data
	LEFT JOIN
		parcels_csv AS csv USING(ain)
	WHERE 
		jurisdiction IS NOT NULL
	GROUP BY 
		ain
;

UPDATE addresses SET jurisdiction = proper(jurisdiction) WHERE jurisdiction IS NOT NULL;

UPDATE addresses SET postal_city = proper(postal_city) WHERE postal_city IS NOT NULL;

-- Add a hyphen to 9-character zip codes
CREATE INDEX addresses_postal_idx ON addresses(postal);
UPDATE addresses
	SET postal = printf('%s-%s', substr(postal, 1, 5), substr(postal, -4, 4))
	WHERE instr(postal, '-') = 0
	AND length(postal) = 9
;

-- Recover the geom_src column to get spatialite to recognize it
SELECT err_if(RecoverGeometryColumn('addresses', 'geom_src', 2229, 'MULTIPOLYGON', 'XY') = 0, 'RecoverGeometryColumn(): validation failed');

-- Use the parcel geometry to create a point for the main geom_4326 column
SELECT AddGeometryColumn('addresses', 'geom_4326', 4326, 'POINT', 'XY');
UPDATE addresses SET geom_4326 = ST_Transform(ST_Centroid(geom_src), 4326);

-- Create and populate the geom_3857 column
SELECT AddGeometryColumn('addresses', 'geom_3857', 3857, 'POINT', 'XY');
UPDATE addresses SET geom_3857 = ST_Transform(geom_4326, 3857);

-- Create and populate the geohash column
ALTER TABLE addresses ADD COLUMN geohash text;
UPDATE addresses SET geohash = ST_GeoHash(geom_4326, 20);
