package lru

type bintreeLRUItem struct {
	key             string
	value           interface{}
	hits            int
	parent          *bintreeLRUItem
	left            *bintreeLRUItem
	right           *bintreeLRUItem
	morePopularNode *bintreeLRUItem
	lessPopularNode *bintreeLRUItem
}

type bintreeLRU struct {
	tip            *bintreeLRUItem
	capacity       int
	size           int
	popularityTail *bintreeLRUItem
}

// NewBintreeLRU creates an instance of the LRU cache with the binary tree as a backend
func NewBintreeLRU(capacity int) LRU {
	if capacity <= 0 {
		capacity = 1
	}

	return &bintreeLRU{capacity: capacity, size: 0}
}

func (l *bintreeLRU) Get(key string) (found bool, value interface{}) {
	tip := l.tip
	for {
		if tip == nil {
			return false, nil
		}

		if tip.key == key {
			tip.hits++
			l.swap(tip)

			return true, tip.value
		}

		if tip.key > key {
			tip = tip.left
		} else {
			tip = tip.right
		}
	}
}

func (l *bintreeLRU) Set(key string, value interface{}) {
	if l.tip == nil {
		l.tip = &bintreeLRUItem{
			key:             key,
			value:           value,
			hits:            1,
			parent:          nil,
			left:            nil,
			right:           nil,
			morePopularNode: nil,
			lessPopularNode: nil,
		}

		l.size = 1
		l.popularityTail = l.tip
		return
	}

	node := l.tip
	for {
		switch {
		case node.key == key:
			node.hits++
			l.swap(node)
			return

		case node.key > key:
			if node.left != nil {
				node = node.left
				continue
			}

			if l.size == l.capacity {
				l.evict()
				l.Set(key, value)
				return
			}

			// add
			node.left = l.newNode(key, value)
			node.left.parent = node

			// rebalance
			l.rebalance(node)

			l.size++
			return

		case node.key < key:
			if node.right != nil {
				node = node.right
				continue
			}

			if l.size == l.capacity {
				l.evict()
				l.Set(key, value)
				return
			}

			// add
			node.right = l.newNode(key, value)
			node.right.parent = node

			// rebalance
			l.rebalance(node)
			l.size++
			return

		}
	}
}

func (l *bintreeLRU) Size() int {
	return l.size
}

func (l *bintreeLRU) swap(node *bintreeLRUItem) {
	for {
		if node == nil || node.morePopularNode == nil || node.hits <= node.morePopularNode.hits {
			return
		}

		currentNextNode := node.morePopularNode
		nextNextNode := currentNextNode.morePopularNode

		currentNextNode.morePopularNode = node
		currentNextNode.lessPopularNode = node.lessPopularNode
		node.morePopularNode = nextNextNode
		node.lessPopularNode = currentNextNode

		if nextNextNode != nil {
			nextNextNode.lessPopularNode = node
		}

		if currentNextNode.lessPopularNode != nil {
			currentNextNode.lessPopularNode.morePopularNode = currentNextNode
		}

		if l.popularityTail == node {
			l.popularityTail = currentNextNode
		}
	}
}

func (l *bintreeLRU) findBiggestInSubtree(node *bintreeLRUItem) *bintreeLRUItem {
	tip := node
	for {
		if tip.right == nil {
			return tip
		}

		tip = tip.right
	}
}

func (l *bintreeLRU) evict() {
	if l.size < l.capacity || l.popularityTail == nil {
		return
	}

	if l.size == 1 {
		l.tip = nil
		l.popularityTail = nil
		l.size = 0
		return
	}

	nodeToDelete := l.popularityTail
	var newParent *bintreeLRUItem

	if l.tip.key == nodeToDelete.key {
		// remove the tip
		if l.tip.left != nil {
			newParent = l.findBiggestInSubtree(l.tip.left)

			newParentsParent := newParent.parent
			if newParentsParent != l.tip {
				newParentsParent.right = newParent.left
				if newParentsParent.right != nil {
					newParentsParent.right.parent = newParentsParent
				}

				newParent.left = l.tip.left
				newParent.left.parent = newParent
			}

			newParent.right = l.tip.right
			if newParent.right != nil {
				newParent.right.parent = newParent
			}
		} else {
			newParent = l.tip.right
		}

		l.popularityTail = l.popularityTail.morePopularNode
		if l.popularityTail != nil {
			l.popularityTail.lessPopularNode = nil
		}

		l.tip = newParent
		l.tip.parent = nil
	} else if nodeToDelete.left == nil && nodeToDelete.right == nil {
		// this is a leaf node, simply unlink it from parent

		if nodeToDelete.parent.left == nodeToDelete {
			nodeToDelete.parent.left = nil
		} else {
			nodeToDelete.parent.right = nil
		}

		l.popularityTail = l.popularityTail.morePopularNode
		if l.popularityTail != nil {
			l.popularityTail.lessPopularNode = nil
		}

		newParent = nodeToDelete.parent

	} else {
		// this isn't a tip or a leaf node. we need to find a new parent for the remaining nodes

		parentNode := nodeToDelete.parent
		if parentNode == nil {
			panic("parentNode can't be nil, something is wrong with the tree")
		}

		if nodeToDelete.left != nil {
			newParent = l.findBiggestInSubtree(nodeToDelete.left)
			newParentsParent := newParent.parent
			if newParentsParent != nodeToDelete {
				newParentsParent.right = newParent.left
				if newParentsParent.right != nil {
					newParentsParent.right.parent = newParentsParent
				}

				newParent.left = nodeToDelete.left
				newParent.left.parent = newParent
			}

			newParent.right = nodeToDelete.right
			if newParent.right != nil {
				newParent.right.parent = newParent
			}
		} else {
			newParent = nodeToDelete.right
		}

		l.popularityTail = l.popularityTail.morePopularNode
		if l.popularityTail != nil {
			l.popularityTail.lessPopularNode = nil
		}

		if parentNode.left == nodeToDelete {
			parentNode.left = newParent
			parentNode.left.parent = parentNode
		} else {
			parentNode.right = newParent
			parentNode.right.parent = parentNode
		}
	}

	l.size--
	l.rebalance(newParent)
}

func (l *bintreeLRU) newNode(key string, value interface{}) *bintreeLRUItem {
	newNode := &bintreeLRUItem{
		key:             key,
		value:           value,
		hits:            1,
		left:            nil,
		right:           nil,
		morePopularNode: l.popularityTail,
		lessPopularNode: nil,
	}

	l.popularityTail = newNode
	if newNode.morePopularNode != nil {
		newNode.morePopularNode.lessPopularNode = newNode
	}

	return newNode
}

func (l *bintreeLRU) maxDepth(node *bintreeLRUItem) int {
	if node == nil {
		return 0
	}

	leftDepth := l.maxDepth(node.left)
	rightDepth := l.maxDepth(node.right)

	if leftDepth > rightDepth {
		return leftDepth + 1
	}

	return rightDepth + 1
}

func (l *bintreeLRU) rebalance(node *bintreeLRUItem) {
	if node == nil {
		return
	}

	leftDepth := l.maxDepth(node.left)
	rightDepth := l.maxDepth(node.right)

	nodeParent := node.parent

	if leftDepth > rightDepth+1 {
		newParent := node.left
		newParent.parent = nodeParent

		node.left = newParent.right
		if node.left != nil {
			node.left.parent = node
		}

		newParent.right = node
		newParent.right.parent = newParent

		if nodeParent == nil {
			l.tip = newParent
			l.tip.parent = nil

		} else if nodeParent.left == node {
			nodeParent.left = newParent
			nodeParent.left.parent = nodeParent

		} else if nodeParent.right == node {
			nodeParent.right = newParent
			nodeParent.right.parent = nodeParent
		}

	} else if rightDepth > leftDepth+1 {
		newParent := node.right
		newParent.parent = nodeParent

		node.right = newParent.left
		if node.right != nil {
			node.right.parent = node
		}

		newParent.left = node
		newParent.left.parent = newParent

		if nodeParent == nil {
			l.tip = newParent
			l.tip.parent = nil

		} else if nodeParent.left == node {
			nodeParent.left = newParent
			nodeParent.left.parent = nodeParent

		} else if nodeParent.right == node {
			nodeParent.right = newParent
			nodeParent.right.parent = nodeParent
		}
	}

	l.rebalance(nodeParent)
}
