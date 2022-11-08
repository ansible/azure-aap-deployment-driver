package templates

import (
	"errors"
)

// A node in this graph is just a string, so a nodeset is a map whose
// keys are the nodes that are present.
type nodeset map[string]struct{}

// depmap tracks the nodes that have some dependency relationship to
// some other node, represented by the key of the map.
type depmap map[string]nodeset

type Graph struct {
	nodes nodeset

	// `dependencies` tracks child -> parents.
	dependencies depmap
	// `dependents` tracks parent -> children.
	dependents depmap
	// Keep track of the nodes of the graph themselves.
}

func NewDependencyGraph() *Graph {
	return &Graph{
		dependencies: make(depmap),
		dependents:   make(depmap),
		nodes:        make(nodeset),
	}
}

func (g *Graph) AddDependency(child, parent string) error {
	if child == parent {
		return errors.New("self-referential dependencies not allowed")
	}

	if g.hasDependencyOn(parent, child) {
		return errors.New("circular dependencies not allowed")
	}

	// Add nodes.
	g.nodes[parent] = struct{}{}
	g.nodes[child] = struct{}{}

	// Add edges.
	addNodeToNodeset(g.dependents, parent, child)
	addNodeToNodeset(g.dependencies, child, parent)

	return nil
}

// GetAllDependenciesSorted returns a slice of all of the graph nodes in topological sort order. That is,
// if `B` depends on `A`, then `A` is guaranteed to come before `B` in the sorted output.
// The graph is guaranteed to be cycle-free because cycles are detected while building the
// graph.
// The output is grouped into "layers", which are guaranteed to not have any dependencies within each layer.
func (g *Graph) GetAllDependenciesSorted() [][]string {
	layers := [][]string{}

	// Copy the graph
	shrinkingGraph := g.clone()
	for {
		leaves := shrinkingGraph.leafs()
		if len(leaves) == 0 {
			break
		}

		layers = append(layers, leaves)
		for _, leafNode := range leaves {
			shrinkingGraph.remove(leafNode)
		}
	}

	return layers
}

func (g *Graph) hasDependencyOn(child, parent string) bool {
	deps := g.getAllDependencies(child)
	_, ok := deps[parent]
	return ok
}

func (g *Graph) getAllDependencies(child string) nodeset {
	return g.buildTransitive(child, g.immediateDependencies)
}

func (g *Graph) immediateDependencies(node string) nodeset {
	return g.dependencies[node]
}

func (g *Graph) leafs() []string {
	leafs := make([]string, 0)

	for node := range g.nodes {
		if _, ok := g.dependencies[node]; !ok {
			leafs = append(leafs, node)
		}
	}

	return leafs
}

func (g *Graph) remove(node string) {
	// Remove edges from things that depend on `node`.
	for dependent := range g.dependents[node] {
		removeFromDepmap(g.dependencies, dependent, node)
	}
	delete(g.dependents, node)

	// Remove all edges from node to the things it depends on.
	for dependency := range g.dependencies[node] {
		removeFromDepmap(g.dependents, dependency, node)
	}
	delete(g.dependencies, node)

	// Finally, remove the node itself.
	delete(g.nodes, node)
}

func (g *Graph) clone() *Graph {
	return &Graph{
		dependencies: copyDepmap(g.dependencies),
		dependents:   copyDepmap(g.dependents),
		nodes:        copyNodeset(g.nodes),
	}
}

// buildTransitive starts at `root` and continues calling `nextFn` to keep discovering more nodes until
// the graph cannot produce any more. It returns the set of all discovered nodes.
func (g *Graph) buildTransitive(root string, nextFn func(string) nodeset) nodeset {
	if _, ok := g.nodes[root]; !ok {
		return nil
	}

	out := make(nodeset)
	searchNext := []string{root}
	for len(searchNext) > 0 {
		// List of new nodes from this layer of the dependency graph. This is
		// assigned to `searchNext` at the end of the outer "discovery" loop.
		discovered := []string{}
		for _, node := range searchNext {
			// For each node to discover, find the next nodes.
			for nextNode := range nextFn(node) {
				// If we have not seen the node before, add it to the output as well
				// as the list of nodes to traverse in the next iteration.
				if _, ok := out[nextNode]; !ok {
					out[nextNode] = struct{}{}
					discovered = append(discovered, nextNode)
				}
			}
		}
		searchNext = discovered
	}

	return out
}

func removeFromDepmap(dm depmap, key, node string) {
	nodes := dm[key]
	if len(nodes) == 1 {
		// The only element in the nodeset must be `node`, so we
		// can delete the entry entirely.
		delete(dm, key)
	} else {
		// Otherwise, remove the single node from the nodeset.
		delete(nodes, node)
	}
}

func copyNodeset(s nodeset) nodeset {
	out := make(nodeset, len(s))
	for k, v := range s {
		out[k] = v
	}
	return out
}

func copyDepmap(m depmap) depmap {
	out := make(depmap, len(m))
	for k, v := range m {
		out[k] = copyNodeset(v)
	}
	return out
}

func addNodeToNodeset(dm depmap, key, node string) {
	nodes, ok := dm[key]
	if !ok {
		nodes = make(nodeset)
		dm[key] = nodes
	}
	nodes[node] = struct{}{}
}
