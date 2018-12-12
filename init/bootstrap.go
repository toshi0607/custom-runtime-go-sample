package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	log.Println("bootstrap実行中")
	out, err := exec.Command("pwd").Output()
	if err != nil {
		log.Println(err)
	}
	log.Printf("pwd: %s\n", string(out))

	binaryPath := os.Getenv("LAMBDA_TASK_ROOT")+"/"+os.Getenv("_HANDLER")
	log.Println(binaryPath)
	if err := exec.Command(binaryPath).Run(); err != nil {
		log.Println(err)
	}
}
