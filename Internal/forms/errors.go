package forms

type errors map[string][]string

// Add an erro mesaage for given form filed
func (e errors) Add(filed, message string) {
	e[filed] = append(e[filed], message)
}

// Retuns the first eroor message
func (e errors) Get(filed string) string {
	es := e[filed]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
