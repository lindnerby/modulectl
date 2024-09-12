//go:build e2e

package create_test

import (
	"github.com/kyma-project/lifecycle-manager/api/shared"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/ocireg"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/github"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/localblob"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/ociartifact"
	ocmmetav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	compdescv2 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/versions/v2"
	ocmocireg "github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ocireg"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test 'create' command", Ordered, func() {
	//var initialDir string
	//var workDir string

	// TODO remove debugging env
	// err := os.Setenv("OCI_REPOSITORY_URL", "http://k3d-oci.localhost:5001")
	// err = os.Setenv("TEST_REPOSITORY_URL", "https://github.com/lindnerby/template-operator.git")
	// err = os.Setenv("MODULE_TEMPLATE_PATH", "/tmp/module-config-template.yaml")
	// err = os.Setenv("MODULE_TEMPLATE_VERSION", "1.0.0")

	// TODO decide what should be configurable
	//ociRepoURL := os.Getenv("OCI_REPOSITORY_URL")
	//testRepoURL := os.Getenv("TEST_REPOSITORY_URL")
	//templatePath := os.Getenv("MODULE_TEMPLATE_PATH")
	//moduleTemplateVersion := os.Getenv("MODULE_TEMPLATE_VERSION")

	// TODO adapt path for checked out repo on pipeline
	moduleRepoPath := "./testdata/template-operator/"
	configFilePath := "./testdata/template-operator/module-config.yaml"
	localRegistry := "http://k3d-oci.localhost:5001"
	templateOutput := "./testdata/template-operator/template.yaml"
	moduleTemplateVersion := "1.0.0"
	securityScanConfigFile := "./testdata/template-operator/sec-scanners-config.yaml"
	changedSecScanConfigFile := "./testdata/template-operator/sec-scanners-config-changed.yaml"
	//changedSecScannerConfigFile := "../../../../template-operator/sec-scanners-config-changed.yaml"

	// TODO see if we can use elegant dir reference to template-operator
	//setup := func() {
	//	var err error
	//	initialDir, err = os.Getwd()
	//	Expect(err).ToNot(HaveOccurred())
	//	workDir = resolveWorkingDirectory()
	//	err = os.Chdir(workDir)
	//	Expect(err).ToNot(HaveOccurred())
	//}
	//
	//teardown := func() {
	//	err := os.Chdir(initialDir)
	//	Expect(err).ToNot(HaveOccurred())
	//	cleanupWorkingDirectory(workDir)
	//	workDir = ""
	//	initialDir = ""
	//}

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked without any args", func() {
			//print current dir
			currDir, _ := os.Getwd()
			println(currDir)
			cmd = createCmd{}
		})

		It("Then the command should fail", func() {
			cmd = createCmd{}
			err := cmd.execute()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Error: \"--module-config-file\" flag is required"))

			// TODO
			//By("And no files are generated")
			//Expect(filesIn(workDir)).Should(BeEmpty())
		})
	})

	Context("Given `modulectl create` command", func() {
		var cmd createCmd
		It("When invoked with --module-config-file and --path", func() {
			cmd = createCmd{
				moduleConfigFile: configFilePath,
				path:             moduleRepoPath,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())

			By("And no module template file should be generated")
			Expect(filesIn("./testdata/template-operator/")).Should(Not(ContainElement("template.yaml")))
		})
	})

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked with existing '--registry' and '--insecure' flag", func() {
			cmd = createCmd{
				moduleConfigFile: configFilePath,
				path:             moduleRepoPath,
				registry:         localRegistry,
				insecure:         true,
				output:           templateOutput,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())

			By("And module template file should be generated")
			Expect(filesIn("./testdata/template-operator/")).Should(ContainElement("template.yaml"))
		})
		It("Then module template should contain the expected content", func() {
			template, err := readModuleTemplate(templateOutput)
			Expect(err).ToNot(HaveOccurred())
			descriptor := getDescriptor(template)
			Expect(descriptor).ToNot(BeNil())
			Expect(descriptor.SchemaVersion()).To(Equal(compdescv2.SchemaVersion))

			By("And annotations should be correct")
			annotations := template.Annotations
			Expect(annotations[shared.ModuleVersionAnnotation]).To(Equal(moduleTemplateVersion))
			Expect(annotations[shared.IsClusterScopedAnnotation]).To(Equal("false"))

			By("And descriptor.component.repositoryContexts should be correct")
			Expect(descriptor.RepositoryContexts).To(HaveLen(1))
			repo := descriptor.GetEffectiveRepositoryContext()
			Expect(repo.Object["baseUrl"]).To(Equal(localRegistry))
			Expect(repo.Object["componentNameMapping"]).To(Equal(string(ocmocireg.OCIRegistryURLPathMapping)))
			Expect(repo.Object["type"]).To(Equal(ocireg.Type))

			By("And descriptor.component.resources should be correct")
			Expect(descriptor.Resources).To(HaveLen(2))
			resource := descriptor.Resources[0]
			Expect(resource.Name).To(Equal("template-operator"))
			Expect(resource.Relation).To(Equal(ocmmetav1.ExternalRelation))
			Expect(resource.Type).To(Equal("ociImage"))
			Expect(resource.Version).To(Equal(moduleTemplateVersion))

			resource = descriptor.Resources[1]
			Expect(resource.Name).To(Equal(rawManifestLayerName))
			Expect(resource.Relation).To(Equal(ocmmetav1.LocalRelation))
			Expect(resource.Type).To(Equal(typeYaml))
			Expect(resource.Version).To(Equal(moduleTemplateVersion))

			By("And descriptor.component.resources[0].access should be correct")
			resourceAccessSpec0, err := ocm.DefaultContext().AccessSpecForSpec(descriptor.Resources[0].Access)
			Expect(err).ToNot(HaveOccurred())
			ociArtifactAccessSpec, ok := resourceAccessSpec0.(*ociartifact.AccessSpec)
			Expect(ok).To(BeTrue())
			Expect(ociArtifactAccessSpec.GetType()).To(Equal(ociartifact.Type))
			Expect(ociArtifactAccessSpec.ImageReference).To(Equal("europe-docker.pkg.dev/kyma-project/prod/template-operator:1.0.0"))

			By("And descriptor.component.resources[1].access should be correct")
			resourceAccessSpec1, err := ocm.DefaultContext().AccessSpecForSpec(descriptor.Resources[1].Access)
			Expect(err).ToNot(HaveOccurred())
			localBlobAccessSpec, ok := resourceAccessSpec1.(*localblob.AccessSpec)
			Expect(ok).To(BeTrue())
			Expect(localBlobAccessSpec.GetType()).To(Equal(localblob.Type))
			Expect(localBlobAccessSpec.LocalReference).To(ContainSubstring("sha256:"))

			By("And descriptor.component.sources should be correct")
			Expect(len(descriptor.Sources)).To(Equal(1))
			source := descriptor.Sources[0]
			sourceAccessSpec, err := ocm.DefaultContext().AccessSpecForSpec(source.Access)
			Expect(err).ToNot(HaveOccurred())
			githubAccessSpec, ok := sourceAccessSpec.(*github.AccessSpec)
			Expect(ok).To(BeTrue())
			Expect(github.Type).To(Equal(githubAccessSpec.Type))
			Expect(githubAccessSpec.RepoURL).To(ContainSubstring("template-operator.git"))

			By("And spec.mandatory should be false")
			Expect(template.Spec.Mandatory).To(BeFalse())

			By("And security scan labels should be correct")
			secScanLabels := flatten(descriptor.Sources[0].Labels)
			Expect(secScanLabels).To(HaveKeyWithValue("git.kyma-project.io/ref", "refs/heads/main"))
			Expect(secScanLabels).To(HaveKeyWithValue("scan.security.kyma-project.io/rc-tag", "1.0.0"))
			Expect(secScanLabels).To(HaveKeyWithValue("scan.security.kyma-project.io/language", "golang-mod"))
			Expect(secScanLabels).To(HaveKeyWithValue("scan.security.kyma-project.io/dev-branch", "main"))
			Expect(secScanLabels).To(HaveKeyWithValue("scan.security.kyma-project.io/subprojects", "false"))
			Expect(secScanLabels).To(HaveKeyWithValue("scan.security.kyma-project.io/exclude", "**/test/**,**/*_test.go"))
		})
	})

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked with '--module-archive-version-overwrite' flag", func() {
			cmd = createCmd{
				path:                          moduleRepoPath,
				registry:                      localRegistry,
				moduleConfigFile:              configFilePath,
				version:                       moduleTemplateVersion,
				moduleArchiveVersionOverwrite: true,
				secScanConfig:                 securityScanConfigFile,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())
		})
	})

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked with same version module and same content without '--module-archive-version-overwrite' flag", func() {
			cmd = createCmd{
				path:             moduleRepoPath,
				registry:         localRegistry,
				moduleConfigFile: configFilePath,
				version:          moduleTemplateVersion,
				secScanConfig:    securityScanConfigFile,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())
		})
	})

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked with same version module, but different content without '--module-archive-version-overwrite' flag", func() {
			cmd = createCmd{
				path:             moduleRepoPath,
				registry:         localRegistry,
				moduleConfigFile: configFilePath,
				version:          moduleTemplateVersion,
				secScanConfig:    changedSecScanConfigFile,
			}
		})
		It("Then the command should fail with same version exists message", func() {
			err := cmd.execute()
			Expect(err).To(Equal(errCreateModuleFailedWithSameVersion))
		})
	})
})
