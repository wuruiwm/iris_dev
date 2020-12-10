package main

import (
	_ "iris_dev/app/model"
	_ "iris_dev/common"
	_ "iris_dev/config"
	"iris_dev/router"
)

func main(){
	router.HttpInit()
}