// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"flag"
	"fmt"
	"os"
)

// RunCLI 运行命令行接口
func RunCLI() {
	// 定义命令行参数
	demoType := flag.String("demo", "k12", "Demo type: k12, real, coze, registry, or all")
	flag.Parse()

	fmt.Println("==============================================")
	fmt.Println("   Coze Studio Adapter Demo")
	fmt.Println("==============================================")
	fmt.Println()

	switch *demoType {
	case "k12":
		DemoK12Adapter()
	case "real":
		DemoRealK12Adapter()
	case "coze":
		DemoCozeK12Adapter()
	case "registry":
		DemoAdapterRegistry()
	case "all":
		DemoK12Adapter()
		fmt.Println("\n" + string(make([]byte, 50, 50)) + "\n")
		DemoAdapterRegistry()
	default:
		fmt.Printf("Unknown demo type: %s\n", *demoType)
		fmt.Println("Available options: k12, real, coze, registry, all")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("   Demo Finished Successfully!")
	fmt.Println("==============================================")
}
