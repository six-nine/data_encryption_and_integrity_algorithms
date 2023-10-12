prime = [True for _ in range(100000)]
for i in range(2, 100000):
    if prime[i]:
        j = i
        while j + i < 100000:
            j += i
            prime[j] = False

for i in range(2, 100):
    if prime[i]:
        print(i)


for i in range(2, 100000):
    if prime[i] and i % 4 == 3:
        print(i)