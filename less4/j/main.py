import sys
from collections import defaultdict
from typing import List, Tuple
from math import gcd
from itertools import chain


def norm_frac(a: int, b: int) -> Tuple[int, int]:
    if b == 0:
        raise ValueError("Denominator is zero")
    if b < 0:
        a, b = -a, -b
    g = gcd(a, b)
    return (a // g, b // g)


def solve(l: int, w: int, cars: List[Tuple[int, int, int, int]]) -> List[int]:
    n = len(cars)
    # (time, event_type, point), event_type: 0 = Collision, 1 = Finish
    events = []
    collision = defaultdict(list)  # point -> list of car indices
    finish = defaultdict(list)     # time  -> list of car indices
    eliminated = [False] * n

    def add_collision(t, x, y, *car_indices):
        if t[0] <= 0:  # t <= 0
            return
        point = (t, x, y)
        events.append((t, 0, point))
        collision[point].extend(car_indices)

    def add_finish(t, car_index):
        if t[0] <= 0:  # t <= 0
            return
        events.append((t, 1, t))
        finish[t].append(car_index)

    # Pairwise collisions
    for i in range(n):
        x1, y1, vx1, vy1 = cars[i]
        for j in range(i + 1, n):
            x2, y2, vx2, vy2 = cars[j]

            # if vx1 * vy2 == vy1 * vx2:
            #     continue

            dx = x1 - x2
            dy = y1 - y2
            dvx = vx1 - vx2
            dvy = vy1 - vy2

            if dvx == 0 and dvy == 0:
                continue

            if dx * dvy != dy * dvx:
                continue

            # Compute t
            if vx1 != vx2:
                t = norm_frac(dx, vx2 - vx1)
            else:
                t = norm_frac(dy, vy2 - vy1)

            if t[0] <= 0:
                continue

            # x = x1 + t * vx1
            x = norm_frac(x1*t[1] + t[0]*vx1, t[0])
            y = norm_frac(y1*t[1] + t[0]*vy1, t[0])
            add_collision(t, x, y, i, j)

    # Borders and finish
    for i, (x0, y0, vx, vy) in enumerate(cars):
        # Bottom (y=0) or top (y=w)
        if vy != 0:
            t = norm_frac(-y0, vy)
            if t[0] <= 0:
                t = norm_frac(w - y0, vy)
            x = norm_frac(x0*t[1] + t[0]*vx, t[1])
            y = norm_frac(y0*t[1] + t[0]*vy, t[1])
            add_collision(t, x, y, i)

        # Finish line x = l
        if vx != 0:
            t = norm_frac(l - x0, vx)
            if t[0] > 0:
                add_finish(t, i)

    # Sort events: by time (exact comparison), then Collision (0) before Finish (1)
    def event_key(e):
        t, etype, _ = e
        return (t[0] / t[1], etype)  # convert time to float for sorting

    events.sort(key=event_key)

    # Process events
    for _, etype, key in events:
        if etype == 0:  # Collision
            participants = [c for c in collision[key] if not eliminated[c]]
            t, x, y = key
            on_border = (y[0] == 0) or (y[0] == w * y[1])
            if len(participants) > 1 or on_border:
                for c in participants:
                    eliminated[c] = True
        else:  # Finish
            candidates = [c + 1 for c in finish[key] if not eliminated[c]]
            if candidates:
                return candidates

    return []


def main():
    it = chain.from_iterable(line.split() for line in sys.stdin.buffer)

    n = int(next(it))
    l = int(next(it))
    w = int(next(it))

    cars = []
    for _ in range(n):
        x = int(next(it))
        y = int(next(it))
        vx = int(next(it))
        vy = int(next(it))
        cars.append((x, y, vx, vy))

    result = solve(l, w, cars)
    sys.stdout.write(str(len(result)) + "\n")
    if result:
        sys.stdout.write(" ".join(map(str, result)) + "\n")


if __name__ == "__main__":
    main()
