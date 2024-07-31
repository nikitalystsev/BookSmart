CREATE SCHEMA IF NOT EXISTS bs;

CREATE TABLE IF NOT EXISTS bs.book
(
    id              UUID PRIMARY KEY NOT NULL,
    title           TEXT             NOT NULL,
    author          TEXT             NOT NULL,
    publisher       TEXT             NOT NULL,
    copies_number   INT              NOT NULL CHECK (copies_number > 0),
    rarity          TEXT             NOT NULL,
    genre           TEXT             NOT NULL,
    publishing_year INT              NOT NULL,
    language        TEXT             NOT NULL,
    age_limit       INT              NOT NULL CHECK (age_limit >= 0)
);

CREATE TABLE IF NOT EXISTS bs.reader
(
    id           UUID PRIMARY KEY NOT NULL,
    fio          TEXT             NOT NULL,
    phone_number VARCHAR(20)      NOT NULL UNIQUE,
    age          INT              NOT NULL CHECK (age > 0 AND age < 100),
    password     TEXT             NOT NULL
);

CREATE TABLE IF NOT EXISTS bs.lib_card
(
    id            UUID PRIMARY KEY NOT NULL,
    reader_id     UUID             NOT NULL,
    lib_card_num  VARCHAR(13)      NOT NULL UNIQUE,
    validity      INT              NOT NULL,
    issue_date    DATE             NOT NULL,
    action_status BOOLEAN          NOT NULL,
    FOREIGN KEY (reader_id) REFERENCES reader (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS bs.favorite_books
(
    book_id   UUID,
    reader_id UUID,
    PRIMARY KEY (book_id, reader_id),
    FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (reader_id) REFERENCES reader (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS bs.reservation
(
    id          UUID PRIMARY KEY NOT NULL,
    reader_id   UUID             NOT NULL,
    book_id     UUID             NOT NULL,
    issue_date  DATE             NOT NULL,
    return_date DATE             NOT NULL,
    state       TEXT             NOT NULL,
    FOREIGN KEY (reader_id) REFERENCES reader (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CHECK (issue_date < return_date)
);
