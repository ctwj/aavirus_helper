package lib

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ===================================
// Debug 函数接收任意类型的参数并打印调试信息
func Debug(args ...interface{}) {
	for _, arg := range args {
		printArg(arg)
	}
	fmt.Println()
}

// printArg 用于打印单个参数
func printArg(arg interface{}) {
	// 如果参数是 map 或者未知类型，将其转换为 JSON 格式
	if argType := fmt.Sprintf("%T", arg); argType == "map[string]interface {}" {
		jsonStr, err := json.MarshalIndent(arg, "", "  ")
		if err == nil {
			fmt.Printf("\033[0;33m[DEBUG]\033[0m (JSON) %v ", string(jsonStr))
		} else {
			fmt.Printf("\033[0;33m[DEBUG]\033[0m (Error marshaling to JSON) %v ", err)
		}
	} else {
		// 使用反射获取参数的类型和值
		argType := fmt.Sprintf("%v", reflect.TypeOf(arg))
		argValue := fmt.Sprintf("%v", reflect.ValueOf(arg))

		// 打印类型和值
		fmt.Printf("\033[0;33m[DEBUG]\033[0m (%v) %v ", argType, argValue)
	}
}
