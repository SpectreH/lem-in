<h1 align="center">LEM-IN</h1>

## About The Project
Lem-in is an algorithm project, where ants have to go the fastest as posible from start-room to end-room.

## Installation
```
git clone https://github.com/SpectreH/lem-in.git
cd lem-in
```

## Usage
```
go run . [FILE]
```

## Examples
```
$ cat maps/example03.txt
4
4 5 4
##start
0 1 4
1 3 6
##end
5 6 4
2 3 4
3 3 1
0-1
2-4
1-4
0-2
4-5
3-0
4-3

$ go run . maps/example03.txt
L1-3
L1-4 L2-3
L1-5 L2-4 L3-3
L2-5 L3-4 L4-3
L3-5 L4-4
L4-5
```

## Additional information

Only standard go packages were in use. In <code>maps</code> folder you can find several presets to test algorithm and generate paths. Also you can find bad presets examples.

## Author

* SpectreH (https://github.com/SpectreH)
