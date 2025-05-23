package solver

import (
	"slices"
	"fmt"
)

// Represents set of all possible assignable to variable
type Domain struct {
	values   []int
	modified bool
}

func NewDomain(values ...int) *Domain {
	domain := &Domain{
		values:   []int{},
		modified: false,
	}

	domain.values = append(domain.values, values...)

	return domain
}

func (d *Domain) Copy() *Domain {
	domainCopy := &Domain{
		values:   slices.Clone(d.values),
		modified: false,
	}

	return domainCopy
}

// Accessors
func (d *Domain) Contains(value int) bool {
	return slices.Contains(d.values, value)
}

func (d *Domain) Size() int {
	return len(d.values)
}

func (d *Domain) Empty() bool {
	return len(d.values) == 0
}

func (d *Domain) Modified() bool {
	return d.modified
}

// Mutators
func (d *Domain) Expand(value int) {
	if !d.Contains(value) {
		d.values = append(d.values, value)
	}
}

// modified -> true when values emptied out
func (d *Domain) Remove(value int) bool {
	if d.Empty() || !d.Contains(value) { return false }

	for index, item := range d.values {
		if item == value {
			d.values = append(d.values[:index], d.values[index+1:]...)
			d.modified = true

			break
		}
	}

	return true
}

func (d *Domain) SetModified(modifer bool) {
	d.modified = modifer
}


/*
Domain:
  values:   [1 2 3 ...]
  modified: true
*/
func (d *Domain) String() string {
	res := "Domain:\n"
	res += fmt.Sprintf("  values:   %v\n", d.values)
	res += fmt.Sprintf("  modified: %v", d.modified)

	return res
}

