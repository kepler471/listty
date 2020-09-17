package main

type Position struct {
	Depth int
	Row   int
}

type PositionStack struct {
	AllPositions []Position
}

func (pStack *PositionStack) IsEmpty() bool {
	return &pStack.AllPositions == nil
}

func (pStack *PositionStack) GetLast() *Position {
	if pStack.IsEmpty() {
		return &Position{
			Depth: 0,
			Row:   0,
		}
	}

	return pStack.GetPosition(pStack.Size())
}

func (pStack *PositionStack) GetPosition(depth int) *Position {
	if pStack.IsEmpty() {
		return &Position{
			Depth: 0,
			Row:   0,
		}
	}

	return &pStack.AllPositions[depth]
}

func (pStack *PositionStack) GetRow(depth int) int {
	if pStack.IsEmpty() {
		return 0
	}
	return pStack.GetPosition(depth).Row
}

func (pStack *PositionStack) SetRow(depth, row int) {
	if pStack.IsEmpty() {
		return
	}
	pStack.GetPosition(depth).Row = row
}

func (pStack *PositionStack) AddPosition(depth, row int) {
	pStack.Push(Position{
		Depth: depth,
		Row:   row,
	})
}

func (pStack *PositionStack) Push(p Position) {
	pStack.AllPositions = append(pStack.AllPositions, p)
}

func (pStack *PositionStack) Pop() {
	if pStack.Size() == 0 || pStack.IsEmpty() {
		return
	}

	pStack.AllPositions = pStack.AllPositions[:pStack.Size()]
}

func (pStack *PositionStack) Size() int {
	if pStack.IsEmpty() {
		return 0
	}
	return len(pStack.AllPositions) - 1
}
