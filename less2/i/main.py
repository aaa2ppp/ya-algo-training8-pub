import sys
from typing import TextIO, BinaryIO, Iterator


def solve(vals, n, m):
    total = n * m
    # Создаём список индексов [0, 1, 2, ..., total-1]
    indices = list(range(total))
    # Сортируем индексы по значению в vals
    indices.sort(key=vals.__getitem__)  # ← ключевой момент!

    visited = [0] * total
    maximum = 0
    dirs = (-1, 0, 1, 0, 0, -1, 0, 1)  # упаковано для скорости

    for idx in indices:
        val = vals[idx]
        i = idx // m
        j = idx % m
        visited[idx] = 1

        # Распаковываем направления без создания кортежей
        for d in range(4):
            di = dirs[d * 2]
            dj = dirs[d * 2 + 1]
            ni = i + di
            nj = j + dj
            if 0 <= ni < n and 0 <= nj < m:
                nidx = ni * m + nj
                if vals[nidx] == val - 1:
                    if visited[nidx] + 1 > visited[idx]:
                        visited[idx] = visited[nidx] + 1
        if visited[idx] > maximum:
            maximum = visited[idx]

    return maximum


def run(input_stream: BinaryIO, output_stream: TextIO) -> None:
    data = input_stream.read().split()
    n = int(data[0])
    m = int(data[1])
    total = n * m
    # Парсим всё в плоский список
    vals = list(map(int, data[2:2 + total]))
    ans = solve(vals, n, m)
    output_stream.write(str(ans))


if __name__ == "__main__":
    run(sys.stdin.buffer, sys.stdout)
