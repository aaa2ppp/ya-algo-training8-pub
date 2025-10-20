import sys
from itertools import islice


def solve(a, b):
    n = len(a)

    # добавляем терминальные элементы, чтобы не делать лишние проверки в основном цикле
    a.append(0)
    b.append(0)

    # Бинарный поиск по k от 0 до n
    lo, hi = 0, n
    while lo < hi:
        k = (lo + hi) // 2

        ri, ai = 0, 0
        rv = a[0]
        av = b[0]
        while True:
            if ai > ri:
                if ai == n or ai - ri > k:
                    break
            else:
                if ri == n:
                    break
                if ri - ai > k:
                    ai += 1
                    av = b[ai]
            if rv < av:
                av -= rv
                ri += 1
                rv = a[ri]
            elif rv > av:
                rv -= av
                ai += 1
                av = b[ai]
            else:
                ri += 1
                rv = a[ri]
                ai += 1
                av = b[ai]

        if ri == n:
            hi = k
        else:
            lo = k + 1

    return -1 if lo == n else lo


def main():
    data = sys.stdin.read().split()

    n = int(data[0])
    a = list(map(int, islice(data, 1, n + 1)))
    b = list(map(int, islice(data, n + 1, n + 1 + n)))

    ans = solve(a, b)
    sys.stdout.write(str(ans))


if __name__ == "__main__":
    main()
