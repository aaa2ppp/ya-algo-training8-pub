import time
from collections import defaultdict

n = 1_000_000
trees = [(i % 1000, i // 1000) for i in range(n)]

# defaultdict
start = time.time()
forest1 = defaultdict(int)
for x, y in trees:
    forest1[(x, y)] += 1
t1 = time.time() - start

# dict.get
start = time.time()
forest2 = {}
for x, y in trees:
    key = (x, y)
    forest2[key] = forest2.get(key, 0) + 1
t2 = time.time() - start

print(f"defaultdict: {t1:.4f}s")
print(f"dict.get:    {t2:.4f}s")
