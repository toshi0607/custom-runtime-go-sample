package main

import (
	"log"
	"os"
	"os/exec"
)

var(
	handler string
)

func init() {
	handler = os.Getenv("LAMBDA_TASK_ROOT")+"/"+os.Getenv("_HANDLER")
}

func main() {
	log.Println("bootstrap実行中")
	out, err := exec.Command("pwd").Output()
	if err != nil {
		log.Println(err)
	}
	log.Printf("pwd: %s\n", string(out))

	log.Println(handler)
	if err := exec.Command(handler).Run(); err != nil {
		log.Println(err)
	}
}
