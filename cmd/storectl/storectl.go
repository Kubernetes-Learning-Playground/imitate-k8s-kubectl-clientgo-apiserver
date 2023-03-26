package main

import storectl "practice_ctl/pkg/storectl/cmd"

func main() {
	storectl.RunCmd()
}

/*
	测试用命令：
	go run storectl.go create cars ../../json/car.json
	go run storectl.go delete cars initCar
	go run storectl.go describe apples initApple
	go run storectl.go apply cars ../../json/car.json
	go run storectl.go list cars
*/