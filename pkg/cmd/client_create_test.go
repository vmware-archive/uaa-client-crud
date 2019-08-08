package cmd_test

import (
	"bufio"
	"bytes"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Client create", func() {
	var b bytes.Buffer
	var out *bufio.Writer
	var c *cobra.Command

	BeforeEach(func() {
		b = bytes.Buffer{}
		out = bufio.NewWriter(&b)
		c = cmd.NewCreateClientCmd(out)
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
