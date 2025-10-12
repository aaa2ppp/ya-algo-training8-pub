import sys
import os
from collections import deque

DEBUG = os.getenv("DEBUG") is not None

# Флаги состояния
GROCERIES_IN_HANDS = 1 << 0
PARCEL_IN_HANDS     = 1 << 1
GROCERIES_DELIVERED = 1 << 2
PARCEL_DELIVERED    = 1 << 3

CARRYING_MASK   = GROCERIES_IN_HANDS | PARCEL_IN_HANDS
DELIVERED_MASK  = GROCERIES_DELIVERED | PARCEL_DELIVERED

# Места:
# 0 — дом
# 1 — супермаркет
# 2 — пункт выдачи

def solve(a, b, c, v):
    # visited[state] = (hops, time)
    # state = (place, flags)
    visited = {}
    queue = deque()

    def remember_state(place, flags, hops, time):
        state = (place, flags)
        if state not in visited or time < visited[state][1]:
            visited[state] = (hops, time)
            queue.append(state)

    def go_to_next_place(place, flags, hops, time, next_hops):
        # Определяем скорость в зависимости от того, что несём
        carrying = flags & CARRYING_MASK
        if carrying == 0:
            speed = v[0]
        elif carrying == CARRYING_MASK:
            speed = v[2]
        else:  # ровно один предмет в руках
            speed = v[1]

        for next_place, dist in next_hops:
            new_time = time + dist / speed
            remember_state(next_place, flags, hops + 1, new_time)

    INF = 100500.0
    ans = INF

    # Начальное состояние: дома, ничего не несём, ничего не доставлено
    start_state = (0, 0)
    visited[start_state] = (0, 0.0)
    queue.append(start_state)

    while queue:
        place, flags = queue.popleft()
        hops, time = visited[(place, flags)]

        if DEBUG:
            print(f"vis: hops={hops}, time={time:.6f}, place={place}, flags={flags:04b}", file=sys.stderr)

        if place == 0:  # дом
            if (flags & DELIVERED_MASK) == DELIVERED_MASK:
                if DEBUG:
                    print(f"bingo! time={time:.6f}", file=sys.stderr)
                ans = min(ans, time)

            # Доставка продуктов
            if flags & GROCERIES_IN_HANDS:
                new_flags = (flags & ~GROCERIES_IN_HANDS) | GROCERIES_DELIVERED
                remember_state(0, new_flags, hops + 1, time)

            # Доставка посылки
            if flags & PARCEL_IN_HANDS:
                new_flags = (flags & ~PARCEL_IN_HANDS) | PARCEL_DELIVERED
                remember_state(0, new_flags, hops + 1, time)

            # Уход из дома
            go_to_next_place(place, flags, hops, time, [(1, a), (2, b)])

        elif place == 1:  # супермаркет
            # Взять продукты, если ещё не взяты и не доставлены
            if not (flags & (GROCERIES_DELIVERED | GROCERIES_IN_HANDS)):
                new_flags = flags | GROCERIES_IN_HANDS
                remember_state(1, new_flags, hops + 1, time)

            # Уйти из супермаркета
            go_to_next_place(place, flags, hops, time, [(0, a), (2, c)])

        elif place == 2:  # пункт выдачи
            # Взять посылку, если ещё не взята и не доставлена
            if not (flags & (PARCEL_DELIVERED | PARCEL_IN_HANDS)):
                new_flags = flags | PARCEL_IN_HANDS
                remember_state(2, new_flags, hops + 1, time)

            # Уйти из пункта выдачи
            go_to_next_place(place, flags, hops, time, [(0, b), (1, c)])

        else:
            raise ValueError(f"Unknown place: {place}")

    return ans


def main():
    data = sys.stdin.read().strip().split()
    if not data:
        return

    # Чтение a, b, c
    a = float(data[0])
    b = float(data[1])
    c = float(data[2])
    # Чтение v0, v1, v2
    v0 = float(data[3])
    v1 = float(data[4])
    v2 = float(data[5])

    result = solve(a, b, c, [v0, v1, v2])
    # Вывод с достаточной точностью (как в Go по умолчанию через 'g')
    print("{:.15g}".format(result))


if __name__ == "__main__":
    main()
