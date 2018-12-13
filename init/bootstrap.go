package main

import (
	"log"
	"os"
	"os/exec"
)

var(
	handlerPath string
)

func init() {
	handlerPath = os.Getenv("LAMBDA_TASK_ROOT")+"/"+os.Getenv("_HANDLER")
}

func main() {
	log.Println("bootstrap started")
	out, err := exec.Command("pwd").Output()
	if err != nil {
		log.Fatalf("failed to exec pwd. error: %v", err)
	}
	log.Printf("pwd: %s\n", string(out))

	log.Println(handlerPath)
	if err := exec.Command(handlerPath).Run(); err != nil {
		log.Fatalf("failed to exec handler. error: %v", err)
	}
}
