package storage

import "github.com/AdarshAnand26/Student-API/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents()([]types.Student, error)
}