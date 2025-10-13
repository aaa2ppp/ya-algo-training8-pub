import sys
import os

# Константы
MAX_PRICE = 1000
MAX_LENGTH = 100
OVER_COST = MAX_LENGTH * MAX_PRICE + 500  # 100500 :)

# Флаг отладки
DEBUG = os.getenv("DEBUG") is not None

class Offer:
    """Предложение на рынке."""
    __slots__ = ("P", "R", "Q", "F")
    def __init__(self, P, R, Q, F):
        self.P = P  # цена за единицу, если количество < R
        self.R = R  # пороговое количество
        self.Q = Q  # цена за единицу, если количество >= R
        self.F = F  # максимальное доступное количество

    def cost(self, n):
        """Рассчитывает стоимость покупки `n` единиц по этому предложению."""
        if n > self.F:
            return OVER_COST
        if n < self.R:
            return self.P * n
        return self.Q * n

def solve(l, offers):
    """
    Решает задачу оптимальной покупки `l` единиц товара по заданным предложениям.
    
    Args:
        l: Требуемое количество товара.
        offers: Список объектов Offer.
        
    Returns:
        Кортеж (минимальная_стоимость, список_количеств_по_каждому_предложению).
        Если покупка невозможна, возвращает (-1, []).
    """
    n = len(offers)
    if n == 0:
        return (0, []) if l == 0 else (-1, [])
    
    # Находим максимальный порог R и общее количество на рынке
    max_r = max(o.R for o in offers)
    total_market = sum(o.F for o in offers)
    
    # Если всего на рынке меньше, чем нужно, решение невозможно
    if total_market < l:
        if DEBUG:
            print(f"too few on the market {total_market} < {l}", file=sys.stderr)
        return -1, []
    
    # `m` - максимальная сумма, которую может понадобиться рассмотреть
    m = max_r + l

    # Инициализация DP-таблицы.
    # `dp[i][j]` будет хранить минимальную стоимость для покупки `j` единиц,
    # используя первые `i+1` предложений.
    # Для восстановления ответа также храним, сколько куплено у i-го продавца.
    dp_cost = [[OVER_COST] * (m + 1) for _ in range(n)]
    dp_choice = [[0] * (m + 1) for _ in range(n)]
    
    # Базовый случай: используем только первое предложение
    for j in range(m + 1):
        dp_cost[0][j] = offers[0].cost(j)
        dp_choice[0][j] = j

    # Заполняем DP-таблицу
    for i in range(1, n):
        for j in range(m + 1):
            # Пытаемся купить `k` единиц у i-го продавца
            for k in range(j + 1):
                prev_cost = dp_cost[i-1][j - k]
                curr_cost = offers[i].cost(k)
                total_cost = prev_cost + curr_cost
                
                if total_cost < dp_cost[i][j]:
                    dp_cost[i][j] = total_cost
                    dp_choice[i][j] = k

    # Находим минимальную стоимость для покупки как минимум `l` единиц
    min_cost = OVER_COST
    best_j = -1
    for j in range(l, m + 1):
        if dp_cost[n-1][j] < min_cost:
            min_cost = dp_cost[n-1][j]
            best_j = j

    if min_cost == OVER_COST:
        # В теории, этого не должно произойти, если total_market >= l
        raise RuntimeError("it's impossible!")

    # Восстанавливаем ответ: сколько куплено у каждого продавца
    chunks = [0] * n
    current_j = best_j
    for i in range(n - 1, -1, -1):
        bought = dp_choice[i][current_j]
        chunks[i] = bought
        current_j -= bought

    return min_cost, chunks

def main():
    """Основная функция для чтения ввода, решения задачи и вывода результата."""
    data = sys.stdin.read().split()
    if not data:
        return

    # Парсим входные данные
    n = int(data[0])
    l = int(data[1])
    offers = []
    idx = 2
    for _ in range(n):
        p = int(data[idx])
        r = int(data[idx+1])
        q = int(data[idx+2])
        f = int(data[idx+3])
        idx += 4
        offers.append(Offer(p, r, q, f))

    # Решаем задачу
    cost, chunks = solve(l, offers)

    # Выводим результат
    print(cost)
    print(" ".join(map(str, chunks)))

if __name__ == "__main__":
    main()
