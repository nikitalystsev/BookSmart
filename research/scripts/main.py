import parser


def main() -> None:
    """
    Главная функция
    """
    parser.parse_stats_history(
        '../locust_stats/with_balance_stats_history.csv',
        '../data/with_balance.csv',
        "../md/with_balance.txt"
    )

    parser.parse_stats_history(
        '../locust_stats/without_balance_stats_history.csv',
        '../data/without_balance.csv',
        "../md/without_balance.txt"
    )

    parser.build_comparative_graphic(
        '../data/with_balance.csv',
        '../data/without_balance.csv',
        "../graphics/compare_graph.svg"
    )


if __name__ == '__main__':
    main()
