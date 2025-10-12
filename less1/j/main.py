import sys

class ListWrapper:
    def __init__(self, items):
        self.items = items

    def sublist(self, idx, size):
        return SubList(self, idx, size)

    def set(self, idx, val):
        self.items[idx] = val

    def add(self, val):
        self.items.append(val)

    def get(self, idx):
        return self.items[idx]

class SubList:
    def __init__(self, parent_list, idx, size):
        self.parent = parent_list
        self.idx = idx
        self.size = size

    def sublist(self, idx, size):
        return self.parent.sublist(self.idx + idx, size)

    def set(self, idx, val):
        self.parent.set(self.idx + idx, val)

    def add(self, val):
        # Stub — sublists don't support add
        pass

    def get(self, idx):
        return self.parent.get(self.idx + idx)

class Figon:
    def __init__(self):
        self.lists = {}

    def parse_list_name(self, s):
        s = s.strip()
        first_space = s.find(' ')
        s = s[first_space + 1:]
        second_space = s.find(' ')
        return s[:second_space]

    def parse_dot_name(self, s):
        s = s.strip()
        return s[:s.find('.')]

    def parse_args(self, s):
        start = s.find('(')
        end = s.find(')', start)
        args_str = s[start + 1:end]
        return [int(x.strip()) for x in args_str.split(',')]

    def parse_new_list(self, s):
        dst_name = self.parse_list_name(s)
        eq_pos = s.find('=')
        rhs = s[eq_pos + 1:]
        args = self.parse_args(rhs)
        self.lists[dst_name] = ListWrapper(args[:])  # make a copy

    def parse_sub_list(self, s):
        dst_name = self.parse_list_name(s)
        eq_pos = s.find('=')
        rhs = s[eq_pos + 1:]
        args = self.parse_args(rhs)
        src_name = self.parse_dot_name(rhs)
        from_idx, to_idx = args[0], args[1]
        # convert to 0-based, inclusive range → size = to - from + 1
        sub = self.lists[src_name].sublist(from_idx - 1, to_idx - from_idx + 1)
        self.lists[dst_name] = sub

    def parse_set(self, s):
        name = self.parse_dot_name(s)
        args = self.parse_args(s)
        idx, val = args[0], args[1]
        self.lists[name].set(idx - 1, val)

    def parse_add(self, s):
        name = self.parse_dot_name(s)
        args = self.parse_args(s)
        val = args[0]
        self.lists[name].add(val)

    def parse_get(self, s):
        name = self.parse_dot_name(s)
        args = self.parse_args(s)
        idx = args[0]
        return self.lists[name].get(idx - 1)

def solve(n, text):
    figon = Figon()
    ans = []

    for line in text:
        s = line.strip()
        if s.startswith("List "):
            if " new " in s:
                figon.parse_new_list(s)
            else:
                figon.parse_sub_list(s)
        elif ".set(" in s:
            figon.parse_set(s)
        elif ".add(" in s:
            figon.parse_add(s)
        elif ".get(" in s:
            ans.append(figon.parse_get(s))
        else:
            raise ValueError(f"unknown syntax: {s!r}")

    return ans

def main():
    data = sys.stdin.read().splitlines()
    if not data:
        return

    n = int(data[0])
    text = data[1:1 + n]

    result = solve(n, text)
    for val in result:
        print(val)

if __name__ == "__main__":
    main()
