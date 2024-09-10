//go:build e2e

package create_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test 'create' Command", Ordered, func() {
	var initialDir string
	var workDir string

	// TODO adapt path for checked out repo on pipeline
	moduleRepoPath := "../../../../template-operator"
	configFilePath := "../../../../template-operator/module-config.yaml"
	//secScannerConfigFile := "../../../../template-operator/sec-scanners-config.yaml"
	//
	//changedSecScannerConfigFile := "../../../../template-operator/sec-scanners-config-changed.yaml"

	// TODO see if we can use elegant dir reference to template-operator
	setup := func() {
		var err error
		initialDir, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())
		workDir = resolveWorkingDirectory()
		err = os.Chdir(workDir)
		Expect(err).ToNot(HaveOccurred())
	}

	teardown := func() {
		err := os.Chdir(initialDir)
		Expect(err).ToNot(HaveOccurred())
		cleanupWorkingDirectory(workDir)
		workDir = ""
		initialDir = ""
	}

	Context("Given `modulectl create` command", func() {
		var cmd createCmd
		It("When invoked without any args", func() {
			//print current dir
			currDir, _ := os.Getwd()
			By("Current dir: " + currDir)
			cmd = createCmd{}
		})

		It("Then the command should fail", func() {
			cmd = createCmd{}
			err := cmd.execute()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Error: \"--module-config-file\" flag is required"))

			//By("And no files are generated")
			//Expect(filesIn(workDir)).Should(BeEmpty())
		})
	})

	Context("Given `modulectl create` command", func() {
		BeforeAll(func() { setup() })
		AfterAll(func() { teardown() })

		var cmd createCmd
		It("When invoked with --module-config-file and --path", func() {
			cmd = createCmd{
				moduleConfigFile: configFilePath,
				path:             moduleRepoPath,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())

			By("And module template file should be generated")
			Expect(filesIn(workDir)).Should(HaveLen(1))
			Expect(filesIn(workDir)).Should(ContainElement("template.yaml"))
		})
	})
})
