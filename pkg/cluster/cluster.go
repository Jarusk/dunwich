package cluster

import (
	"net/netip"

	"github.com/hashicorp/memberlist"

	log "github.com/sirupsen/logrus"
)

type Cluster struct {
	ml *memberlist.Memberlist
}

type MemberlistHandler struct{}

func (m MemberlistHandler) NotifyConflict(existing *memberlist.Node, added *memberlist.Node) {
	log.Infof("%s(%s) and %s(%s) joined, are same node", existing, existing.Addr, added, added.Addr)
}

func (c Cluster) JoinCluster(listenPort int, bootstrapNodes []netip.AddrPort) {

	handler := MemberlistHandler{}

	mlConfig := memberlist.DefaultLocalConfig()
	mlConfig.Conflict = handler
	mlConfig.BindPort = listenPort

	list, err := memberlist.Create(mlConfig)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to create memberlist")
	}
	log.Infof("listening on %v for gossip", mlConfig.BindPort)

	c.ml = list

	if len(bootstrapNodes) > 0 {
		log.Debugf("attempting to join via %v", bootstrapNodes)
		_, err := c.ml.Join(nodeListToStringList(bootstrapNodes))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("failed to join gossip")
		}
	} else {
		log.Debug("no join node set, acting as bootstrap")
	}
}

func nodeListToStringList(nodes []netip.AddrPort) []string {
	nodeStringList := []string{}

	for _, v := range nodes {
		v := v
		nodeStringList = append(nodeStringList, v.String())
	}

	return nodeStringList
}
