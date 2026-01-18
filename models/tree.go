package models

// NodeType indicates whether a node is a category (folder) or module (leaf).
type NodeType int

const (
	NodeCategory NodeType = iota // Expandable folder
	NodeModule                   // Leaf node (installable)
)

// TreeNode represents a single node in the tree hierarchy.
type TreeNode struct {
	ID       string      // Path-based ID (e.g., "claude/workflow/tdd_guard")
	Name     string      // Display name
	Type     NodeType    // Category or Module
	Depth    int         // Nesting level (0 = root category)
	Expanded bool        // Whether children are visible (categories only)
	Children []*TreeNode // Child nodes
	Module   Module      // Non-nil for module nodes
	Parent   *TreeNode   // Parent node (nil for root)
}

// TreeState holds the complete tree state for the TUI.
type TreeState struct {
	Roots       []*TreeNode     // Top-level category nodes
	FlatVisible []*TreeNode     // Flattened list of visible nodes for navigation
	Expanded    map[string]bool // Track expanded state by node ID
}

// NewTreeState creates a new empty tree state.
func NewTreeState() *TreeState {
	return &TreeState{
		Roots:       make([]*TreeNode, 0),
		FlatVisible: make([]*TreeNode, 0),
		Expanded:    make(map[string]bool),
	}
}

// IsExpanded returns whether a node is expanded.
func (t *TreeState) IsExpanded(id string) bool {
	return t.Expanded[id]
}

// SetExpanded sets the expanded state of a node.
func (t *TreeState) SetExpanded(id string, expanded bool) {
	t.Expanded[id] = expanded
}

// ToggleExpanded toggles the expanded state of a node.
func (t *TreeState) ToggleExpanded(id string) {
	t.Expanded[id] = !t.Expanded[id]
}

// RebuildFlatVisible rebuilds the flattened list of visible nodes.
func (t *TreeState) RebuildFlatVisible() {
	t.FlatVisible = make([]*TreeNode, 0)
	for _, root := range t.Roots {
		t.flattenNode(root)
	}
}

func (t *TreeState) flattenNode(node *TreeNode) {
	t.FlatVisible = append(t.FlatVisible, node)
	if node.Type == NodeCategory && t.IsExpanded(node.ID) {
		for _, child := range node.Children {
			t.flattenNode(child)
		}
	}
}

// GetInstalledCount returns (installed, total) module counts for a category node.
func (node *TreeNode) GetInstalledCount() (int, int) {
	if node.Type == NodeModule {
		if node.Module != nil && node.Module.IsInstalled() {
			return 1, 1
		}
		return 0, 1
	}

	installed, total := 0, 0
	for _, child := range node.Children {
		ci, ct := child.GetInstalledCount()
		installed += ci
		total += ct
	}
	return installed, total
}

// HasModules returns true if the node contains any modules (directly or in descendants).
func (node *TreeNode) HasModules() bool {
	if node.Type == NodeModule {
		return true
	}
	for _, child := range node.Children {
		if child.HasModules() {
			return true
		}
	}
	return false
}

// IsLastChild returns true if this node is the last child of its parent.
func (node *TreeNode) IsLastChild() bool {
	if node.Parent == nil {
		return true
	}
	children := node.Parent.Children
	return len(children) > 0 && children[len(children)-1] == node
}
