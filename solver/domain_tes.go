package solver

import (
	"testing"
	"fmt"
)

func TestBlankCellDomain(t *testing.T) {
	const DIMENSION int = 9
	startingDomainValues := make([]int, DIMENSION)
	for num := 1; num <= DIMENSION; num++ {
		startingDomainValues[num-1] = num
	}


	message := fmt.Sprintf("\nstarting domain values:\n%v\n", startingDomainValues)
	t.Log(message)

	domain := NewDomain(startingDomainValues...)

	message = fmt.Sprintf("\nstarting domain:\n%v\n", domain.String())
	t.Logf(message)
}

func TestDomainInit(t *testing.T) {
	domain_1 := NewDomain()
	domain_2 := NewDomain([]int{}...)
	domain_3 := NewDomain([]int{0}...)
	domain_4 := NewDomain(2,4,5)
	domain_5 := NewDomain([]int{1, 2, 3}...)

	t.Log(domain_1)
	t.Log(domain_2)
	t.Log(domain_3)
	t.Log(domain_4)
	t.Log(domain_5)
}

func TestAdd(t *testing.T) {
	domain := NewDomain()

	message := fmt.Sprintf("\ndomain is empty: %v\n%v\n", domain.Empty(), domain.String())
	t.Logf(message)

	domain.Add(1)
	message = fmt.Sprintf("\nAdded 1\ndomain is empty: %v\n%v\n", domain.Empty(), domain.String())
	t.Logf(message)

	domain.Add(2)
	message = fmt.Sprintf("\nAdded 2\ndomain is empty: %v\n%v\n", domain.Empty(), domain.String())
	t.Logf(message)

	domain.Add(1)
	message = fmt.Sprintf("\nAdded 1 again\ndomain is empty: %v\n%v\n", domain.Empty(), domain.String())
	t.Logf(message)
}

func TestRemove(t *testing.T) {
	domain := NewDomain()
	popped := domain.Remove(100)

	message := fmt.Sprintf("\nvalue was removed: %v\ndomain is empty: %v\n%v\n", popped, domain.Empty(), domain.String())
	t.Log(message)
	
	domain.Add(3)
	domain.Add(4)
	message = fmt.Sprintf("\ndomain is empty: %v\n%v\n", domain.Empty(), domain.String())
	t.Log(message)

	popped  = domain.Remove(3)
	message = fmt.Sprintf("\nvalue was removed: %v\ndomain is empty: %v\n%v\n", popped, domain.Empty(), domain.String())
	t.Log(message)

	popped = domain.Remove(3)
	message = fmt.Sprintf("\nvalue was removed: %v\ndomain is empty: %v\n%v\n", popped, domain.Empty(), domain.String())
	t.Log(message)

	popped = domain.Remove(4)
	message = fmt.Sprintf("\nvalue was removed: %v\ndomain is empty: %v\n%v\n", popped, domain.Empty(), domain.String())
	t.Log(message)
}

func TestContains(t *testing.T) {
	domain := NewDomain()

	check := domain.Contains(1)
	message := fmt.Sprintf("\nempty domain contains 1: %v\n", check)
	t.Log(message)

	domain.Add(1)
	message = fmt.Sprintf("\nadded 1\n%v\n", domain.values)
	t.Log(message)

	check = domain.Contains(1)
	message = fmt.Sprintf("\ndomain contains 1: %v\n", check)
	t.Log(message)

	check = domain.Contains(2)
	message = fmt.Sprintf("\ndomain contains 2: %v\n", check)
	t.Log(message)
}

