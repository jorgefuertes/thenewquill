package console

func (i *input) historyAdd() {
	if len(i.history) >= historySize {
		i.history = i.history[1:]
	}

	for _, h := range i.history {
		if string(h) == string(i.current) {
			return
		}
	}

	i.history = append(i.history, i.current)
	i.index = len(i.history) - 1
}

func (i *input) historyPrev() {
	i.index--
	if i.index < 0 {
		i.index = 0
	}

	i.current = i.history[i.index]
}

func (i *input) historyNext() {
	i.index++
	if i.index >= len(i.history) {
		i.index = len(i.history) - 1
	}

	i.current = i.history[i.index]
}
