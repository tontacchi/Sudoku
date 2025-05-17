package solver

type trailEntry struct {
	variable *Variable
	domain   *Domain
}

// Represents changes for easier forward propagation
type Trail struct {
	stack     []trailEntry
	markers   []int
	numPushes int
	numUndoes int
}

func NewTrail() *Trail {
	return &Trail{
		stack:   []trailEntry{},
		markers: []int{},
	}
}

// Accessors
func (t *Trail) Size() int {
	return len(t.stack)
}

// Mutators
func (t *Trail) PlaceMarker() {
	t.markers = append(t.markers, len(t.stack))
}

func (t *Trail) Push(variable *Variable) {
	entry := trailEntry{
		variable: variable,
		domain:   variable.domain.Copy(),
	}

	t.stack = append(t.stack, entry)
	t.numPushes++
}

func (t *Trail) Undo() {
	if (len(t.markers) == 0) { return }

	targetSize := t.markers[len(t.markers) - 1]
	t.markers   = t.markers[:len(t.markers) - 1]

	for len(t.stack) > targetSize {
		entry   := t.stack[len(t.stack) - 1]
		t.stack  = t.stack[:len(t.stack) - 1]

		entry.variable.domain = entry.domain
		entry.variable.modified = false
		entry.variable.Unassign()
	}

	t.numUndoes++
}

func (t *Trail) Clear() {
	t.stack   = []trailEntry{}
	t.markers = []int{}

	t.numPushes, t.numUndoes = 0, 0
}

