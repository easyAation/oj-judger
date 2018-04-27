import string
while True:
    try:
        s = raw_input()
        arr = string.split(s,' ')
        a = int(arr[0])
        b = int(arr[1])
        print a + b
    except EOFError:
        exit(0)