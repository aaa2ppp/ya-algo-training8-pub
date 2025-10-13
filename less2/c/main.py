import sys

# Типы событий: END раньше BEGIN при одинаковом времени
EVENT_END = 0
EVENT_BEGIN = 1

def solve(intervals):
    n = len(intervals)
    if n == 0:
        return 0.0

    events = [None] * (2 * n)
    for i, (b, e, w) in enumerate(intervals):
        events[2 * i] = (b, EVENT_BEGIN, i)
        events[2 * i + 1] = (e, EVENT_END, i)

    events.sort()

    cum_weight = [0.0] * n
    max_weight = 0.0
    prev_max_weight = 0.0

    for _, typ, idx in events:
        if typ == EVENT_BEGIN:
            w_val = prev_max_weight + intervals[idx][2]
            cum_weight[idx] = w_val
            if w_val > max_weight:
                max_weight = w_val
        else:  # EVENT_END
            if cum_weight[idx] > prev_max_weight:
                prev_max_weight = cum_weight[idx]

    return max_weight

def main():
    data = sys.stdin.read().split()
    if not data:
        print("0")
        return

    n = int(data[0])
    intervals = []
    idx = 1
    for _ in range(n):
        b = float(data[idx])
        e = float(data[idx + 1])
        w = float(data[idx + 2])
        idx += 3
        intervals.append((b, e, w))

    result = solve(intervals)
    sys.stdout.write(f"{result:g}\n")

if __name__ == "__main__":
    main()
