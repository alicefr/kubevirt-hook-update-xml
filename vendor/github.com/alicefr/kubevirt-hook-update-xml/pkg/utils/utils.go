package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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

	mv, merr := mxj.NewMapXml(rawXML)
	if merr != nil {
		log.Log.Reason(merr).Errorf("Failed to unmarshal xml: %v", merr)
		return []byte{}, merr

	}
	var v interface{}
	v, err = mv.ValuesForPath("domain.metadata.kubevirt.uid")
	if err != nil {
		log.Log.Reason(err).Errorf("Failed parsing old uid in the xml value:%v : %v", v, err)
		return []byte{}, err
	}
	oldUID := v.([]interface{})[0].(string)
	log.Log.Infof("Replace old UID:%s with new UID:%s", oldUID, newUID)
	newXML := strings.ReplaceAll(string(rawXML), oldUID, newUID)
	return []byte(newXML), nil
}
