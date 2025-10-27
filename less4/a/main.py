import sys


def solve(front, back):
    Arrival, Departure = 0, 1  # Порядок важен! для сортировки

    A, B = 0, 1  # points
    points = [0, 0]

    events = []

    for dep, arr in front:
        events.append((dep, Departure, A))
        events.append((arr, Arrival, B))

    for dep, arr in back:
        events.append((dep, Departure, B))
        events.append((arr, Arrival, A))

    events.sort()

    count = 0
    for _, etype, point in events:
        if etype == Arrival:
            points[point] += 1
        else:  # Departure
            if points[point] > 0:
                points[point] -= 1
            else:
                count += 1

    return count


def main():
    data = sys.stdin.read().split()

    it = iter(data)
    n = int(next(it))
    front = []
    for _ in range(n):
        s = next(it)
        dep, arr = s.split('-')
        front.append((dep, arr))

    m = int(next(it))
    back = []
    for _ in range(m):
        s = next(it)
        dep, arr = s.split('-')
        back.append((dep, arr))

    ans = solve(front, back)
    print(ans)


if __name__ == "__main__":
    main()
