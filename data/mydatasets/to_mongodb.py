import csv
from uuid import UUID, uuid4

from bson import Binary, Int64
from pymongo import MongoClient

mongo_uri = "mongodb://nikitalystsev:zhpiix69@localhost:27017/?directConnection=true&authSource=booksmart"

client = MongoClient(mongo_uri)

db = client['booksmart']
collection = db['book']
collection_reader = db['reader']
csv_file_path = 'books.csv'

file = open(csv_file_path, newline='', encoding="utf8", errors="ignore")
reader = csv.DictReader(file)


def convert_to_bson_binary(id_str):
    try:
        uuid = UUID(id_str)
        return Binary(uuid.bytes, subtype=4)
    except ValueError:
        print(f"Warning: Invalid UUID string {id_str}. Skipping this entry.")
        return None


def convert_to_int64(data_str):
    try:
        data_int = int(data_str)
        return Int64(data_int)
    except ValueError:
        print(f"Warning: Invalid convert to int64 {data_str}. Skipping this entry.")
        return None


# Список для хранения документов
books = []

for row in reader:
    book = {
        "_id": convert_to_bson_binary(row['id']),
        "title": row['title'],
        "author": row['author'],
        "publisher": row['publisher'],
        "copies_number": convert_to_int64(row['copiesNumber']),
        "rarity": row['rarity'],
        "genre": row['genres'],
        "publishing_year": convert_to_int64(row['publishDate']),
        "language": row['language'],
        "age_limit": convert_to_int64(row['ageLimit']),
    }
    if book["_id"] is not None:  # Проверяем, что _id валиден
        books.append(book)

# Вставляем все книги за один раз
if books:
    collection.insert_many(books)

# Создание читателя
reader_id = uuid4()
reader = {
    "_id": Binary(reader_id.bytes, subtype=4),
    "fio": "Лысцев Никита Дмитриевич",
    "phone_number": "89314022581",
    "age": convert_to_int64("21"),
    "password": "$2a$10$xDzRFS0ClhEcosyFVQEPCev8AXakZyYau4Hk8iN3dyTXJYXUj1coO",
    "role": "Admin"
}

collection_reader.insert_one(reader)

file.close()

print("Data successfully migrated to MongoDB!")

client.close()
