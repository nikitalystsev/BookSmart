import csv
import random

from locust import HttpUser, task, constant


def read_id_column(file_path):
    ids = []
    with open(file_path, mode='r', newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            ids.append(row['id'])
    return ids


class BookSmartTestUser(HttpUser):
    wait_time = constant(0)

    book_ids = read_id_column("../../data/mydatasets/books.csv")

    name = "test_requests"

    @task
    def get_page_books(self):
        """
        Получить страницу книг
        """
        self.client.get(f"/api/v1/books?page_number=1", name=self.name)

    @task
    def get_book_by_id(self):
        """
        Получить книгу по id
        """
        book_id = random.choice(self.book_ids)
        self.client.get(f"/api/v1/books/{book_id}", name=self.name)

    @task
    def get_rating_book_by_id(self):
        """
        Получить рейтинг книги по id
        """
        book_id = random.choice(self.book_ids)
        self.client.get(f"/api/v1/books/{book_id}/ratings", name=self.name)

# locust --host=http://localhost --headless --csv=../locust_stats/with_balance -u 500 -r 10 -t 1m
