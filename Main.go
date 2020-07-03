package main

func main() {
	// connecting to database
	app := Config{}
	app.Conn("localhost", 5432, "postgres", "packform")
	app.Start(":8000")

}