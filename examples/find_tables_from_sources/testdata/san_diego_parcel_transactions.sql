/*
	Description: San Diego County Parcels Transactions Transform Script
	Table names:
	- parcel_transactions_data: the incoming assessor data.
	- parcel_transactions: the final transaction data.
*/

BEGIN;

	-- Create the overlays table
	CREATE TABLE parcel_transactions AS
		SELECT
	        '06'     AS state_fips_code,   -- California
    	    '073'    AS county_fips_code,  -- San Diego
			(
				TRIM(h1) || '_' ||
				COALESCE(TRIM(h37), '') || '_' || -- doc type 
				COALESCE(TRIM(h38),   '') || '_' || -- doc number
				COALESCE(printf('%06d', h39), '') || '_' || -- doc date
				COALESCE(printf('%06d', h63), '') -- transaction date
			) AS reference_id,

			-- assessor parcel number
			h1 AS parcel_apn,

			-- transaction date parsing
			CASE
				WHEN h63 IS NULL OR h63 = 0 THEN NULL
				ELSE date(
				(CASE
					WHEN CAST(substr(printf('%06d', h63), 5, 2) AS INT) >= 70
					THEN '19' || substr(printf('%06d', h63), 5, 2)
					ELSE '20' || substr(printf('%06d', h63), 5, 2)
				END)
				|| '-' || substr(printf('%06d', h63), 1, 2)  -- MM
				|| '-' || substr(printf('%06d', h63), 3, 2)  -- DD
				)
			END AS transaction_date,

			-- document identifiers
			h38 AS document_number,

            -- document date parsing
			CASE
				WHEN h39 IS NULL OR h39 = 0 
                    THEN NULL
				ELSE 
                    date(
                        (CASE
                            WHEN CAST(substr(printf('%06d', h39), 5, 2) AS INT) >= 70
                            THEN '19' || substr(printf('%06d', h39), 5, 2)
                            ELSE '20' || substr(printf('%06d', h39), 5, 2)
                        END)
                        || '-' || substr(printf('%06d', h39), 1, 2)  -- MM
                        || '-' || substr(printf('%06d', h39), 3, 2)  -- DD
                    )
			END AS recording_date,

			-- document type
			CASE
				WHEN h37 = 0 THEN 'unresearched'
				WHEN h37 = 1 THEN 'regular owner change'
				WHEN h37 = 2 THEN 'quit claim'
				WHEN h37 = 3 THEN 'unrecorded deed'
				WHEN h37 = 4 THEN 'death certificate'
				WHEN h37 = 5 THEN 'unrecorded death certificate'
				WHEN h37 = 6 THEN 'other'
				WHEN h37 = 7 THEN 'unknown or multiple documents'
				WHEN h37 = 8 THEN 'recorded contract'
			END AS document_type,

			-- transaction amounts
			null AS amount_cents,
			null AS transfer_fee_cents

		FROM parcel_transactions_data;

    -- Nullify dates that SQLite "rolled" forward (e.g. Feb 29 1991 -> March 1)
    UPDATE parcel_transactions
    SET 
        recording_date = CASE 
            WHEN date(recording_date) IS NOT NULL AND date(recording_date, '+0 days') != recording_date 
            	THEN NULL
            ELSE recording_date
        END,
        transaction_date = CASE 
            WHEN date(transaction_date) IS NOT NULL AND date(transaction_date, '+0 days') != transaction_date 
            	THEN NULL 
            ELSE transaction_date 
        END
    WHERE 
        recording_date IS NOT NULL 
        OR transaction_date IS NOT NULL
    ;

	-- remove duplicates.
    DELETE FROM 
		parcel_transactions
	WHERE rowid NOT IN (
		SELECT MIN(rowid)
		FROM parcel_transactions
		GROUP BY reference_id
	);

	-- add a unique index to prove uniqueness
	CREATE UNIQUE INDEX idx_parcel_transactions_reference_id 
		ON parcel_transactions(reference_id);

	-- build parcel transaction owners table
	CREATE TABLE parcel_transaction_owners AS
		WITH base AS (
			SELECT
				pt.reference_id AS parcel_transaction_id,
				CAST(owners.key AS INTEGER) + 1 AS owner_ordinal,

				json_extract(owners.value, '$.name')      AS owner_name,
				json_extract(owners.value, '$.is_entity') AS is_entity,

				CASE
					WHEN json_extract(owners.value, '$.interest') IS NULL THEN NULL
					WHEN TRIM(json_extract(owners.value, '$.interest')) = '' THEN NULL
					ELSE CAST(REPLACE(TRIM(json_extract(owners.value, '$.interest')), '%', '') AS REAL)
				END AS percent_interest,

				NULL AS ownership_code,

				parse_san_diego_assessor_mailing_address(src.h12) AS addr_json,

				src.h13 AS postal

			FROM parcel_transactions pt
			JOIN parcel_transactions_data src
				ON (
					TRIM(src.h1) || '_' ||
					COALESCE(TRIM(src.h37), '') || '_' ||
					COALESCE(TRIM(src.h38), '') || '_' ||
					COALESCE(CASE WHEN src.h39 IS NULL OR src.h39 = 0 THEN '' ELSE printf('%06d', src.h39) END, '') || '_' ||
					COALESCE(CASE WHEN src.h63 IS NULL OR src.h63 = 0 THEN '' ELSE printf('%06d', src.h63) END, '')
				) = pt.reference_id
			CROSS JOIN json_each(
				parse_san_diego_assessor_owners(src.h10)
			) AS owners
		)
		SELECT
			parcel_transaction_id,
			owner_ordinal,
			owner_name,
			is_entity,
			percent_interest,
			ownership_code,
			addr_json,

			CASE WHEN owner_ordinal = 1 THEN addr_json->>'$.lines[0]' END AS address_line1,
			CASE WHEN owner_ordinal = 1 THEN addr_json->>'$.lines[1]' END AS address_line2,
			CASE WHEN owner_ordinal = 1 THEN addr_json->>'$.lines[2]' END AS address_line3,
			CASE WHEN owner_ordinal = 1 THEN addr_json->>'$.lines[3]' END AS address_line4,

			CASE
				WHEN owner_ordinal = 1 THEN 
					postal
				ELSE 
					NULL
			END AS postal
		FROM base
	;

COMMIT;
