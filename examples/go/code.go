package main

import (
	"fmt"
)

const MaxNameLen = 32

type Student struct {
	ID    int
	Name  string
	Grade float64
}

func printStudent(s *Student) {
	if s == nil {
		fmt.Println("invalid student pointer")
		return
	}
	fmt.Printf("id: %d\nname: %s\ngrade: %.2f\n\n", s.ID, s.Name, s.Grade)
}

func averageGrade(students []Student) float64 {
	if len(students) == 0 {
		return 0.0
	}

	var sum float64
	for _, s := range students {
		sum += s.Grade
	}
	return sum / float64(len(students))
}

func main() {
	students := []Student{
		{ID: 1, Name: "alice", Grade: 91.2},
		{ID: 2, Name: "bob", Grade: 85.5},
		{ID: 3, Name: "charlie", Grade: 77.8},
	}

	fmt.Println("=== student records ===")
	for i := range students {
		printStudent(&students[i])
	}

	avg := averageGrade(students)
	fmt.Printf("average grade: %.2f\n", avg)

	// comment
	newStudent := &Student{ID: 4, Name: "diana", Grade: 88.9}

	fmt.Println("\n=== new student ===")
	printStudent(newStudent)
}
