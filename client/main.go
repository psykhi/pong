package main

func main() {
	done := make(chan struct{})
	c := NewClient()
	c.Start()
	defer func() {
		c.Stop()
	}()
	<-done
}
