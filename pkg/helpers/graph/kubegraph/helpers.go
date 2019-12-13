package kubegraph

import (
	"sort"

	"github.com/openshift/oc/pkg/cli/helpers/revision"
	osgraph "github.com/openshift/oc/pkg/helpers/graph/genericgraph"
	kubegraph "github.com/openshift/oc/pkg/helpers/graph/kubegraph/nodes"
)

// RelevantDeployments returns the active deployment and a list of inactive deployments (in order from newest to oldest)
func RelevantDeployments(g osgraph.Graph, dNode *kubegraph.DeploymentNode) (*kubegraph.ReplicaSetNode, []*kubegraph.ReplicaSetNode) {
	allDeployments := []*kubegraph.ReplicaSetNode{}
	uncastDeployments := g.SuccessorNodesByEdgeKind(dNode, DeploymentEdgeKind)
	if len(uncastDeployments) == 0 {
		return nil, []*kubegraph.ReplicaSetNode{}
	}

	for i := range uncastDeployments {
		allDeployments = append(allDeployments, uncastDeployments[i].(*kubegraph.ReplicaSetNode))
	}

	sort.Sort(RecentDeploymentReferences(allDeployments))

	deploymentRevision, _ := revision.Revision(dNode.Deployment)
	firstRSRevision, _ := revision.Revision(allDeployments[0].ReplicaSet)

	if deploymentRevision == firstRSRevision {
		return allDeployments[0], allDeployments[1:]
	}

	return nil, allDeployments
}

type RecentDeploymentReferences []*kubegraph.ReplicaSetNode

func (m RecentDeploymentReferences) Len() int      { return len(m) }
func (m RecentDeploymentReferences) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m RecentDeploymentReferences) Less(i, j int) bool {
	firstRev, _ := revision.Revision(m[i].ReplicaSet)
	secondRev, _ := revision.Revision(m[j].ReplicaSet)
	return firstRev > secondRev
}
