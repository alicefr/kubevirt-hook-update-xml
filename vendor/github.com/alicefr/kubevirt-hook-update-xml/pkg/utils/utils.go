package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	types "github.com/alicefr/kubevirt-hook-update-xml/pkg/types"
	vmSchema "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/log"
)

func MergeKubeVirtXMLWithProvidedXML(file string, vmiJSON []byte) ([]byte, error) {
	vmiSpec := vmSchema.VirtualMachineInstance{}
	err := json.Unmarshal(vmiJSON, &vmiSpec)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to unmarshal given VMI spec: %s", vmiJSON)
		panic(err)
	}
	uid := string(vmiSpec.ObjectMeta.UID)
	log.Log.Infof("UID: %s", uid)

	return ioutil.ReadFile(fmt.Sprintf("%s/%s", types.Dir, file))
}
