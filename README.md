# amilib
<H1>Asterisk AMI library for amiproxy</H1>

The premises is you will have a cluster of asterisk servers , This library will allow to queue commands in any of the asterisks of the cluster or in an AWS loadbalance etc, etc. The goal is reduce the load in a system by accessing the asterisk using a proxy which is part of the project to queue the commands  instead the AMI directly. 
There are libraries similar to this one using direct access to the port or using a different protocols, mostly ARI  that require more call handling or event handling. The idea of this library is push and forget, trying to archieve high volumen of processing calls or events without needing a lot of ports to get this.
