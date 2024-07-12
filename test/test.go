package main

import (
	"fmt"
	"strconv"
)

func main() {
	var str string = "5"
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
	} else {
		fmt.Printf("%.1f %v", floatValue, floatValue)
	}

}
