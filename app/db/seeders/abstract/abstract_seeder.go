package abstract

type Seeder interface {
	Seed()
	Model() interface{}
}
