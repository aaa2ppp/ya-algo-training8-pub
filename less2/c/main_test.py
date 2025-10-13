import sys
import os
import math
import io

# Импортируем функцию solve из main.py
from main import solve

def almost_equal(a: float, b: float, rel_tol=1e-4, abs_tol=1e-4) -> bool:
    return abs(a - b) <= max(rel_tol * max(abs(a), abs(b)), abs_tol)

def test_solve():
    test_cases = [
        # (название, вход_как_список_интервалов, ожидаемый_ответ)
        ("1 - базовый случай", [(0, 1, 1), (0.5, 1.5, 1.5), (1, 2, 1)], 2.0),
        ("2 - отрицательные координаты", [(-2, -1, 1), (-1.5, -0.5, 1.5), (-1, 0, 1)], 2.0),
        ("3 - zero intervals", [], 0.0),
        ("4 - single interval", [(0, 1, 2.5)], 2.5),
        ("5 - same start/end", [(0, 1, 1), (0, 2, 2), (1, 2, 1)], 2.0),
        ("6 - nested", [(0, 5, 1), (1, 4, 2), (2, 3, 3)], 3.0),
        ("7 - tiny disjoint", [(0, 0.0001, 1), (0.0001, 0.0002, 1), (0.0002, 0.0003, 1)], 3.0),
        ("8 - large weights", [(0, 1, 1000), (0.5, 1.5, 2000), (1, 2, 3000)], 4000.0),
        ("9 - reverse disjoint", [(3, 4, 1), (2, 3, 1), (1, 2, 1)], 3.0),
        ("10 - same endpoints", [(0, 1, 1)] * 5, 1.0),
        ("11 - negative coords", [(-5, -3, 1), (-4, -2, 2), (-3, -1, 3)], 4.0),
        ("12 - float precision disjoint", [(0.1, 0.2, 0.1), (0.2, 0.3, 0.1), (0.3, 0.4, 0.1)], 0.3),
        ("13 - disjoint", [(0, 1, 1), (2, 3, 1), (4, 5, 1)], 3.0),
        ("14 - all start together", [(0, 1, 1), (0, 2, 1), (0, 3, 1), (0, 4, 1)], 1.0),
        ("15 - all end together", [(0, 10, 1), (1, 10, 1), (2, 10, 1), (3, 10, 1)], 1.0),
        ("16 - one contains all", [(0, 10, 1), (1, 2, 2), (3, 4, 3), (5, 6, 4)], 9.0),
        ("17 - close boundaries", [(0, 1, 1), (1, 2, 1), (1, 1.0000001, 1)], 2.0),
        ("18 - all intersect", [(0, 100, 0.01)] * 100, 0.01),
        ("19 - complex choice", [(0, 2, 1), (1, 3, 2), (2, 4, 3), (3, 5, 2), (4, 6, 1)], 5.0),
        ("20 - two heavy non-overlapping", [(0, 3, 5), (1, 2, 3), (2, 4, 3), (3, 5, 5)], 10.0),
        ("21 - touching points", [(0, 1, 2), (1, 2, 2), (2, 3, 2)], 6.0),
        ("22 - heavy vs light", [(0, 2, 10), (1, 3, 1), (2, 4, 1), (3, 5, 1)], 11.0),
        ("23 - greedy fails", [(0, 3, 3), (2, 4, 3), (3, 5, 3)], 6.0),
        ("24 - chain", [(0, 1, 2), (1, 2, 3), (2, 3, 4), (3, 4, 5)], 14.0),
        ("25 - alternative paths", [(0, 2, 5), (1, 3, 4), (2, 4, 3), (3, 5, 2)], 8.0),
    ]

    for name, intervals, expected in test_cases:
        result = solve(intervals)
        if not almost_equal(result, expected):
            print(f"❌ {name}: got {result}, want {expected}", file=sys.stderr)
            sys.exit(1)
        else:
            print(f"✅ {name}")

    print("🎉 Все тесты пройдены!")

if __name__ == "__main__":
    test_solve()
