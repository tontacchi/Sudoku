package solver

/*
Represents variable's possible assignable values
*/
type Domain struct {
	values   []int
	modified bool
}

func (d *Domain) InitValues(values []int) {
	d.values = make([]int, len(values))

	for i, item := range values {
		d.values[i] = item
	}

	d.modified = false
}

// Accessors
func (d *Domain) Contains(value int) bool {
	for item := range d.values {
		if item == value {
			return true
		}
	}

	return false
}

func (d *Domain) Size() int {
	return len(d.values)
}

func (d *Domain) IsEmpty() bool {
	return d.Size() == 0
}

func (d *Domain) HasBeenModified() bool {
	return d.modified
}

// Modifiers
func (d *Domain) Add(value int) bool {
	for _, item := range d.values {
		if item == value { return false }
	}

	d.values = append(d.values, value)
	return true
}

func (d *Domain) Remove(value int) bool {
	for index, item := range d.values {
		if item == value {
			d.values = append(d.values[:index], d.values[index:]...)
			return true
		}
	}

	return false
}

func (d *Domain) SetModifiedTrue() {
	d.modified = true
}

func (d *Domain) SetModifiedFalse() {
	d.modified = false
}

