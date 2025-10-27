import sys
import bisect
import array


def main():
    stdin = sys.stdin.buffer
    stdout = sys.stdout

    n = int(stdin.readline())

    b = array.array('i', (0,)) * n
    t = array.array('i', (0,)) * n

    for i in range(n):
        line = stdin.readline()
        i_space = line.find(b' ')
        b[i] = int(line[:i_space])
        t[i] = int(line[i_space+1:])        

    m = int(stdin.readline())

    for _ in range(m):
        q = int(stdin.readline())
        j = bisect.bisect_left(b, q)
        stdout.write(str(t[j - 1] * q))
        stdout.write('\n')


if __name__ == '__main__':
    main()
