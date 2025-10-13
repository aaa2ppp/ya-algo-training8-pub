import random
import time
from main import solve

def benchmark_solve():
    n, m = 1000, 1000
    print(f"Generating {n}x{m} matrix...")
    mx = [random.randint(1, 10 * n * m) for _ in range(n*m)]
    
    print("Warming up...")
    solve(mx, n, m)  # прогрев
    
    print("Running benchmark...")
    start = time.perf_counter()
    iterations = 5
    for _ in range(iterations):
        solve(mx, n, m)
    elapsed = time.perf_counter() - start
    
    print(f"✅ {iterations} runs in {elapsed:.3f} sec")
    print(f"⏱️  Avg: {elapsed/iterations*1000:.1f} ms per run")

if __name__ == "__main__":
    benchmark_solve()
