//go:build e2e

package create_test

import (
	"errors"
	"io/fs"
	"os"

	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/ocireg"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/github"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/localblob"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/ociartifact"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	ocmmetav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	compdescv2 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/versions/v2"
	ocmocireg "github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ocireg"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/kyma-project/lifecycle-manager/api/shared"
	"github.com/kyma-project/lifecycle-manager/api/v1beta2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var errCreateModuleFailedWithSameVersion = errors.New(
	"version 1.0.0 already exists with different content, please use --module-archive-version-overwrite flag to overwrite it")

var _ = Describe("Test 'create' command", Ordered, func() {
	// _ = os.Setenv("OCI_REPOSITORY_URL", "http://k3d-oci.localhost:5001")
	// _ = os.Setenv("MODULE_TEMPLATE_PATH", "/tmp/module-config-template.yaml")
	// _ = os.Setenv("MODULE_TEMPLATE_VERSION", "1.0.0")

	// ociRegistry := os.Getenv("OCI_REPOSITORY_URL")
	// testRepoURL := os.Getenv("TEST_REPOSITORY_URL")
	// templateOutputPath := os.Getenv("MODULE_TEMPLATE_PATH")
	// moduleTemplateVersion := os.Getenv("MODULE_TEMPLATE_VERSION")

	ociRegistry := "http://k3d-oci.localhost:5001"
	moduleRepoPath := "./testdata/template-operator/"
	moduleTemplateVersion := "1.0.0"
	moduleConfigFilePath := moduleRepoPath + "module-config.yaml"
	templateOutputPath := moduleRepoPath + "template.yaml"
	securityScanConfigFile := moduleRepoPath + "sec-scanners-config.yaml"
	changedSecScanConfigFile := moduleRepoPath + "sec-scanners-config-changed.yaml"

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked without any args", func() {
			cmd = createCmd{}
		})

		It("Then the command should fail", func() {
			err := cmd.execute()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Error: \"--module-config-file\" flag is required"))

			By("And no module template.yaml is generated")
			Expect(filesIn(moduleRepoPath)).Should(Not(ContainElement("template.yaml")))
		})
	})

	Context("Given 'modulectl create' command", func() {
		var cmd createCmd
		It("When invoked with '--module-config-file' and '--path'", func() {
			cmd = createCmd{
				moduleConfigFile: moduleConfigFilePath,
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
				moduleConfigFile: moduleConfigFilePath,
				path:             moduleRepoPath,
				registry:         ociRegistry,
				insecure:         true,
				output:           templateOutputPath,
			}
		})
		It("Then the command should succeed", func() {
			Expect(cmd.execute()).To(Succeed())

			By("And module template file should be generated")
			Expect(filesIn("./testdata/template-operator/")).Should(ContainElement("template.yaml"))
		})
		It("Then module template should contain the expected content", func() {
			template, err := readModuleTemplate(templateOutputPath)
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
			Expect(repo.Object["baseUrl"]).To(Equal(ociRegistry))
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
			Expect(resource.Name).To(Equal("raw-manifest"))
			Expect(resource.Relation).To(Equal(ocmmetav1.LocalRelation))
			Expect(resource.Type).To(Equal("yaml"))
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
				registry:                      ociRegistry,
				moduleConfigFile:              moduleConfigFilePath,
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
				registry:         ociRegistry,
				moduleConfigFile: moduleConfigFilePath,
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
				registry:         ociRegistry,
				moduleConfigFile: moduleConfigFilePath,
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

// Test helper functions

func readModuleTemplate(filepath string) (*v1beta2.ModuleTemplate, error) {
	moduleTemplate := &v1beta2.ModuleTemplate{}
	moduleFile, err := os.ReadFile(filepath)
	if err != nil && len(moduleFile) > 0 {
		return nil, err
	}
	err = yaml.Unmarshal(moduleFile, moduleTemplate)
	if err != nil {
		return nil, err
	}
	return moduleTemplate, err
}

func getDescriptor(template *v1beta2.ModuleTemplate) *v1beta2.Descriptor {
	if template.Spec.Descriptor.Object != nil {
		desc, ok := template.Spec.Descriptor.Object.(*v1beta2.Descriptor)
		if !ok || desc == nil {
			return nil
		}
		return desc
	}
	ocmDesc, err := compdesc.Decode(
		template.Spec.Descriptor.Raw,
		[]compdesc.DecodeOption{compdesc.DisableValidation(true)}...)
	if err != nil {
		return nil
	}
	template.Spec.Descriptor.Object = &v1beta2.Descriptor{ComponentDescriptor: ocmDesc}
	desc, ok := template.Spec.Descriptor.Object.(*v1beta2.Descriptor)
	if !ok {
		return nil
	}

	return desc
}

func flatten(labels ocmmetav1.Labels) map[string]string {
	labelsMap := make(map[string]string)
	for _, l := range labels {
		var value string
		_ = yaml.Unmarshal(l.Value, &value)
		labelsMap[l.Name] = value
	}
	return labelsMap
}

func filesIn(dir string) []string {
	fi, err := os.Stat(dir)
	Expect(err).ToNot(HaveOccurred())
	Expect(fi.IsDir()).To(BeTrueBecause("The provided path should be a directory: %s", dir))

	dirFs := os.DirFS(dir)
	entries, err := fs.ReadDir(dirFs, ".")
	Expect(err).ToNot(HaveOccurred())

	var res []string
	for _, ent := range entries {
		if ent.Type().IsRegular() {
			res = append(res, ent.Name())
		}
	}

	return res
}
