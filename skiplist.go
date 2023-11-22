package skiplist

import (
	"cmp"
	"fmt"
	"math/bits"
	"math/rand"
)

type Record[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

type Node[K cmp.Ordered, V any] struct {
	record *Record[K, V]
	next   []*Node[K, V]
}

type SkipList[K cmp.Ordered, V any] struct {
	head  *Node[K, V]
	level int
	size  int
}

func (s *SkipList[K, V]) String() string {
	var result string

	for i := s.level; i >= 0; i-- {
		x := s.head.next[i]
		for x != nil {
			result = result + fmt.Sprintf("level: %d - node: %v, %v\n", i, x.record.Key, x.record.Value)
			x = x.next[i]
		}
	}

	return result
}

func NewRecord[K cmp.Ordered, V any](key K, value V) *Record[K, V] {
	return &Record[K, V]{
		Key:   key,
		Value: value,
	}
}

func NewNode[K cmp.Ordered, V any](key K, value V, level int) *Node[K, V] {
	return &Node[K, V]{
		record: NewRecord(key, value),
		next:   make([]*Node[K, V], level+1),
	}
}

func NewHeaderNode[K cmp.Ordered, V any](level int) *Node[K, V] {
	return &Node[K, V]{
		record: nil,
		next:   make([]*Node[K, V], level+1),
	}
}

func NewSkipList[K cmp.Ordered, V any]() *SkipList[K, V] {
	return &SkipList[K, V]{
		head:  NewHeaderNode[K, V](0),
		level: 0,
		size:  0,
	}
}

func (s *SkipList[K, V]) findNode(key K) (*Node[K, V], bool) {
	x := s.head

	// for each level if I found a greater key I pass to a lower level
	// if it's inferior I go to the next node
	// if it's the same I've found the node
	for i := s.level; i >= 0; i-- {
		for {
			if x.next[i] == nil || x.next[i].record.Key > key {
				break
			} else if x.next[i].record.Key == key {
				return x.next[i], true
			} else {
				x = x.next[i]
			}
		}
	}

	return nil, false
}

func (s *SkipList[K, V]) Find(key K) (V, bool) {
	node, found := s.findNode(key)

	if found {
		return node.record.Value, true
	} else {
		return *new(V), false
	}
}

func (s *SkipList[K, V]) getRandomLevel() int {
	// level := 0
	// for rand.Int31()%2 == 0 {
	// 	level += 1
	// }
	return bits.TrailingZeros64(rand.Uint64())
}

func (s *SkipList[K, V]) adjustLevel(level int) {
	s.head.next = append(s.head.next, make([]*Node[K, V], level-s.level)...)
	s.level = level
}

func (s *SkipList[K, V]) Insert(key K, value V) {
	x := s.head
	updates := make([]*Node[K, V], s.level+1)

	// for each level, store in updates the nodes that should be just before the new node
	// in the loop, search will start from the last found node, no need to restart from head
	for i := s.level; i >= 0; i-- {
		for x.next[i] != nil && x.next[i].record.Key < key {
			x = x.next[i]
		}
		updates[i] = x
	}

	// replace value if key already present
	if x.next[0] != nil && x.next[0].record.Key == key {
		x.next[0].record.Value = value
	} else {
		newLevel := s.getRandomLevel()

		// add new level
		if newLevel > s.level {
			for i := s.level; i < newLevel; i++ {
				updates = append(updates, s.head)
			}
			s.adjustLevel(newLevel)
		}

		newNode := NewNode[K, V](key, value, newLevel)

		// node references swap: the already stored nodes with inferior keys will point to the new node
		// the new node will point to what the already stored nodes pointed before
		for i := 0; i <= newLevel; i++ {
			newNode.next[i] = updates[i].next[i]
			updates[i].next[i] = newNode
		}

		s.size += 1
	}
}

func (s *SkipList[K, V]) Delete(key K) {
	x := s.head

	// similar to find, but when I found the node I update the reference to the node pointed
	// by the node to be deleted for each level
	for i := s.level; i >= 0; i-- {
		for {
			if x.next[i] == nil || x.next[i].record.Key > key {
				break
			} else if x.next[i].record.Key == key {
				x.next[i] = x.next[i].next[i]
			} else {
				x = x.next[i]
			}
		}
	}
}
