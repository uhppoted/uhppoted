CREATE TABLE raw(event TEXT);

CREATE TABLE events(deviceID  INTEGER,
                    eventID   INTEGER,
                    timestamp TEXT,
                    card      INTEGER,
                    door      INTEGER,
                    granted   INTEGER,
                    result    INTEGER,
                    PRIMARY KEY(deviceID, eventID, timestamp));

CREATE INDEX eventdate ON events(DATE(timestamp));
