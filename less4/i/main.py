import sys
from typing import List, Tuple
from itertools import chain
from collections import defaultdict


def calc_squares(d: int) -> List[int]:
    """Возвращает список всех квадратов (начиная с 0), <= d, в порядке возрастания."""
    q = []
    i = 0
    while i * i <= d:
        q.append(i * i)
        i += 1
    return q


def calc_square_sums(d: int) -> List[Tuple[int, int]]:
    """
    Возвращает список пар (a, b), таких что a^2 + b^2 == d и a <= b.
    Здесь a и b — целые неотрицательные числа (корни, не квадраты!).
    """
    q = calc_squares(d)
    qs = []
    i, j = 0, len(q) - 1
    while i <= j:
        c = q[i] + q[j]
        if c < d:
            i += 1
        elif c > d:
            j -= 1
        else:
            # q[i] = i*i, q[j] = j*j → пара (i, j)
            qs.append((i, j))
            j -= 1
    return qs


def calc_offsets(d: int) -> List[Tuple[int, int]]:
    """Возвращает все целочисленные смещения (dx, dy), такие что dx^2 + dy^2 == d."""
    pairs = calc_square_sums(d)
    offsets = []
    for a, b in pairs:
        if a == 0:
            # (0, ±b), (±b, 0)
            offsets.extend([(0, b), (0, -b), (b, 0), (-b, 0)])
        elif a == b:
            # (±a, ±a) — 4 варианта
            offsets.extend([(a, a), (a, -a), (-a, -a), (-a, a)])
        else:
            # Все комбинации знаков и перестановок — 8 штук
            offsets.extend([
                (a, b), (a, -b), (-a, -b), (-a, b),
                (b, a), (b, -a), (-b, -a), (-b, a)
            ])
    return offsets


def solve(d, forest) -> int:
    """
    Считает количество неориентированных пар деревьев, расстояние между которыми в квадрате равно d.
    Поддерживает дубликаты (много деревьев в одной точке).
    """
    offsets = calc_offsets(d)
    if not offsets:
        return 0

    count = 0
    for x, y in forest.keys():
        for dx, dy in offsets:
            nx, ny = x + dx, y + dy
            if (nx, ny) in forest:
                count += forest[(nx, ny)]

    return count // 2


def main() -> None:
    it = chain.from_iterable(line.split() for line in sys.stdin.buffer)

    n = int(next(it))
    d = int(next(it))

    forest = defaultdict(int)
    for _ in range(n):
        x = int(next(it))
        y = int(next(it))
        # Считаем количество деревьев в каждой точке
        forest[(x, y)] += 1

    ans = solve(d, forest)
    print(ans)


if __name__ == "__main__":
    main()
