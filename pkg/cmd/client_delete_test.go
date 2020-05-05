package cmd_test

import (
	"bufio"
	"bytes"
	"os"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"

	"code.cloudfoundry.org/credhub-cli/credhub/permissions"
	"github.com/cloudfoundry-community/go-uaa"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces"
	"github.com/cf-platform-eng/uaa-client-crud/pkg/interfaces/interfacesfakes"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Client delete", func() {
	var fakeUaaClient *interfacesfakes.FakeUaaAPI
	var fakeCredHubClient *interfacesfakes.FakeCredHubAPI

	var b bytes.Buffer
	var out *bufio.Writer
	var c *cobra.Command
	var uaaEndpoint string
	var credhubEndpoint string

	uaaApiFactory := func(target string, zoneID string, adminClientIdentity string, adminClientSecret string) interfaces.UaaAPI {
		uaaEndpoint = target
		return fakeUaaClient
	}

	credHubFactory := func(target string, skipTLS bool, clientID string, clientSecret string, uaaEndpoint string) (interfaces.CredHubAPI, error) {
		credhubEndpoint = target
		return fakeCredHubClient, nil
	}

	BeforeEach(func() {
		b = bytes.Buffer{}
		out = bufio.NewWriter(&b)

		fakeUaaClient = &interfacesfakes.FakeUaaAPI{}
		fakeCredHubClient = &interfacesfakes.FakeCredHubAPI{}

		c = cmd.NewDeleteClientCmd(uaaApiFactory, credHubFactory, out)
		c.Flags().Set("uaa-endpoint", "bob")
		c.Flags().Set("admin-identity", "monkey123")
		c.Flags().Set("admin-pwd", "bob")
		c.Flags().Set("target-client-identity", "monkey123")
	})

	It("happy path client does exist and gets deleted", func() {

		fakeUaaClient.GetClientReturns(&uaa.Client{}, nil)
		fakeUaaClient.DeleteClientReturns(&uaa.Client{}, nil)
		fakeCredHubClient.GetPermissionByPathActorReturns(&permissions.Permission{UUID: "123"}, nil)
		fakeCredHubClient.DeletePermissionReturns(&permissions.Permission{}, nil)

		c.Flags().Set("credhub-identity", "bob")
		c.Flags().Set("credhub-secret", "monkey123")
		c.Flags().Set("credhub-endpoint", "bob")
		c.Flags().Set("credential-path", "monkey123")

		c.PreRun(c, []string{})
		err := c.RunE(c, []string{})
		out.Flush()

		Expect(err).To(BeNil())
		Expect(fakeUaaClient.DeleteClientCallCount()).To(Equal(1))
		Expect(fakeUaaClient.DeleteClientArgsForCall(0)).To(Equal("monkey123"))
		Expect(fakeCredHubClient.DeletePermissionCallCount()).To(Equal(1))
		Expect(fakeCredHubClient.DeletePermissionArgsForCall(0)).To(Equal("123"))

	})

	It("happy path client does exist and gets deleted with credentials", func() {

		fakeUaaClient.GetClientReturns(&uaa.Client{}, nil)
		fakeUaaClient.DeleteClientReturns(&uaa.Client{}, nil)
		fakeCredHubClient.GetPermissionByPathActorReturns(&permissions.Permission{UUID: "123"}, nil)
		fakeCredHubClient.DeletePermissionReturns(&permissions.Permission{}, nil)
		fakeCredHubClient.FindByPathReturns(credentials.FindResults{
			Credentials: []credentials.Base{{
				Name:             "/c/ksm/blah",
				VersionCreatedAt: "",
			},
			},
		}, nil)

		c.Flags().Set("credhub-identity", "bob")
		c.Flags().Set("credhub-secret", "monkey123")
		c.Flags().Set("credhub-endpoint", "bob")
		c.Flags().Set("credential-path", "/c/ksm/*")
		c.Flags().Set("delete-credhub-path", "true")
		c.Flags().Set("target-client-secret", "the secret")

		c.PreRun(c, []string{})
		err := c.RunE(c, []string{})
		out.Flush()

		Expect(err).To(BeNil())
		Expect(fakeUaaClient.DeleteClientCallCount()).To(Equal(1))
		Expect(fakeUaaClient.DeleteClientArgsForCall(0)).To(Equal("monkey123"))
		Expect(fakeCredHubClient.DeletePermissionCallCount()).To(Equal(1))
		Expect(fakeCredHubClient.DeletePermissionArgsForCall(0)).To(Equal("123"))
		Expect(fakeCredHubClient.FindByPathCallCount()).To(Equal(1))
		Expect(fakeCredHubClient.DeleteCredentialCallCount()).To(Equal(1))
		Expect(fakeCredHubClient.FindByPathArgsForCall(0)).To(Equal("/c/ksm/"))
	})

	It("setting env vars for credhub are passed to clientDelete", func() {

		os.Setenv("CREDHUB_CLIENT_ID", "notbob")
		os.Setenv("CREDHUB_CLIENT_SECRET", "monkey123")
		os.Setenv("CREDHUB_URL", "https://credhub.endpoint")
		os.Setenv("CREDHUB_CRED_PATH", "/path")
		os.Setenv("DELETE_CREDHUB_PATH", "true")

		cc := cmd.NewDeleteClientCmd(uaaApiFactory, credHubFactory, out)

		Expect(cc.Flag("credhub-identity").Value.String()).To(Equal("notbob"))
		Expect(cc.Flag("credhub-secret").Value.String()).To(Equal("monkey123"))
		Expect(cc.Flag("credhub-endpoint").Value.String()).To(Equal("https://credhub.endpoint"))
		Expect(cc.Flag("credential-path").Value.String()).To(Equal("/path"))
		Expect(cc.Flag("delete-credhub-path").Value.String()).To(Equal("true"))
	})

	It("credhub and uaa endpoints gets appended https://", func() {

		fakeUaaClient.GetClientReturns(&uaa.Client{}, nil)
		fakeUaaClient.DeleteClientReturns(&uaa.Client{}, nil)
		fakeCredHubClient.GetPermissionByPathActorReturns(&permissions.Permission{UUID: "123"}, nil)
		fakeCredHubClient.DeletePermissionReturns(&permissions.Permission{}, nil)

		os.Setenv("CREDHUB_CLIENT_ID", "notbob")
		os.Setenv("CREDHUB_CLIENT_SECRET", "monkey123")
		os.Setenv("CREDHUB_URL", "credhub.endpoint")
		os.Setenv("CREDHUB_CRED_PATH", "/path")

		cc := cmd.NewDeleteClientCmd(uaaApiFactory, credHubFactory, out)
		cc.Flags().Set("uaa-endpoint", "bob")

		cc.PreRun(cc, []string{})
		err := cc.RunE(cc, []string{})

		Expect(err).To(BeNil())
		Expect(credhubEndpoint).To(Equal("https://credhub.endpoint"))
		Expect(uaaEndpoint).To(Equal("https://bob"))

	})
})
