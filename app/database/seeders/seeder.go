package seeders

type Seeder interface {
	Seed()
	Model() interface{}
}
