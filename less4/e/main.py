import sys
from itertools import chain

def solve(a, trips, k):
    # a is 1-indexed; a[0] is dummy
    n = len(a) - 1

    # Prefix sum of a
    sumA = [0] * (n + 1)
    for i in range(1, n + 1):
        sumA[i] = sumA[i - 1] + a[i]

    # Frequency difference array for trip coverage
    sumT = [0] * (n + 2)  # index up to n+1
    total = 0

    for trip in trips:
        l, r = trip
        sumT[l] += 1
        if r + 1 <= n:
            sumT[r + 1] -= 1
        total += sumA[r] - sumA[l - 1]

    # Convert diff array to prefix sum (actual coverage count per position)
    for i in range(1, n + 1):
        sumT[i] += sumT[i - 1]

    # Create list of positions 1..n
    positions = list(range(1, n + 1))

    # Sort positions by coverage count descending (most covered first)
    positions.sort(key=lambda i: sumT[i], reverse=True)

    # Greedily remove as much as possible from most covered positions
    for i in positions:
        if k <= 0:
            break
        if k >= a[i]:
            total -= sumT[i] * a[i]
            k -= a[i]
        else:
            total -= sumT[i] * k
            k = 0

    return total


def main():
    it = chain.from_iterable(line.split() for line in sys.stdin.buffer)

    n = int(next(it))
    m = int(next(it))
    k = int(next(it))

    # Read a[1..n]; make a 1-indexed (a[0] unused)
    a = [0] * (n+1)
    for i in range(1, n+1):
        a[i] = int(next(it))

    trips = [None] * m
    for i in range(m):
        l = int(next(it))
        r = int(next(it))
        trips[i] = (l, r)

    ans = solve(a, trips, k)
    print(ans)


if __name__ == "__main__":
    main()
