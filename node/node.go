package node

type Node struct {
	NodeID string
	NodeIP string
	Type   NodeType
}

type NodeType int

const (
	WHOCARES string = ""

	Client NodeType = iota + 1
	Proxy
)

// String - Creating common behavior - give the type a String function
func (w NodeType) String() string {
	return [...]string{"Client", "Proxy"}[w-1]
}

func NewNode() *Node {
	return &Node{
		Type:   Client,
		NodeIP: WHOCARES,
	}
}
