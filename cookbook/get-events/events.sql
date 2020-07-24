-- IMPORT EVENT FROM EVENT FILES

.import /tmp/uhppoted-405419896.events raw
.import /tmp/uhppoted-303986753.events raw

INSERT OR IGNORE INTO events 
       SELECT TRIM(SUBSTR(event,1,11))  AS deviceID,
              TRIM(SUBSTR(event,12,7))  AS eventID,
              TRIM(SUBSTR(event,19,20)) AS timestamp,
              TRIM(SUBSTR(event,39,13)) AS card,
              TRIM(SUBSTR(event,52,2))  AS door,
              UPPER(TRIM(SUBSTR(event,54,6))) = 'TRUE' AS granted,
              TRIM(SUBSTR(event,60,4))  AS result
              FROM raw;

DELETE FROM raw;

-- GENERATE EVENT SUMMARY

.headers on
.mode tabs
.once /tmp/events.tsv

SELECT p.day AS Date,Total,Granted,Refused FROM
       ( SELECT DATE(timestamp) AS day,COUNT(*) AS Total
                FROM events
                GROUP BY day
       ) AS p
       LEFT JOIN ( SELECT DATE(timestamp) AS day,COUNT(*) AS Granted
                          FROM events
                          WHERE granted=1
                          GROUP BY day
                 ) AS q
       ON p.day=q.day
       LEFT JOIN ( SELECT DATE(timestamp) AS day,COUNT(*) AS Refused
                          FROM events
                          WHERE granted=0
                          GROUP BY day
                 ) AS r
       ON p.day=r.day
       ORDER BY Date;

-- PRUNE EVENTS OLDER THAN A YEAR

SELECT * FROM events WHERE timestamp < date('NOW','start of month','-12 month');

-- TRUNCATE EVENT FILES

.once /tmp/uhppoted-405419896.events
SELECT * FROM event WHERE FALSE;

.once /tmp/uhppoted-303986753.events
SELECT * FROM event WHERE FALSE;
