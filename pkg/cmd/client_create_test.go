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

	uaaApiFactory := func(target string, zoneID string, adminClientIdentity string, adminClientSecret string) (interfaces.UaaAPI, error) {
		uaaEndpoint = target
		return fakeUaaClient, nil
	}

	credHubFactory := func(target string, skipTLS bool, clientID string, clientSecret string, uaaEndpoint string) (interfaces.CredHubAPI, error) {
		credhubEndpoint = target
		return fakeCredHubClient, nil
	}

	Context("No environment variables set", func() {
		BeforeEach(func() {
			b = bytes.Buffer{}
			out = bufio.NewWriter(&b)

			fakeUaaClient = &interfacesfakes.FakeUaaAPI{}
			fakeCredHubClient = &interfacesfakes.FakeCredHubAPI{}

			c = cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
		})

		AfterEach(func() {
			os.Unsetenv("CREDHUB_CLIENT_ID")
			os.Unsetenv("CREDHUB_CLIENT_SECRET")
			os.Unsetenv("CREDHUB_URL")
			os.Unsetenv("CREDHUB_CRED_PATH")
			uaaEndpoint = ""
			credhubEndpoint = ""
		})

		It("client does not already exist creates client", func() {

			fakeUaaClient.GetClientReturns(nil, errors.New("client does not exist"))
			fakeUaaClient.CreateClientReturns(&uaa.Client{ClientID: "monkey123"}, nil)
			fakeCredHubClient.GetPermissionByPathActorReturns(nil, errors.New("permission not found"))
			fakeCredHubClient.AddPermissionReturns(nil, nil)

			c.Flags().Set("uaa-endpoint", "bob")
			c.Flags().Set("admin-identity", "monkey123")
			c.Flags().Set("admin-secret", "bob")
			c.Flags().Set("target-client-identity", "monkey123")
			c.Flags().Set("target-client-secret", "p@ssw0rD")
			c.Flags().Set("authorities", "auth1, auth2")
			c.Flags().Set("credhub-identity", "bob")
			c.Flags().Set("credhub-secret", "monkey123")
			c.Flags().Set("credhub-endpoint", "bob")
			c.Flags().Set("credential-path", "monkey123")
			c.PreRun(c, []string{})
			err := c.RunE(c, []string{})
			out.Flush()

			Expect(err).To(BeNil())
			Expect(fakeUaaClient.CreateClientCallCount()).To(Equal(1))
			Expect(fakeCredHubClient.AddPermissionCallCount()).To(Equal(1))
		})

		It("Does not append https when credhub is empty", func() {
			fakeUaaClient.GetClientReturns(nil, errors.New("client does not exist"))
			fakeUaaClient.CreateClientReturns(&uaa.Client{ClientID: "monkey123"}, nil)

			c := cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
			c.Flags().Set("uaa-endpoint", "bob")
			c.PreRun(c, []string{})
			err := c.RunE(c, []string{})
			Expect(err).To(BeNil())
			Expect(fakeCredHubClient.GetPermissionByPathActorCallCount()).To(BeZero())
			Expect(uaaEndpoint).To(Equal("https://bob"))
			Expect(credhubEndpoint).To(Equal(""))

		})
	})

	Context("Environment variables set", func() {
		BeforeEach(func() {
			b = bytes.Buffer{}
			out = bufio.NewWriter(&b)

			fakeUaaClient = &interfacesfakes.FakeUaaAPI{}
			fakeCredHubClient = &interfacesfakes.FakeCredHubAPI{}

		})

		AfterEach(func() {
			os.Unsetenv("CREDHUB_CLIENT_ID")
			os.Unsetenv("CREDHUB_CLIENT_SECRET")
			os.Unsetenv("CREDHUB_URL")
			os.Unsetenv("CREDHUB_CRED_PATH")
			uaaEndpoint = ""
			credhubEndpoint = ""
		})

		It("setting env vars for credhub are passed to clientCreate", func() {

			os.Setenv("CREDHUB_CLIENT_ID", "notbob")
			os.Setenv("CREDHUB_CLIENT_SECRET", "monkey123")
			os.Setenv("CREDHUB_URL", "https://credhub.endpoint")
			os.Setenv("CREDHUB_CRED_PATH", "/path")

			c = cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
			Expect(c.Flag("credhub-identity").Value.String()).To(Equal("notbob"))
			Expect(c.Flag("credhub-secret").Value.String()).To(Equal("monkey123"))
			Expect(c.Flag("credhub-endpoint").Value.String()).To(Equal("https://credhub.endpoint"))
			Expect(c.Flag("credential-path").Value.String()).To(Equal("/path"))

		})

		It("credhub and uaa endpoints gets appended https://", func() {

			fakeUaaClient.GetClientReturns(nil, errors.New("client does not exist"))
			fakeUaaClient.CreateClientReturns(&uaa.Client{ClientID: "monkey123"}, nil)
			fakeCredHubClient.GetPermissionByPathActorReturns(nil, errors.New("permission not found"))
			fakeCredHubClient.AddPermissionReturns(nil, nil)

			os.Setenv("CREDHUB_CLIENT_ID", "notbob2")
			os.Setenv("CREDHUB_CLIENT_SECRET", "monkey1234")
			os.Setenv("CREDHUB_URL", "credhub.endpoint")
			os.Setenv("CREDHUB_CRED_PATH", "/path")

			c = cmd.NewCreateClientCmd(uaaApiFactory, credHubFactory, out)
			c.Flags().Set("uaa-endpoint", "bob")

			c.PreRun(c, []string{})
			err := c.RunE(c, []string{})

			Expect(err).To(BeNil())
			Expect(credhubEndpoint).To(Equal("https://credhub.endpoint"))
			Expect(uaaEndpoint).To(Equal("https://bob"))

		})
	})

})
