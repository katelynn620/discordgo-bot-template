package repo

type BaseRepo interface {
	FindOrCreate()
	Update()
}
