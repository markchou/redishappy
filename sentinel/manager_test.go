package sentinel

import (
	"testing"

	"github.com/mdevilliers/redishappy/configuration"
	"github.com/mdevilliers/redishappy/types"
)

func TestBasicEventChannel(t *testing.T) {

	switchmasterchannel := make(chan types.MasterSwitchedEvent)

	manager := NewManager(switchmasterchannel, configuration.NewConfigurationManager(configuration.Configuration{}))
	defer manager.ClearState()
	manager.Notify(&SentinelAdded{Sentinel: types.Sentinel{Host: "10.1.1.1", Port: 12345}})

	responseChannel := make(chan SentinelTopology)

	manager.GetState(TopologyRequest{ReplyChannel: responseChannel})
	topologyState := <-responseChannel

	if len(topologyState.Sentinels) != 1 {
		t.Error("Topology count should be 1")
	}

	manager2 := NewManager(switchmasterchannel, configuration.NewConfigurationManager(configuration.Configuration{}))
	manager2.Notify(&SentinelAdded{Sentinel: types.Sentinel{Host: "10.1.1.2", Port: 12345}})

	manager2.GetState(TopologyRequest{ReplyChannel: responseChannel})

	topologyState = <-responseChannel

	if len(topologyState.Sentinels) != 2 {
		t.Errorf("Topology count should be 2 : it is %d", len(topologyState.Sentinels))
	}
}

func TestAddingAndLoseingASentinel(t *testing.T) {

	switchmasterchannel := make(chan types.MasterSwitchedEvent)
	manager := NewManager(switchmasterchannel, configuration.NewConfigurationManager(configuration.Configuration{}))
	defer manager.ClearState()

	sentinel := types.Sentinel{Host: "10.1.1.5", Port: 12345}

	manager.Notify(&SentinelAdded{Sentinel: sentinel})
	manager.Notify(&SentinelLost{Sentinel: sentinel})

	responseChannel := make(chan SentinelTopology)

	manager.GetState(TopologyRequest{ReplyChannel: responseChannel})
	topologyState := <-responseChannel

	if len(topologyState.Sentinels) != 1 {
		t.Error("Topology count should be 1")
	}
}

func TestAddingSentinelMultipleTimes(t *testing.T) {

	switchmasterchannel := make(chan types.MasterSwitchedEvent)
	manager := NewManager(switchmasterchannel, configuration.NewConfigurationManager(configuration.Configuration{}))
	defer manager.ClearState()

	sentinel := types.Sentinel{Host: "10.1.1.6", Port: 12345}

	manager.Notify(&SentinelAdded{Sentinel: sentinel})

	ping := &SentinelPing{Sentinel: sentinel}
	ping2 := &SentinelPing{Sentinel: sentinel}
	manager.Notify(ping)
	manager.Notify(ping2)

	responseChannel := make(chan SentinelTopology)

	manager.GetState(TopologyRequest{ReplyChannel: responseChannel})
	topologyState := <-responseChannel

	_, ok := topologyState.FindSentinelInfo(sentinel)

	if !ok {
		t.Error("Added sentinel not found")
	}
}
