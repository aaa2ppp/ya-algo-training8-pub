import sys
from collections import deque

def solve(a, b):
    n = len(a)
    
    # Создаём очереди один раз
    req_idx = deque()
    req_val = deque()
    av_idx = deque()
    av_val = deque()
    
    def check(k):
        # Быстрая очистка — без аллокаций
        req_idx.clear()
        req_val.clear()
        av_idx.clear()
        av_val.clear()
        
        for i in range(n):
            req_idx.append(i)
            req_val.append(a[i])
            av_idx.append(i)
            av_val.append(b[i])
            
            if i - req_idx[0] > k:
                return False
                
            while av_idx and i - av_idx[0] > k:
                av_idx.popleft()
                av_val.popleft()
            
            while req_val and av_val:
                if req_val[0] < av_val[0]:
                    av_val[0] -= req_val.popleft()
                    req_idx.popleft()
                elif req_val[0] > av_val[0]:
                    req_val[0] -= av_val.popleft()
                    av_idx.popleft()
                else:
                    req_val.popleft()
                    req_idx.popleft()
                    av_val.popleft()
                    av_idx.popleft()
        
        return not req_val

    lo, hi = 0, n
    while lo < hi:
        mid = (lo + hi) // 2
        if check(mid):
            hi = mid
        else:
            lo = mid + 1

    return -1 if lo == n else lo


def main():
    data = sys.stdin.read().split()
    if not data:
        return
    
    n = int(data[0])
    a = list(map(int, data[1:n+1]))
    b = list(map(int, data[n+1:n+1+n]))
    
    ans = solve(a, b)
    sys.stdout.write(str(ans))


if __name__ == "__main__":
    main()
