import sys
from collections import defaultdict

sys.setrecursionlimit(10**6) 

def solve(parent, queries):
    n = len(parent)
    m = len(queries)
    ans = [0] * m

    # Строим мапу вопросов: для каждого узла — список запросов вида (индекс_запроса, предок)
    query_map = defaultdict(list)
    for i, (a, b) in enumerate(queries):
        query_map[b].append((i, a))

    # Строим дерево
    tree = [[] for _ in range(n)]
    for node in range(1, n):
        tree[parent[node]].append(node)

    # Множество текущих предков в DFS
    current_ancestors = set()

    def dfs(node):
        # Отвечаем на все запросы, где текущий узел — это "b"
        if node in query_map:
            for idx, ancestor in query_map[node]:
                if ancestor in current_ancestors:
                    ans[idx] = 1

        current_ancestors.add(node)
        for child in tree[node]:
            dfs(child)
        current_ancestors.discard(node)

    # Запускаем DFS от всех корней (у которых parent[node] == 0)
    for node in range(1, n):
        if parent[node] == 0:
            dfs(node)

    return ans

def main():
    data = sys.stdin.read().split()
    if not data:
        return

    it = iter(data)
    n = int(next(it))
    # В Go parent[0] не используется, и ввод идёт для p[1..n]
    # Здесь parent[0] останется 0, как в оригинале
    parent = [0] * (n + 1)
    for i in range(1, n + 1):
        parent[i] = int(next(it))

    m = int(next(it))
    queries = []
    for _ in range(m):
        a = int(next(it))
        b = int(next(it))
        queries.append((a, b))

    result = solve(parent, queries)

    # Выводим каждый результат на новой строке
    sys.stdout.write("\n".join(str(x) for x in result) + "\n")

if __name__ == "__main__":
    main()
