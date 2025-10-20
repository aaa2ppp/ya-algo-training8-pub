import sys

def solve1(a, b):
    # a[0] и b[0] не используются — индексация с 1
    if len(a) <= 1 or len(b) <= 1:
        return 0

    max_a = max(a)
    max_b = max(b)

    # sum_i[x] = сумма индексов i, где a[i] == x
    # cnt_i[x] = количество таких i
    sum_a = [0] * (max_a + 1)
    cnt_a = [0] * (max_a + 1)

    sum_b = [0] * (max_b + 1)
    cnt_b = [0] * (max_b + 1)

    for i in range(1, len(a)):
        val = a[i]
        sum_a[val] += i
        cnt_a[val] += 1

    for j in range(1, len(b)):
        val = b[j]
        sum_b[val] += j
        cnt_b[val] += 1

    ans = 0
    # Перебираем только возможные значения из диапазонов
    for ai in range(len(sum_a)):
        if cnt_a[ai] == 0:
            continue
        for bj in range(len(sum_b)):
            if cnt_b[bj] == 0:
                continue
            ans += (sum_a[ai] * cnt_b[bj] - sum_b[bj] * cnt_a[ai]) * abs(ai - bj)

    return ans

def main():
    data = sys.stdin.read().split()
    if not data:
        return

    it = iter(data)
    n = int(next(it))
    a = [0] + [int(next(it)) for _ in range(n)]

    m = int(next(it))
    b = [0] + [int(next(it)) for _ in range(m)]

    result = solve1(a, b)
    print(result)

if __name__ == "__main__":
    main()
