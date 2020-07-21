.import /var/uhppoted/events/405419896.log raw
.import /var/uhppoted/events/303986753.log raw

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

.headers on
.mode tabs
.output /tmp/events.tsv

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

