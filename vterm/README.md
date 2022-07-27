FIXME
=====
- Move all control functions into 1 big switch statement, measure performance
- If damage contains a simple reference to the screen, possible that by
  the time vman processes it, it's being modified by further controls.
  Do I care?
- Make tests parallel, see https://github.com/golang/go/wiki/TableDrivenTests
- Scrollbar?
- Benchmark against tmux.
- Array vs list: https://baptiste-wicht.com/posts/2012/12/cpp-benchmark-vector-list-deque.html

Terminology
===========
Vars:
- cur, start, stop, upper, lower: offsets
- beg, end: coordinates
- row, col: data structures
- x, y, w, h
Method types:
- v: vector, takes pairs
- o: offset, takes int
- e: extent, takes start offset and size
