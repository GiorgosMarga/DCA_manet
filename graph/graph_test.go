package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeights(t *testing.T) {
	g := NewGraph()
	assert.NotNil(t, g)
	err := g.AddNode(10)
	assert.Nil(t, err)
	assert.NotNil(t, g.AddNode(10))
}

func TestConnectNodesInvalidIds(t *testing.T) {
	g := NewGraph()
	assert.NotNil(t, g)
	err := g.AddNode(10)
	assert.Nil(t, err)
	err = g.AddNode(15)
	assert.Nil(t, err)
	err = g.ConnectNodes(2, 1)
	assert.Error(t, err)
}
func TestConnectNodes(t *testing.T) {
	g := NewGraph()
	assert.NotNil(t, g)
	err := g.AddNode(10)
	assert.Nil(t, err)
	err = g.AddNode(15)
	assert.Nil(t, err)
	err = g.ConnectNodes(1, 2)
	assert.Nil(t, err)
}
