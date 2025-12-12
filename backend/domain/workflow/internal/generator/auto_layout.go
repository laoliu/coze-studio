package generator

// AutoLayoutEngine 自动布局引擎
type AutoLayoutEngine struct {
	nodeWidth     float64
	nodeHeight    float64
	horizontalGap float64
	verticalGap   float64
}

// NewAutoLayoutEngine 创建自动布局引擎
func NewAutoLayoutEngine() *AutoLayoutEngine {
	return &AutoLayoutEngine{
		nodeWidth:     200,
		nodeHeight:    80,
		horizontalGap: 100,
		verticalGap:   150,
	}
}

// AutoLayout 自动计算节点位置
func (e *AutoLayoutEngine) AutoLayout(workflow *ConfiguredWorkflow) *ConfiguredWorkflow {
	if len(workflow.Nodes) == 0 {
		return workflow
	}

	// 1. 构建图结构
	graph := e.buildGraph(workflow)

	// 2. 拓扑排序，计算节点层级
	levels := e.topologicalSort(graph, workflow.Nodes)

	// 3. 层次布局
	e.hierarchicalLayout(workflow.Nodes, levels)

	return workflow
}

// buildGraph 构建图的邻接表
func (e *AutoLayoutEngine) buildGraph(workflow *ConfiguredWorkflow) map[string][]string {
	graph := make(map[string][]string)

	// 初始化所有节点
	for _, node := range workflow.Nodes {
		graph[node.ID] = make([]string, 0)
	}

	// 构建边
	for _, edge := range workflow.Edges {
		graph[edge.From] = append(graph[edge.From], edge.To)
	}

	return graph
}

// topologicalSort 拓扑排序，计算每个节点的层级
func (e *AutoLayoutEngine) topologicalSort(graph map[string][]string, nodes []*InternalNode) map[string]int {
	levels := make(map[string]int)
	visited := make(map[string]bool)

	// 计算入度
	inDegree := make(map[string]int)
	for nodeID := range graph {
		inDegree[nodeID] = 0
	}
	for _, neighbors := range graph {
		for _, neighbor := range neighbors {
			inDegree[neighbor]++
		}
	}

	// BFS 计算层级
	queue := make([]string, 0)
	for _, node := range nodes {
		if inDegree[node.ID] == 0 {
			queue = append(queue, node.ID)
			levels[node.ID] = 0
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visited[current] = true

		currentLevel := levels[current]

		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				// 更新层级
				if level, exists := levels[neighbor]; !exists || currentLevel+1 > level {
					levels[neighbor] = currentLevel + 1
				}

				// 减少入度
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// 处理未访问的节点（可能存在环或孤立节点）
	currentLevel := 0
	for _, node := range nodes {
		if _, exists := levels[node.ID]; !exists {
			levels[node.ID] = currentLevel
			currentLevel++
		}
	}

	return levels
}

// hierarchicalLayout 层次布局算法
func (e *AutoLayoutEngine) hierarchicalLayout(nodes []*InternalNode, levels map[string]int) {
	// 统计每层的节点数量
	levelNodes := make(map[int][]string)
	maxLevel := 0

	for _, node := range nodes {
		level := levels[node.ID]
		if level > maxLevel {
			maxLevel = level
		}
		levelNodes[level] = append(levelNodes[level], node.ID)
	}

	// 为每层的节点分配位置
	startX := 50.0
	startY := 50.0

	for level := 0; level <= maxLevel; level++ {
		nodesInLevel := levelNodes[level]
		if len(nodesInLevel) == 0 {
			continue
		}

		// 计算起始 X 坐标
		currentX := startX

		for i, nodeID := range nodesInLevel {
			// 查找节点
			for _, node := range nodes {
				if node.ID == nodeID {
					// 计算位置
					x := currentX + float64(i)*(e.nodeWidth+e.horizontalGap)
					y := startY + float64(level)*e.verticalGap

					// 如果该层只有一个节点，居中显示
					if len(nodesInLevel) == 1 {
						x = startX + 200 // 居中偏移
					}

					node.Position = &Position{
						X: x,
						Y: y,
					}
					break
				}
			}
		}
	}
}

// OptimizeLayout 优化布局（减少交叉边）
func (e *AutoLayoutEngine) OptimizeLayout(workflow *ConfiguredWorkflow) {
	// 简单实现：按层排序节点，减少交叉
	// 更复杂的优化可以使用 Sugiyama 算法
	// TODO: 实现更高级的优化算法
}
