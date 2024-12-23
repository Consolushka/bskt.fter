package main

import "FTER/app/cmd"

func main() {
	//cmd.GamePdf("0022400377")
	//dsn := "host=localhost user=postgres password=postgres dbname=fter port=5432 sslmode=disable"
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//fmt.Print(db.Name())
	cmd.Execute()
}
