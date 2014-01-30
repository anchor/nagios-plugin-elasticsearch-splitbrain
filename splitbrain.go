/* 
 * check_elasticsearch_splitbrain
 *
 * Nagios check for splitbrain in an Elasticsearch cluster. Takes one
 * mandatory option, -nodes, followed by a comma-separated node list.
*/

package main

import (
	"github.com/fractalcat/nagiosplugin"
	"os/exec"
	"flag"
	"strings"
	"fmt"
	"regexp"
)

// Given a node's FQDN, grab its topology (via nrpe) and write it to 
// a channel. 
func getTopology(node string, c chan string)  {
	cmd := exec.Command(
		"/usr/lib64/nagios/plugins/check_nrpe",
		"-H",
		node,
		"-u",
		"-c",
		"check_elasticsearch_topology",
		"-t",
		"40",
	)
	topology, err := cmd.CombinedOutput()
	if err != nil {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("Could not get topology for %v: %v", node, err))
	}
	c <- string(topology)
}

// Given a topology and a list of nodes intended to be in the cluster
// (for sanity-checking), returns the FQDN of the master of the 
// topology.
func getMaster(topology string, nodes []string) string {
	nodeMap := make(map[string]bool)
	for _, node := range nodes {
		nodeMap[node] = true
	}
	topologyLines := strings.Split(topology, "\n")
	for _, nodeLine := range topologyLines {
		nameBytes := []byte(strings.Trim(nodeLine, " \n"))
		if len(nameBytes) < 2 {
			continue
		}
		name := string(nodeLine[2:])
		if _, ok := nodeMap[name]; !ok {
			continue
		}
		if match, _ := regexp.Match(`[\.m]\s+[a-zA-Z0-9\.\-_]+`, nameBytes); !match {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("Could not parse node name:", name))
		}
		if nameBytes[0] == []byte("m")[0] {
			return name
		}
	}
	nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("Could not locate a master"))
	return ""
}

func main() {
	nodeList := flag.String("nodes", "", "Comma-separated list of node names in the cluster")
	flag.Parse()
	nodes := strings.Split(*nodeList, ",")
	flag.Usage = func() {
		flag.PrintDefaults()
	}
	if len(nodes) == 0 || len(*nodeList) == 0 {
		nagiosplugin.Exit(nagiosplugin.UNKNOWN, fmt.Sprintf("No nodes specified"))
	}
	nNodes := len(nodes)
	topologies := make([]string, nNodes)
	masters := make(map[string]bool)
	c := make(chan string, nNodes)
	for _, node := range nodes {
		go getTopology(node, c)
	}
	for i, _ := range nodes {
		topologies[i] = <-c
	}
	masterList := make([]string, 0)
	for _, topology := range topologies {
		topologyMaster := getMaster(topology, nodes)
		// If we haven't seen this master before, add it to the list of 
		// masters to print. 
		if _, ok := masters[topologyMaster]; !ok {
			masterList = append(masterList, topologyMaster)
		}
		masters[topologyMaster] = true
	}
	masterText := strings.Join(masterList, ", ")
	infoText := fmt.Sprintf("%d masters (%s)", len(masters), masterText)
	exitStatus := nagiosplugin.UNKNOWN
	if len(masters) > 1 {
		exitStatus = nagiosplugin.CRITICAL
	}
	if len(masters) == 1 {
		exitStatus = nagiosplugin.OK
	}
	nagiosplugin.Exit(exitStatus, infoText)
}

