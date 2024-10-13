import ast
import csv
import random
import re
import uuid
from datetime import datetime

MIN_COPIES_NUM = 5
MAX_COPIES_NUM = 15


class Generator:
    """
    Класс для генерации и заполнения файлов с данными для базы данных

    https://zenodo.org/records/4265096 -- ссылка на датасет
    """

    def __init__(self):
        pass

    @staticmethod
    def __check_fields(fields_dict: dict):
        """
        Метод проверяет, что у полей есть непустые значения
        """

        for value in fields_dict.values():
            if not value:
                return False

        return True

    @staticmethod
    def __extract_year(date_str: str):
        """
        Метод для извлечения года их даты
        """
        try:
            date_obj = datetime.strptime(date_str, '%m/%d/%y')
        except ValueError:
            try:
                date_str = re.sub(r'(st|nd|rd|th)', '', date_str).strip()
                date_obj = datetime.strptime(date_str, '%B %d %Y')
            except ValueError as e:
                # print(f"Ошибка: не удалось распарсить дату '{date_str}' {e}")
                # time.sleep(5)
                return None

        return date_obj.year

    def generate_books(self, filepath: str):
        """
        Метод позволяет получить датасет с книгами
        """
        csv_file = open(filepath, newline='', encoding="utf8", errors="ignore")
        output_csv_file = open("mydatasets/books.csv", "w", newline='',
                               encoding="utf8", errors="ignore")

        csv_reader = csv.DictReader(csv_file)

        selected_fields = ['title', 'author', 'publisher', 'genres', 'publishDate', 'language']

        book_fields = ['id'] + selected_fields[:3] + ['copiesNumber', 'rarity'] + selected_fields[3:] + ['ageLimit']

        csv_writer = csv.DictWriter(output_csv_file, fieldnames=book_fields)

        cnt = 0
        for row in csv_reader:
            selected_row = {field: row[field] for field in selected_fields if field in row}

            if not self.__check_fields(selected_row):
                continue
            year = self.__extract_year(selected_row['publishDate'])
            if not year:
                continue

            selected_row['genres'] = ','.join(ast.literal_eval(selected_row['genres']))
            if not selected_row['genres']:
                continue

            selected_row['publishDate'] = year

            selected_row['id'] = str(uuid.uuid4())
            selected_row['copiesNumber'] = random.randint(MIN_COPIES_NUM, MAX_COPIES_NUM)
            selected_row['rarity'] = random.choice(['Common', 'Rare', 'Unique'])
            selected_row['ageLimit'] = random.choice([0, 6, 12, 16, 18, 21])

            if cnt == 0:
                print(selected_row)
            csv_writer.writerow(selected_row)

            print(f"Строка №{cnt + 1}")
            cnt += 1

        output_csv_file.close()
        csv_file.close()


def main():
    generator = Generator()
    generator.generate_books('datasets/books.csv')


if __name__ == '__main__':
    main()
