package obj

func (i Item) IsContainer() bool {
	return i.isContainer
}

func (i *Item) WeightTotal() int {
	return recursiveWeight(i)
}

func recursiveWeight(i *Item) int {
	if !i.IsContainer() || len(i.contents) == 0 {
		return i.weight
	}

	w := i.weight
	for _, content := range i.contents {
		w += recursiveWeight(content)
	}

	return w
}

func (i Item) canCarry(w int) bool {
	return i.WeightTotal()+w <= i.maxWeight
}

func (i Item) isFull() bool {
	if !i.IsContainer() {
		return true
	}

	return i.WeightTotal() >= i.maxWeight
}

func (i *Item) Put(a *Item) error {
	if !i.IsContainer() {
		return ErrNotContainer
	}

	if i.isFull() {
		return ErrContainerIsFull
	}

	if !i.canCarry(a.WeightTotal()) {
		return ErrContainerCantCarrySoMuch
	}

	i.contents = append(i.contents, a)

	return nil
}
