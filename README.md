# Connect4 solver

Highly optimized connect4 solver written in Go based on minimax with:
- Alpha-Beta pruning
- Iterative deepening
- Transposition table
- Move ordering
- Game representation as uint64

This is mostly a reimplementation of http://blog.gamesolver.org/solving-connect-four/

# TODO
- [x] Set up initial board
- [x] Print / display board in terminal
- [x] Place stones manually via terminal
- [x] Check for win conditions
- [x] Test win conditions
- [x] Check for draw conditions
- [x] Implement minimax solver
- [x] Optimize engine & solver speed
- [x] Implement alpha-beta pruning
- [x] Implement transposition table
- [x] Implement iterative deepening with null window search
- [x] Implement board with bit mapping
- [x] Anticipating losing moves
- [x] Better move ordering
- [ ] Better move ordering 2 (include pair alignments)
- [ ] Optimize transposition table