/*
 * Copyright 2025 Coze Studio. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nodes

// 注意：循环引用问题修复
// 由于各个节点包（adapter, capability, component）都在自己的 init() 函数中
// 自动注册了节点，所以不需要在这里显式导入和注册
// 导入这些包会造成循环引用：nodes -> nodes/adapter -> nodes
// 
// import (
// 	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/adapter"
// 	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/capability"
// 	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/component"
// )

// RegisterAllNewNodes 注册所有新增的工作流节点
// 注意：由于各节点包已经在 init() 中自动注册，此函数当前为空操作
// 保留此函数是为了向后兼容和未来可能的显式注册需求
func RegisterAllNewNodes() {
	// 节点已经通过各自包的 init() 函数自动注册
	// 不需要显式调用注册函数
	// 
	// 原先的设计：
	// - 注册适配器节点
	// - 注册AI能力节点  
	// - 注册组件节点
}

// 以下函数已废弃 - 各节点包通过 init() 自动注册
// 
// // registerAdapterNodes 注册适配器节点
// func registerAdapterNodes() {
// 	// 节点已在 adapter/adapter_invoke.go 的 init() 中注册
// }
// 
// // registerCapabilityNodes 注册AI能力节点
// func registerCapabilityNodes() {
// 	// 节点已在各 capability 文件的 init() 中注册
// }
// 
// // registerComponentNodes 注册组件节点
// func registerComponentNodes() {
// 	// 节点已在各 component 文件的 init() 中注册
// }

