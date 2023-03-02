package service

type INode interface {
	GetId() int
	GetParent() int
	GetName() string
	GetSort() int
	GetCode() int
	GetType() int
	GetRoot() string
	IsRoot() bool
	GetData() interface{}
}
type INodes []INode

type SystemMenu struct {
	Id     int    `json:"id"`
	Parent int    `json:"parent"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Code   int    `json:"code"`
	Type   int    `json:"type"`
	Root   string `json:"root"`
}

func (m SystemMenu) GetId() int {
	return m.Id
}
func (m SystemMenu) GetParent() int {
	return m.Parent
}
func (m SystemMenu) GetName() string {
	return m.Name
}
func (m SystemMenu) GetSort() int {
	return m.Sort
}
func (m SystemMenu) GetCode() int {
	return m.Code
}
func (m SystemMenu) GetType() int {
	return m.Type
}
func (m SystemMenu) GetRoot() string {
	return m.Root
}
func (m SystemMenu) GetData() interface{} {
	return m
}
func (m SystemMenu) IsRoot() bool {
	// 这里通过FatherId等于0 或者 FatherId等于自身Id表示顶层根节点
	return m.Parent == 0 || m.Parent == m.Id
}

type SystemMenus []SystemMenu

func (s SystemMenus) ConvertToINodeArray() (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}
	return
}

type Tree struct {
	Title           string      `json:"title"`            //节点名字
	Data            interface{} `json:"data"`             //自定义对象
	Leaf            bool        `json:"leaf"`             //叶子节点
	Selected        bool        `json:"checked"`          //选中
	PartialSelected bool        `json:"partial_selected"` //部分选中
	Children        []Tree      `json:"children"`         //子节点
}

func GenerateTree(nodes, selectedNodes []INode) (trees []Tree) {
	trees = []Tree{}
	// 定义顶层根和子节点
	var roots, childs []INode
	for _, v := range nodes {
		if v.IsRoot() {
			// 判断顶层根节点
			roots = append(roots, v)
		}
		childs = append(childs, v)
	}

	for _, v := range roots {
		childTree := &Tree{
			Title: v.GetName(),
			Data:  v.GetData(),
		}
		// 递归之前，根据父节点确认 childTree 的选中状态
		childTree.Selected = nodeSelected(v, selectedNodes, childTree.Children)
		// 递归
		recursiveTree(childTree, childs, selectedNodes)
		// 递归之后，根据子节点确认 childTree 的选中状态
		if !childTree.Selected {
			childTree.Selected = nodeSelected(v, selectedNodes, childTree.Children)
		}
		// 递归之后，根据子节点确认 childTree 的半选中状态
		childTree.PartialSelected = nodePartialSelected(childTree.Children)
		// 递归之后，根据子节确认是否是叶子节点
		childTree.Leaf = len(childTree.Children) == 0
		trees = append(trees, *childTree)
	}
	return
}

func nodeSelected(node INode, selectedNodes []INode, children []Tree) bool {
	for _, v := range selectedNodes {
		if node.GetId() == v.GetId() {
			// 1. 如果选择节点数组中存在当前节点
			return true
		}
	}

	if len(children) == 0 {
		// 2. 不满足前置条件1，且没有子节点
		return false
	}
	selectedNum := 0
	for _, v := range children {
		if v.Selected {
			selectedNum++
		}
	}
	if selectedNum == len(children) {
		// 不满足前置条件1，2 ，且子节点全部是选中状态
		return true
	}
	return false
}

// nodePartialSelected 判断节点的半选中状态
func nodePartialSelected(trees []Tree) bool {
	selectedNum := 0
	for _, v := range trees {
		if v.Selected {
			selectedNum++
		}
	}
	if selectedNum == len(trees) || selectedNum == 0 {
		// 子节点全选中，或一个也没有选中
		return false
	}
	return true
}

// recursiveFindRelationNode 递归查询关联父子节点
// nodeMap 查询结果搜集到map中
// allNodes 所有节点
// node 递归节点
// t 递归查找类型：0 查找父、子节点；1 只查找父节点；2 只查找子节点
func recursiveFindRelationNode(nodeMap map[int]INode, allNodes []INode, node INode, t int) {
	nodeMap[node.GetId()] = node
	for _, v := range allNodes {
		if _, ok := nodeMap[v.GetId()]; ok {
			continue
		}
		// 查找父节点
		if t == 0 || t == 1 {
			if node.GetParent() == v.GetId() {
				nodeMap[v.GetId()] = v
				if v.IsRoot() {
					// 是顶层根节点时，不再进行递归
					continue
				}
				recursiveFindRelationNode(nodeMap, allNodes, v, 1)
			}
		}
		// 查找子节点
		if t == 0 || t == 2 {
			if node.GetId() == v.GetParent() {
				nodeMap[v.GetId()] = v
				recursiveFindRelationNode(nodeMap, allNodes, v, 2)
			}
		}
	}
}

// FindRelationNode 在 allTree 中查询 nodes 中节点的所有父节点
// nodes 要查询父节点的子节点数组
// allTree 所有节点数组
func FindRelationNode(nodes, allNodes []INode) (respNodes []INode) {
	nodeMap := make(map[int]INode)
	for _, v := range nodes {
		recursiveFindRelationNode(nodeMap, allNodes, v, 0)
	}

	for _, v := range nodeMap {
		respNodes = append(respNodes, v)
	}
	//sort.Sort(INodes(respNodes))
	return
}

// recursiveTree 递归生成树结构
// tree 递归的树对象
// nodes 递归的节点
// selectedNodes 选中的节点
func recursiveTree(tree *Tree, nodes, selectedNodes []INode) {
	data := tree.Data.(INode)

	for _, v := range nodes {
		if v.IsRoot() {
			// 如果当前节点是顶层根节点就跳过
			continue
		}
		if data.GetId() == v.GetParent() {
			childTree := &Tree{
				Title: v.GetName(),
				Data:  v.GetData(),
			}
			// 递归之前，根据子节点和父节点确认 childTree 的选中状态
			childTree.Selected = nodeSelected(v, selectedNodes, childTree.Children) || tree.Selected
			recursiveTree(childTree, nodes, selectedNodes)

			if !childTree.Selected {
				// 递归之后，根据子节点确认 childTree 的选中状态
				childTree.Selected = nodeSelected(v, selectedNodes, childTree.Children)
			}
			// 递归之后，根据子节点确认 childTree 的半选中状态
			childTree.PartialSelected = nodePartialSelected(childTree.Children)
			// 递归之后，根据子节确认是否是叶子节点
			childTree.Leaf = len(childTree.Children) == 0
			tree.Children = append(tree.Children, *childTree)
		}
	}
}
