package utils

import (
	"encoding/json"
	"os"

	types "github.com/alicefr/kubevirt-hook-update-xml/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vmSchema "kubevirt.io/client-go/api/v1"
)

var _ = Describe("KubevirtHook", func() {
	BeforeEach(func() {
		path, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		types.Dir = path
	})

	It("can be loaded from JSON", func() {
		vmi := vmSchema.VirtualMachineInstance{
			ObjectMeta: metav1.ObjectMeta{
				UID: "5555938c-bf2d-48ef-a79c-7f52bba79f40",
			},
		}
		bytes, err := json.Marshal(vmi)
		Expect(err).NotTo(HaveOccurred())
		_, err = MergeKubeVirtXMLWithProvidedXML("vmi-fedora.xml", bytes)
		Expect(err).NotTo(HaveOccurred())
	})
})
