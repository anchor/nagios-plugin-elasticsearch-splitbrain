check_elasticsearch_splitbrain
==============================

Nagios plugin to check for splitbrains in Elasticsearch clusters. 

Requirements
============

 - esadmin (https://github.com/anchor/elasticsearchadmin)

This plugin uses a Go package called go-nagios, formerly at 
https://github.com/laziac/go-nagios; it's been taken down, and I'm
waiting to hear back from the author regarding code re-use. If this
doesn't happen I'll rewrite and update the check.
