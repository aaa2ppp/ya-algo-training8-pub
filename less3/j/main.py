import sys
from collections import deque

def solve(a, b):
    n = len(a)
    
    def check(k):
        # Две очереди: required и available
        required = deque()
        available = deque()
        
        for i in range(n):
            required.append([i, a[i]])
            available.append([i, b[i]])
            
            # Если самый старый требуемый элемент выходит за пределы окна k
            if i - required[0][0] > k:
                return False
            
            # Удаляем из available элементы, вышедшие за пределы окна
            while available and i - available[0][0] > k:
                available.popleft()
            
            # Сопоставляем спрос и предложение
            while required and available:
                _, req_seats = required[0]
                _, av_seats = available[0]
                
                if req_seats < av_seats:
                    available[0][1] -= req_seats
                    required.popleft()
                elif req_seats > av_seats:
                    required[0][1] -= av_seats
                    available.popleft()
                else:
                    required.popleft()
                    available.popleft()
        
        return not required

    # Бинарный поиск по k от 0 до n
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
