package main

import (
	"github.com/blumid/gowatch/tasks"
)

func main() {

	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// sig := <-sigchan
	// fmt.Println("Received signal:", sig)

	//start tasks
	tasks.Start()

}
