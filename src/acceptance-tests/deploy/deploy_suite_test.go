package deploy_test

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/cloudfoundry-incubator/etcd-release/src/acceptance-tests/testing/helpers"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/pivotal-cf-experimental/bosh-test/bosh"

	directorCli "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	config        helpers.Config
	boshClient    bosh.Client
	boshCliClient directorCli.Director
)

func TestDeploy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "deploy")
}

var _ = BeforeSuite(func() {
	configPath, err := helpers.ConfigPath()
	Expect(err).NotTo(HaveOccurred())

	config, err = helpers.LoadConfig(configPath)
	Expect(err).NotTo(HaveOccurred())

	boshClient = bosh.NewClient(bosh.Config{
		URL:              config.BOSH.Target,
		Username:         config.BOSH.Username,
		Password:         config.BOSH.Password,
		AllowInsecureSSL: true,
	})
	boshCliClient = getDirector()
})

func getDirector() directorCli.Director {
	directorUrl, err := url.Parse(config.BOSH.Target)
	Expect(err).ToNot(HaveOccurred())

	port, err := strconv.Atoi(directorUrl.Port())
	Expect(err).ToNot(HaveOccurred())

	factoryConfig := directorCli.FactoryConfig{
		Host:         directorUrl.Hostname(),
		Port:         port,
		Client:       config.BOSH.Username,
		ClientSecret: config.BOSH.Password,
		CACert:       config.BOSH.DirectorCACert,
	}
	logger := boshlog.NewLogger(boshlog.LevelError)
	directorFactory := directorCli.NewFactory(logger)
	boshDirector, err := directorFactory.New(factoryConfig, directorCli.NewNoopTaskReporter(), directorCli.NewNoopFileReporter())
	Expect(err).NotTo(HaveOccurred())

	return boshDirector
}
