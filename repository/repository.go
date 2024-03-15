package repository

type Repository interface {
	GetById(id int64)
	Add(interface{})
	Update(id int64, entity interface{})
}
