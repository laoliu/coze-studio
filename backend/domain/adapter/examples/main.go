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
	demoType := flag.String("demo", "k12", "Demo type: k12, registry, or all")
	flag.Parse()

	fmt.Println("==============================================")
	fmt.Println("   Coze Studio Adapter Demo")
	fmt.Println("==============================================")
	fmt.Println()

	switch *demoType {
	case "k12":
		DemoK12Adapter()
	case "registry":
		DemoAdapterRegistry()
	case "all":
		DemoK12Adapter()
		fmt.Println("\n" + string(make([]byte, 50, 50)) + "\n")
		DemoAdapterRegistry()
	default:
		fmt.Printf("Unknown demo type: %s\n", *demoType)
		fmt.Println("Available options: k12, registry, all")
		fmt.Println()
		fmt.Println("Note: 'real' and 'coze' demos require full Coze infrastructure")
		fmt.Println("      and are currently disabled. See coze_*.go files for details.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("==============================================")
	fmt.Println("   Demo Finished Successfully!")
	fmt.Println("==============================================")
}
