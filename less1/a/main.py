import sys
import os

DEBUG = os.getenv("DEBUG") is not None

def solve(a):
    n = len(a)
    vasya = 0
    masha = 0
    min_vasya = float('inf')
    max_masha = 0

    for i, v in enumerate(a):
        if i % 2 == 0:  # even indices (0, 2, 4, ...) - Vasya
            vasya += v
            min_vasya = min(min_vasya, v)
        else:  # odd indices (1, 3, 5, ...) - Masha
            masha += v
            max_masha = max(max_masha, v)

    if DEBUG:
        print(f"{vasya} {min_vasya} {masha} {max_masha}", file=sys.stderr)

    # If we can swap one element to improve Vasya's advantage
    if n >= 2 and min_vasya < max_masha:
        vasya -= min_vasya
        masha -= max_masha
        vasya += max_masha
        masha += min_vasya

    return vasya - masha

def main():
    data = sys.stdin.read().split()
    if not data:
        return
    
    n = int(data[0])
    a = list(map(int, data[1:n+1]))
    
    result = solve(a)
    print(result)

if __name__ == "__main__":
    main()
