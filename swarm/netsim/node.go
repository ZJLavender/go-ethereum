// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package netsim

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

// NodeIDs returns NodeIDs for all nodes in the network.
func (s *Simulation) NodeIDs() (ids []discover.NodeID) {
	nodes := s.Net.GetNodes()
	ids = make([]discover.NodeID, len(nodes))
	for i, node := range nodes {
		ids[i] = node.ID()
	}
	return ids
}

// UpNodeIDs returns NodeIDs for nodeas that are up in the network.
func (s *Simulation) UpNodeIDs() (ids []discover.NodeID) {
	nodes := s.Net.GetNodes()
	for _, node := range nodes {
		if node.Up {
			ids = append(ids, node.ID())
		}
	}
	return ids
}

// AddNodeOption defines the option that can be passed
// to Simulation.AddNode method.
type AddNodeOption func(*adapters.NodeConfig)

// AddNodeWithName sets the optional name for the new node.
func AddNodeWithName(name string) AddNodeOption {
	return func(o *adapters.NodeConfig) {
		o.Name = name
	}
}

// AddNodeWithMsgEvents sets the EnableMsgEvents option
// to NodeConfig.
func AddNodeWithMsgEvents(enable bool) AddNodeOption {
	return func(o *adapters.NodeConfig) {
		o.EnableMsgEvents = enable
	}
}

// AddNode creates a new node with random configuration,
// applies provided options to the config and adds the node to network.
func (s *Simulation) AddNode(opts ...AddNodeOption) (id discover.NodeID, err error) {
	conf := adapters.RandomNodeConfig()
	for _, o := range opts {
		o(conf)
	}
	conf.Services = s.serviceNames
	node, err := s.Net.NewNodeWithConfig(conf)
	if err != nil {
		return id, err
	}
	return node.ID(), s.Net.Start(node.ID())
}

// SetPivotNode sets the NodeID of the network's pivot node.
// Pivot node is just a specific node that should be treated
// differently then other nodes in test. SetPivotNode and
// PivotNodeID are just a convenient functions to set and
// retrieve it.
func (s *Simulation) SetPivotNode(id discover.NodeID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pivotNodeID = &id
}

// PivotNodeID returns NodeID of the pivot node set by
// Simulation.SetPivotNode method.
func (s *Simulation) PivotNodeID() (id *discover.NodeID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pivotNodeID
}

// StopNode stops a node by NodeID.
func (s *Simulation) StopNode(id discover.NodeID) (err error) {
	return s.Net.GetNode(id).Stop()
}

// StopRandomNode stops a random node.
func (s *Simulation) StopRandomNode() (err error) {
	n := s.randomNode()
	if n == nil {
		return ErrNodeNotFound
	}
	return n.Stop()
}

// seed the random generator for Simulation.randomNode.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomNode returns a random SimNode that is up.
// Arguments are NodeIDs for nodes that should not be returned.
func (s *Simulation) randomNode(exclude ...discover.NodeID) *adapters.SimNode {
	ids := s.UpNodeIDs()
	for _, e := range exclude {
		for i, id := range ids {
			if id == e {
				ids = append(ids[:i], ids[i+1:]...)
			}
		}
	}
	n := s.Net.GetNode(ids[rand.Intn(len(ids))])
	node, _ := n.Node.(*adapters.SimNode)
	return node
}
