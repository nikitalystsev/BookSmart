import csv

import matplotlib.pyplot as plt


def parse_stats_history(filepath: str, out_filepath: str, md_out_path: str):
    """
    Функция парсит историю статистики исследования
    """
    stats_file = open(file=filepath, newline='', encoding='utf-8')
    out_file = open(file=out_filepath, mode='w', newline='', encoding='utf-8')
    md_file = open(file=md_out_path, mode='w', newline='', encoding='utf-8')

    reader = csv.DictReader(stats_file)

    out_fieldnames = [
        'User Count',
        'Total Request Count',
        'Total Average Response Time',
        'Requests/s',
    ]

    writer = csv.DictWriter(out_file, fieldnames=out_fieldnames)
    writer.writeheader()

    for row in reader:
        # if row['Name'] == 'Aggregated':
        #     continue

        writer.writerow({
            'User Count': row['User Count'],
            'Total Request Count': row['Total Request Count'],
            'Total Average Response Time': round(float(row['Total Average Response Time']), 3),
            'Requests/s': round(float(row['Requests/s']))
        })

        md_file.write(
            f"{round(float(row['Requests/s']))} {float(row['Total Average Response Time']):.3f}\n"
        )

    md_file.close()
    out_file.close()
    stats_file.close()


def process_data_for_graphic(filepath: str):
    """
    Функция для получения данных для сравнительного графика
    """
    file = open(file=filepath, newline='', encoding='utf-8')
    file_reader = csv.DictReader(file)

    required_data = [[], []]

    for row in file_reader:
        required_data[0].append(float(row['Total Average Response Time']))
        required_data[1].append(float(row['Requests/s']))

    file.close()

    return required_data


def build_comparative_graphic(
        with_balance_filepath: str,
        without_balance_filepath: str,
        output_svg: str
):
    """
    Функция строит графики исследований
    """
    with_balance_data = process_data_for_graphic(with_balance_filepath)
    without_balance_data = process_data_for_graphic(without_balance_filepath)

    plt.plot(with_balance_data[1], with_balance_data[0], label='C исп-м балансировки', color='blue', marker='x')
    plt.plot(without_balance_data[1], without_balance_data[0], label='Без исп-я балансировки', color='red', marker='*')

    plt.ylabel('Среднее время ответа, мс')
    plt.xlabel('Число запросов в секунду')

    plt.legend()
    plt.grid()
    # plt.show()

    # Устанавливаем одинаковый масштаб для осей
    plt.axis('equal')

    plt.savefig(output_svg, format='svg')

# locust --host=http://localhost --headless --csv=../locust_stats/with_balance -u 500 -r 10 -t 1m
