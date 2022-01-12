# GoRobohash
Robohash Project Go Version, origin from [e1ven/Robohash](https://github.com/e1ven/Robohash)

```
import hashlib

hash = hashlib.sha512()
string="test"
string = string.encode('utf-8')
hash.update(string)
hexdigest = hash.hexdigest()
#print(hexdigest)
hexdigest = "e6509f1865e882eba7ded8250a5194b6bb8e4b685861019dd4947a54966fd4df8349af07f3327380cf4767f49013e0d6371bf0b509ba8d706c67d93df7a469f1"
hasharray = []


for i in range(0,4):
	blocksize = int(len(hexdigest) / 4)
	currentstart = (1 + i) * blocksize - blocksize
	currentend = (1 +i) * blocksize
	# print(hexdigest[currentstart:currentend])
	#temp = int(hexdigest[currentstart:currentend],16)
	# print(temp)
	hasharray.append(int(hexdigest[currentstart:currentend],16))

hasharray = hasharray + hasharray
print(hasharray)
```