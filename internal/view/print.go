package view

type contentUpdateMsg struct{}

func (m *model) Println(s string) {
	m.content += s + "\n"

	if m.p != nil {
		m.p.Send(contentUpdateMsg{})
	}
}
