check_elasticsearch_splitbrain
==============================

Nagios plugin to check for splitbrains in Elasticsearch clusters. 

How it works
============

The check is pretty simple - it spawns a goroutine for each node in the
cluster, uses NRPE to retrieve the node's view of its own cluster
topology (including who is master), and then collates the information.
It returns OK if all the nodes have the same master, CRITICAL if there's
more than one master detected, and UNKNOWN otherwise. 

Requirements
============

 - esadmin (https://github.com/anchor/elasticsearchadmin)

This plugin uses a Go package called go-nagios, formerly at 
https://github.com/laziac/go-nagios; it's been taken down, and I'm
waiting to hear back from the author regarding code re-use. If this
doesn't happen I'll rewrite and update the check.

Configuration
=============

This check has two parts - check_elasticsearch_splitbrain and 
check_elasticsearch_topology. The latter is an info-only check (always
returns OK) and is invoked by check_elasticsearch_splitbrain on each
host passed via the -nodes argument.

