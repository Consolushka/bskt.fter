package database

type Seeder interface {
	Seed()
	Model() interface{}
}
