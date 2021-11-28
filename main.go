package main

func main() {
	s := newServer()
	go s.run()
	go s.runHttpService()
	s.runTcp()
}
