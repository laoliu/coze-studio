// Copyright 2025 Coze Studio. All rights reserved.

package examples

import (
	"testing"
)

// TestDemoK12Adapter 测试K12适配器演示
func TestDemoK12Adapter(t *testing.T) {
	// 运行K12适配器演示
	DemoK12Adapter()
}

// TestDemoAdapterRegistry 测试适配器注册表演示
func TestDemoAdapterRegistry(t *testing.T) {
	// 运行适配器注册表演示
	DemoAdapterRegistry()
}

// BenchmarkK12Adapter 性能测试K12适配器
func BenchmarkK12Adapter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DemoK12Adapter()
	}
}
