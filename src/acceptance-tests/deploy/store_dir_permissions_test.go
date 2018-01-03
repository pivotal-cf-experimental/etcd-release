package deploy_test

import (
	"fmt"

	etcdclient "github.com/cloudfoundry-incubator/etcd-release/src/acceptance-tests/testing/etcd"
	"github.com/cloudfoundry-incubator/etcd-release/src/acceptance-tests/testing/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/ops"
)

var _ = Describe("StoreDirPermissions", func() {
	var (
		manifest       string
		manifestName   string
		deploymentName string
		etcdClient     etcdclient.Client
	)

	BeforeEach(func() {
		By("deploying etcd cluster", func() {
			deploymentName = "store-dir-permissions-test"

			var err error
			manifest, err = helpers.NewEtcdManifestWithInstanceCount(deploymentName, 1, false, boshClient)
			Expect(err).NotTo(HaveOccurred())

			manifestName, err = ops.ManifestName(manifest)
			Expect(err).NotTo(HaveOccurred())

			manifest, err = ops.ApplyOp(manifest, ops.Op{
				Type: "replace",
				Path: "/instance_groups/name=etcd/jobs/-",
				Value: map[string]string{
					"name":    "iptables_agent",
					"release": "kubo-etcd",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = boshClient.Deploy([]byte(manifest))
			Expect(err).NotTo(HaveOccurred())

			testConsumerIPs, err := helpers.GetVMIPs(boshClient, manifestName, "testconsumer")
			Expect(err).NotTo(HaveOccurred())

			etcdClient = etcdclient.NewClient(fmt.Sprintf("http://%s:6769", testConsumerIPs[0]))
		})
	})

	AfterEach(func() {
		if !CurrentGinkgoTestDescription().Failed {
			err := boshClient.DeleteDeployment(manifestName)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should make store dir non accessible by other users than vcap", func() {
		host, username, privateKey, err := helpers.GetSSHCreds(manifestName, boshCliClient)
		Expect(err).ToNot(HaveOccurred())

		output, err := helpers.RunSSHCommand(host, 22, username, privateKey, "stat -c 'Owner: %U, Access Rights: %a' /var/vcap/store/etcd")

		Expect(err).ToNot(HaveOccurred())
		Expect(output).To(Equal("Owner: vcap, Access Rights: 700\n"))
	})
})
