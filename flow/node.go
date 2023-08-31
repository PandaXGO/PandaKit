package flow

import (
	"github.com/PandaXGO/PandaKit/utils"
)

type Properties map[string]any

type Node struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	X          int        `json:"x"`
	Y          int        `json:"y"`
	Text       Text       `json:"text"`
	Properties Properties `json:"properties"`
}

type Text struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Value string `json:"value"`
}

func (node *Node) IsStartNode(ty string) bool {
	if node.Type == ty {
		return true
	}
	return false
}

func (node *Node) GetProperties(data any) error {
	if err := utils.Map2Struct(node.Properties, data); err != nil {
		return err
	}
	return nil
}

type NodeFunc func(*Node)

func (node *Node) RunNodeFunc(nodeFunc NodeFunc) {
	nodeFunc(node)
}
