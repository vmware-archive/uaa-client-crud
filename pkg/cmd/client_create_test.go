package cmd_test

import (
	"bufio"
	"bytes"
	"os"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/pkg/errors"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces"
	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces/interfacesfakes"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Client create", func() {
	var b bytes.Buffer
	var out *bufio.Writer
	var c *cobra.Command
	var uaaEndpoint string
	var credhubEndpoint string

	var fakeUaaClient *interfacesfakes.FakeUaaAPI
	var fakeCredHubClient *interfacesfakes.FakeCredHubAPI

	uaaApiFactory := func(target string, zoneID string, adminClientIdentity string, adminClientPwd string) interfaces.UaaAPI {
		uaaEndpoint = target
		return fakeUaaClient
	}

	credHubFactory := func(target string, skipTLS bool, clientID string, clientPwd string, uaaEndpoint string) (interfaces.CredHubAPI, error) {
		credhubEndpoint = target
		return fakeCredHubClient, nil
	}

	BeforeEach(func() {
		b = bytes.Buffer{}
		out = bufio.NewWriter(&b)

		fakeUaaClient = &interfacesfakes.FakeUaaAPI{}
		fakeCredHubClient = &interfacesfakes.FakeCredHubAPI{}

		c = cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
		c.Flags().Set("uaa-endpoint", "bob")
		c.Flags().Set("admin-identity", "monkey123")
		c.Flags().Set("admin-pwd", "bob")
		c.Flags().Set("target-client-identity", "monkey123")
		c.Flags().Set("target-client-password", "p@ssw0rD")
		c.Flags().Set("authorities", "auth1, auth2")
	})

	It("happy path client does not already exist", func() {

		fakeUaaClient.GetClientReturns(nil, errors.New("client does not exist"))
		fakeUaaClient.CreateClientReturns(&uaa.Client{ClientID: "monkey123"}, nil)
		fakeCredHubClient.GetPermissionByPathActorReturns(nil, errors.New("permission not found"))
		fakeCredHubClient.AddPermissionReturns(nil, nil)

		c.Flags().Set("credhub-identity", "bob")
		c.Flags().Set("credhub-password", "monkey123")
		c.Flags().Set("credhub-endpoint", "bob")
		c.Flags().Set("credential-path", "monkey123")
		c.PreRun(c, []string{})
		err := c.RunE(c, []string{})
		out.Flush()

		Expect(err).To(BeNil())
		Expect(fakeUaaClient.CreateClientCallCount()).To(Equal(1))
		Expect(fakeCredHubClient.AddPermissionCallCount()).To(Equal(1))
	})

	It("setting env vars for credhub are passed to clientCreate", func() {

		os.Setenv("CREDHUB_CLIENT_ID", "notbob")
		os.Setenv("CREDHUB_CLIENT_PASSWORD", "monkey123")
		os.Setenv("CREDHUB_URL", "https://credhub.endpoint")
		os.Setenv("CREDHUB_CRED_PATH", "/path")

		cc := cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)

		Expect(cc.Flag("credhub-identity").Value.String()).To(Equal("notbob"))
		Expect(cc.Flag("credhub-password").Value.String()).To(Equal("monkey123"))
		Expect(cc.Flag("credhub-endpoint").Value.String()).To(Equal("https://credhub.endpoint"))
		Expect(cc.Flag("credential-path").Value.String()).To(Equal("/path"))

	})

	It("credhub and uaa endpoints gets appended https://", func() {

		fakeUaaClient.GetClientReturns(nil, errors.New("client does not exist"))
		fakeUaaClient.CreateClientReturns(&uaa.Client{ClientID: "monkey123"}, nil)
		fakeCredHubClient.GetPermissionByPathActorReturns(nil, errors.New("permission not found"))
		fakeCredHubClient.AddPermissionReturns(nil, nil)

		os.Setenv("CREDHUB_CLIENT_ID", "notbob")
		os.Setenv("CREDHUB_CLIENT_PASSWORD", "monkey123")
		os.Setenv("CREDHUB_URL", "credhub.endpoint")
		os.Setenv("CREDHUB_CRED_PATH", "/path")

		cc := cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
		cc.Flags().Set("uaa-endpoint", "bob")

		cc.PreRun(cc, []string{})
		err := cc.RunE(cc, []string{})

		Expect(err).To(BeNil())
		Expect(credhubEndpoint).To(Equal("https://credhub.endpoint"))
		Expect(uaaEndpoint).To(Equal("https://bob"))

	})
})
