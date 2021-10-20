package main

var (
	// Version for build
	Version string
	// Build for build
	Build string
)

func main() {
	a := App{}
	a.Initialize()
	a.Run()
}
