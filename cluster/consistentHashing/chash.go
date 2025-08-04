package chash

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
)

// hashString hashes a string into a uint32
func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type HashRing struct {
	replicas   int
	nodeMap    map[uint32]string
	nodeHashes []uint32
}

// NewHashRing builds the ring with virtual nodes
func NewHashRing(nodes []string, replicas int) *HashRing {
	hr := &HashRing{
		replicas:   replicas,
		nodeMap:    make(map[uint32]string),
		nodeHashes: []uint32{},
	}

	for _, node := range nodes {
		for i := 0; i < replicas; i++ {
			virtualNode := node + "#" + strconv.Itoa(i)
			h := hashString(virtualNode)
			hr.nodeHashes = append(hr.nodeHashes, h)
			hr.nodeMap[h] = node
		}
	}

	sort.Slice(hr.nodeHashes, func(i, j int) bool {
		return hr.nodeHashes[i] < hr.nodeHashes[j]
	})

	return hr
}

// GetNode finds the node for a given key
func (hr *HashRing) GetNode(key string) string {
	h := hashString(key)

	idx := sort.Search(len(hr.nodeHashes), func(i int) bool {
		return hr.nodeHashes[i] >= h
	})

	if idx == len(hr.nodeHashes) {
		idx = 0
	}

	return hr.nodeMap[hr.nodeHashes[idx]]
}

func main() {
	nodes := []string{"node1", "node2", "node3"}
	replicas := 100
	hr := NewHashRing(nodes, replicas)

	keys := []string{ "apple_0", "banana_1", "cherry_2", "date_3", "fig_4", "grape_5",
  "apple_6", "banana_7", "cherry_8", "date_9", "fig_10", "grape_11",
  "apple_12", "banana_13", "cherry_14", "date_15", "fig_16", "grape_17",
  "apple_18", "banana_19", "cherry_20", "date_21", "fig_22", "grape_23",
  "apple_24", "banana_25", "cherry_26", "date_27", "fig_28", "grape_29" }

	for _, key := range keys {
		node := hr.GetNode(key)
		fmt.Printf("Key %q is assigned to node %q\n", key, node)
	}
}
