CREATE TABLE events
(
    id                uuid PRIMARY KEY,
    title             varchar(100),
    begin_time        timestamptz NOT NULL,
    end_time          timestamptz NOT NULL,
    description       text,
    user_id           uuid       NOT NULL,
    notification_time timestamptz NOT NULL,

    CONSTRAINT valid_period CHECK ( begin_time < end_time ),
    CONSTRAINT valid_notification_time CHECK ( notification_time <= begin_time )
);