# Kubo-etcd release

This is a fork of the [cloudfoundry-incubator/etcd-release](https://github.com/cloudfoundry-incubator/etcd-release). This is actively used in the KuBo project (https://github.com/cloudfoundry-incubator/kubo-deployment).
It uses ETCD v3 underneath.

For the original documentation on how to use etcd release, navigate to the [README here](https://github.com/cloudfoundry-incubator/etcd-release/blob/master/README.md)

# Configuring kubo-etcd with TLS

When configuring the original etcd-release with TLS, consul was a required dependency as communication to the etcd nodes were done by consul DNS addresses.

For Kubo, DNS addresses are managed by [BOSH DNS](https://bosh.io/docs/dns.html#dns-release) meaning the certificates generated for peer and client/server TLS communication must be generated accordingly. Given the bosh property `etcd.advertise_urls_dns_suffix` is set  (i.e. etcd.kubo), the certificates must have the following configuration:

* Common Name (CN): etcd.kubo
* Subject Alternative Names (SANs): etcd.kubo, *.etcd.kubo
