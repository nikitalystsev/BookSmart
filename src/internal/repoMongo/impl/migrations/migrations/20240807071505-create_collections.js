module.exports = {
    async up(db, client) {
        await db.createCollection("book", {
            validator: {
                $jsonSchema: {
                    bsonType: "object",
                    required: ["_id", "title", "author", "publisher", "copies_number", "rarity", "genre", "publishing_year", "language", "age_limit"],
                    properties: {
                        _id: {bsonType: "string"},
                        title: {bsonType: "string"},
                        author: {bsonType: "string"},
                        publisher: {bsonType: "string"},
                        copies_number: {bsonType: "int", minimum: 0},
                        rarity: {bsonType: "string"},
                        genre: {bsonType: "string"},
                        publishing_year: {bsonType: "int", minimum: 0},
                        language: {bsonType: "string"},
                        age_limit: {bsonType: "int", minimum: 0}
                    }
                }
            },
            validationLevel: "strict",
            validationAction: "error"
        });
        await db.createCollection("reader", {
            validator: {
                $jsonSchema: {
                    bsonType: "object",
                    required: ["_id", "fio", "phone_number", "age", "password", "role"],
                    properties: {
                        _id: {bsonType: "string"},
                        fio: {bsonType: "string"},
                        phone_number: {bsonType: "string"},
                        age: {bsonType: "int"},
                        password: {bsonType: "string"},
                        role: {bsonType: "string"},
                    }
                }
            },
            validationLevel: "strict",
            validationAction: "error"
        });
        await db.createCollection("lib_card", {
            validator: {
                $jsonSchema: {
                    bsonType: "object",
                    required: ["_id", "reader_id", "lib_card_num", "validity", "issue_date", "action_status"],
                    properties: {
                        _id: {bsonType: "string"},
                        reader_id: {bsonType: "string"},
                        lib_card_num: {bsonType: "string"},
                        validity: {bsonType: "int"},
                        issue_date: {bsonType: "date"},
                        action_status: {bsonType: "bool"},
                    }
                }
            },
            validationLevel: "strict",
            validationAction: "error"
        });
        await db.createCollection("reservation", {
            validator: {
                $jsonSchema: {
                    bsonType: "object",
                    required: ["_id", "reader_id", "book_id", "issue_date", "return_date", "state"],
                    properties: {
                        _id: {bsonType: "string"},
                        reader_id: {bsonType: "string"},
                        book_id: {bsonType: "string"},
                        issue_date: {bsonType: "date"},
                        return_date: {bsonType: "date"},
                        action_status: {bsonType: "string"},
                    }
                }
            },
            validationLevel: "strict",
            validationAction: "error"
        });
        await db.createCollection("favorite_books", {
            validator: {
                $jsonSchema: {
                    bsonType: "object",
                    required: ["_id", "reader_id", "book_id"],
                    properties: {
                        _id: {bsonType: "objectId"},
                        reader_id: {bsonType: "string"},
                        book_id: {bsonType: "string"},
                    }
                }
            },
            validationLevel: "strict",
            validationAction: "error"
        });
    },

    async down(db, client) {
        await db.collection('book').drop();
        await db.collection('lib_card').drop();
        await db.collection('reader').drop();
        await db.collection('reservation').drop();
        await db.collection('favorite_books').drop();
    }
};