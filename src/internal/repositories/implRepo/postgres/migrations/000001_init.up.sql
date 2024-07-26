CREATE TABLE Book
(
    id             UUID PRIMARY KEY NOT NULL,
    title          TEXT             NOT NULL,
    author         TEXT             NOT NULL,
    publisher      TEXT             NOT NULL,
    copiesNumber   INT              NOT NULL,
    rarity         TEXT             NOT NULL,
    genre          TEXT             NOT NULL,
    publishingYear INT              NOT NULL,
    language       TEXT             NOT NULL,
    ageLimit       INT              NOT NULL
);

CREATE TABLE Reader
(
    id          UUID PRIMARY KEY NOT NULL,
    fio         TEXT             NOT NULL,
    phoneNumber VARCHAR(20)      NOT NULL UNIQUE,
    age         INT              NOT NULL,
    password    VARCHAR(10)      NOT NULL
);

CREATE TABLE LibCard
(
    id           UUID PRIMARY KEY NOT NULL,
    readerID     UUID             NOT NULL,
    libCardNum   VARCHAR(13)      NOT NULL UNIQUE,
    validity     INT              NOT NULL,
    issueDate    DATE             NOT NULL,
    actionStatus BOOLEAN          NOT NULL,
    FOREIGN KEY (readerID) REFERENCES Reader (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE FavoriteBooks
(
    bookID   UUID,
    readerID UUID,
    PRIMARY KEY (bookID, readerID),
    FOREIGN KEY (bookID) REFERENCES Book (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (readerID) REFERENCES Reader (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE Reservation
(
    id         UUID PRIMARY KEY NOT NULL,
    readerID   UUID             NOT NULL,
    bookID     UUID             NOT NULL,
    issueDate  DATE             NOT NULL,
    returnDate DATE             NOT NULL,
    state      TEXT             NOT NULL,
    FOREIGN KEY (readerID) REFERENCES Reader (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (bookID) REFERENCES Book (id) ON DELETE CASCADE ON UPDATE CASCADE
);
