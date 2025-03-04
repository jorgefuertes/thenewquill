package item

func (i *Item) WeightTotal() int {
	return recursiveWeight(i)
}

func recursiveWeight(i *Item) int {
	if !i.IsContainer || len(i.Contents) == 0 {
		return i.Weight
	}

	w := i.Weight
	for _, content := range i.Contents {
		w += recursiveWeight(content)
	}

	return w
}

func (i Item) canCarry(w int) bool {
	return i.WeightTotal()+w <= i.MaxWeight
}

func (i Item) isFull() bool {
	if !i.IsContainer {
		return true
	}

	return i.WeightTotal() >= i.MaxWeight
}

func (i *Item) Put(a *Item) error {
	if !i.IsContainer {
		return ErrNotContainer
	}

	if i.isFull() {
		return ErrContainerIsFull
	}

	if !i.canCarry(a.WeightTotal()) {
		return ErrContainerCantCarrySoMuch
	}

	i.Contents = append(i.Contents, a)

	return nil
}
