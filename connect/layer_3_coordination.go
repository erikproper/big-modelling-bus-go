/*
 *
 * Module:    BIG Modelling Bus
 * Package:   Generic
 * Component: Layer 3 - Coordination
 *
 * ..... ... .. .
 *
 * Creator: Henderik A. Proper (e.proper@acm.org), TU Wien, Austria
 *
 * Version of: 05.12.2025
 *
 */

package connect

import "github.com/erikproper/big-modelling-bus.go.v1/generics"

const (
	coordinationPathElement = "coordination"
)

/*
 * Defining topic paths
 */

func (b *TModellingBusConnector) coordinationTopicPath(coordinationID string) string {
	return coordinationPathElement +
		"/" + coordinationID
}

/*
 *
 * Externally visible functionality
 *
 */

/*
 * Posting coordination messages
 */

func (b *TModellingBusConnector) PostCoordination(coordinationID string, json []byte) {
	b.postJSONAsStreamed(b.coordinationTopicPath(coordinationID), json, generics.GetTimestamp())
}

/*
 * Listening to coordination related postings
 */

func (b *TModellingBusConnector) ListenForCoordinationPostings(agentID, coordinationID string, postingHandler func([]byte, string)) {
	b.listenForStreamedPostings(agentID, b.coordinationTopicPath(coordinationID), postingHandler)
}

/*
 * Retrieving coordination messages
 */

func (b *TModellingBusConnector) GetCoordination(agentID, coordinationID string) ([]byte, string) {
	return b.getStreamed(agentID, b.coordinationTopicPath(coordinationID))
}

/*
 * Deleting coordination messages
 */

func (b *TModellingBusConnector) DeleteCoordination(coordinationID string) {
	b.deletePosting(b.coordinationTopicPath(coordinationID))
}
