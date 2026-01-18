package helpers

import (
	"sort"
	"strings"

	"code-template/models"
)

// BuildTree creates a tree structure from a list of modules.
// Modules are organized by their GetPath() which uses "/" as separator.
func BuildTree(modules []models.Module) *models.TreeState {
	tree := models.NewTreeState()

	// Map to track category nodes by their full path
	categoryNodes := make(map[string]*models.TreeNode)

	for _, module := range modules {
		path := module.GetPath()
		parts := strings.Split(path, "/")

		// Build intermediate category nodes
		var parent *models.TreeNode
		var currentPath string

		for i, part := range parts[:len(parts)-1] {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = currentPath + "/" + part
			}

			if existing, ok := categoryNodes[currentPath]; ok {
				parent = existing
				continue
			}

			node := &models.TreeNode{
				ID:       currentPath,
				Name:     part,
				Type:     models.NodeCategory,
				Depth:    i,
				Expanded: false,
				Children: make([]*models.TreeNode, 0),
				Parent:   parent,
			}
			categoryNodes[currentPath] = node

			if parent == nil {
				tree.Roots = append(tree.Roots, node)
			} else {
				parent.Children = append(parent.Children, node)
			}
			parent = node
		}

		// Create the module node (leaf)
		moduleNode := &models.TreeNode{
			ID:     path,
			Name:   module.GetName(),
			Type:   models.NodeModule,
			Depth:  len(parts) - 1,
			Module: module,
			Parent: parent,
		}

		if parent == nil {
			// Module at root level (shouldn't happen with proper paths)
			tree.Roots = append(tree.Roots, moduleNode)
		} else {
			parent.Children = append(parent.Children, moduleNode)
		}
	}

	// Sort children at each level
	sortNodes(tree.Roots)
	for _, node := range categoryNodes {
		sortNodes(node.Children)
	}

	// Filter out empty categories
	tree.Roots = filterEmptyCategories(tree.Roots)

	tree.RebuildFlatVisible()
	return tree
}

// sortNodes sorts nodes alphabetically, with categories before modules.
func sortNodes(nodes []*models.TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		// Categories come before modules
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type == models.NodeCategory
		}
		return nodes[i].Name < nodes[j].Name
	})
}

// filterEmptyCategories removes categories that have no modules.
func filterEmptyCategories(nodes []*models.TreeNode) []*models.TreeNode {
	result := make([]*models.TreeNode, 0)
	for _, node := range nodes {
		if node.Type == models.NodeModule {
			result = append(result, node)
		} else if node.HasModules() {
			node.Children = filterEmptyCategories(node.Children)
			result = append(result, node)
		}
	}
	return result
}

// GetModuleList returns a flat list of all modules registered with the system.
func GetModuleList() []models.Module {
	return moduleRegistry
}

// moduleRegistry is the list of all registered modules.
var moduleRegistry []models.Module

// RegisterModule adds a module to the registry.
func RegisterModule(m models.Module) {
	moduleRegistry = append(moduleRegistry, m)
}

// ClearRegistry clears all registered modules (useful for testing).
func ClearRegistry() {
	moduleRegistry = nil
}
