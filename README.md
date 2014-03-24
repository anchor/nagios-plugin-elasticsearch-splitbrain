nagios-plugin-elasticsearch-splitbrain
======================================

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

 - [esadmin][0]
 - [nagiosplugin][1]

Configuration
=============

This check has two parts - check_elasticsearch_splitbrain and 
check_elasticsearch_topology. The latter is an info-only check (always
returns OK) and is invoked by check_elasticsearch_splitbrain on each
host passed via the -nodes argument.

[0]: https://github.com/anchor/elasticsearchadmin
[1]: https://github.com/fractalcat/nagiosplugin
