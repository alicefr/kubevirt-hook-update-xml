package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	types "github.com/alicefr/kubevirt-hook-update-xml/pkg/types"
	mxj "github.com/clbanning/mxj/v2"
	vmSchema "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/log"
)

const (
	uidSeparatorLeft  = "uid"
	uidSeparatorRight = `\u003c/uid\u003e`
)

type Map map[string]interface{}

func MergeKubeVirtXMLWithProvidedXML(file string, vmiJSON []byte) ([]byte, error) {
	vmiSpec := vmSchema.VirtualMachineInstance{}
	err := json.Unmarshal(vmiJSON, &vmiSpec)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to unmarshal given VMI spec: %s", vmiJSON)
		return []byte{}, err
	}
	newUID := string(vmiSpec.ObjectMeta.UID)
	var rawXML []byte
	rawXML, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", types.Dir, file))
	if err != nil {
		return []byte{}, err
	}

	var mv Map
	mv, err = mxj.NewMapXml(rawXML)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to unmarshal xml: %v", err)
		return []byte{}, err

	}

	if _, ok := objectMap["metadata"]; !ok {
		return []byte{}, fmt.Errorf("Failed metedata not found")

	}
	fmt.Printf("XXX metadata %s \n", string(objectMap["metadata"]))
	log.Log.Infof("New UID: %s", newUID)
	return rawXML, nil
}
