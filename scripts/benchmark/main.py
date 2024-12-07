#! /bin/python

import typing
import matplotlib.pyplot as plt
import os.path
import time


class record:
    timestamp: float
    records: dict[str, float]

    def __init__(self, ts: float, records: dict[str, float]):
        self.timestamp = ts
        self.records = records


def parse_time(stime: str) -> float:
    return get_time(stime.split(".")[0].strip())


def get_time(stime: str) -> float:
    t = time.strptime(stime, "%Y-%m-%d %H:%M:%S")
    return time.mktime(t)


def parse_value(value: str) -> tuple[str, float]:
    items = value.split()

    if 0 == len(items):
        raise Exception("Empty value")
    elif 1 == len(items):
        return "N", float(value.strip())
    elif 2 == len(items):
        return items[1].strip(), float(items[0].strip())
    else:
        raise Exception(f"Broken value: {value}")


def parse_results(filename: str) -> tuple[str, list[record]]:
    records: list[record] = []
    name, _ = os.path.splitext(os.path.basename(filename))

    file = open(filename, "r")
    lines = file.readlines()
    file.close()

    start = parse_time(lines[0])

    for line in lines[1:-2]:
        t, values = line.split(": ")
        records.append(record(
            parse_time(t) - start,
            dict(map(
                lambda x: parse_value(x.strip()),
                values.split('\t')
            ))
        ))

    return name, records


def extract(
    records: list[record],
    f: typing.Callable[[dict[str, float]], float]
) -> tuple[list[float], list[float]]:
    timeline: list[float] = [0.0] * len(records)
    values: list[float] = [0.0] * len(records)

    for i, rec in enumerate(records):
        timeline[i] = i + 1
        values[i] = f(rec.records)

    return timeline, values


def main() -> int:
    _, sqlx = parse_results("../../docker/benchmark/result/sqlx.txt")
    _, gorm = parse_results("../../docker/benchmark/result/gorm.txt")

    _, axes = plt.subplots(2, 2)

    axes[0, 0].plot(*extract(sqlx, lambda x: x[".B/op"] / 2 ** 20), label="sqlx")
    axes[0, 0].plot(*extract(gorm, lambda x: x[".B/op"] / 2 ** 20), label="gorm")
    axes[0, 0].set_xlabel("Замер")
    axes[0, 0].set_ylabel("Память (МБ)")
    axes[0, 0].set_title("Память на сценарий")
    axes[0, 0].grid(True)
    axes[0, 0].legend()

    axes[0, 1].plot(*extract(sqlx, lambda x: x[".allocs/op"]), label="sqlx")
    axes[0, 1].plot(*extract(gorm, lambda x: x[".allocs/op"]), label="gorm")
    axes[0, 1].set_xlabel("Замер")
    axes[0, 1].set_ylabel("Число выделений памяти (шт)")
    axes[0, 1].set_title("Число выделений памяти на сценарий")
    axes[0, 1].grid(True)
    axes[0, 1].legend()

    axes[1, 0].plot(*extract(sqlx, lambda x: x["callns/op"] * 1e-3), label="sqlx")
    axes[1, 0].plot(*extract(gorm, lambda x: x["callns/op"] * 1e-3), label="gorm")
    axes[1, 0].set_xlabel("Замер")
    axes[1, 0].set_ylabel("Время построения запроса (мкс)")
    axes[1, 0].set_title("Время построения запроса")
    axes[1, 0].grid(True)
    axes[1, 0].legend()

    axes[1, 1].plot(*extract(sqlx, lambda x: x["serns/op"] * 1e-6), label="sqlx")
    axes[1, 1].plot(*extract(gorm, lambda x: x["serns/op"] * 1e-6), label="gorm")
    axes[1, 1].set_xlabel("Замер")
    axes[1, 1].set_ylabel("Время сериализации (мс)")
    axes[1, 1].set_title("Время сериализации")
    axes[1, 1].grid(True)
    axes[1, 1].legend()

    plt.show()

    return 0


if __name__ == "__main__":
    exit(main())

