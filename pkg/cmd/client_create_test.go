package cmd_test

import (
	"bufio"
	"bytes"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces"
	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces/interfacesfakes"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var fakeUaaClient *interfacesfakes.FakeUaaAPI
var fakeCredHubClient *interfacesfakes.FakeCredHubAPI

var _ = Describe("Client create", func() {
	var b bytes.Buffer
	var out *bufio.Writer
	var c *cobra.Command

	BeforeEach(func() {
		b = bytes.Buffer{}
		out = bufio.NewWriter(&b)

		fakeUaaClient = &interfacesfakes.FakeUaaAPI{}
		fakeCredHubClient = &interfacesfakes.FakeCredHubAPI{}

		uaaApiFactory := func(target string, zoneID string, adminClientIdentity string, adminClientPwd string) interfaces.UaaAPI {
			return fakeUaaClient
		}

		credHubFactory := func(target string, skipTLS bool, clientID string, clientPwd string, uaaEndpoint string) interfaces.CredHubAPI {
			return fakeCredHubClient
		}

		c = cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
		c.Flags().Set("uaa-endpoint", "bob")
		c.Flags().Set("admin-identity", "monkey123")
		c.Flags().Set("admin-pwd", "bob")
		c.Flags().Set("target-client-identity", "monkey123")
		c.Flags().Set("target-client-password", "p@ssw0rD")
		c.Flags().Set("authorities", "auth1, auth2")
	})

	XIt("happy path", func() {
		c.Flags().Set("credhub-identity", "bob")
		c.Flags().Set("credhub-password", "monkey123")
		c.Flags().Set("credhub-endpoint", "bob")
		c.Flags().Set("credential-path", "monkey123")
		c.PreRun(c, []string{})
		err := c.RunE(c, []string{})
		out.Flush()

		Expect(err).To(BeNil())
	})

	XIt("setting env vars are passed to clientCreate", func() {
		c.Flags().Set("credhub-identity", "bob")
		c.Flags().Set("credhub-password", "monkey123")
		c.Flags().Set("credhub-endpoint", "bob")
		c.Flags().Set("credential-path", "monkey123")
		c.PreRun(c, []string{})
		err := c.RunE(c, []string{})
		out.Flush()
		Expect(err).To(BeNil())
	})
})
