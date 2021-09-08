package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	vmSchema "kubevirt.io/client-go/api/v1"
	"kubevirt.io/client-go/log"
	"kubevirt.io/kubevirt/pkg/hooks"
	hooksInfo "kubevirt.io/kubevirt/pkg/hooks/info"
	hooksV1alpha1 "kubevirt.io/kubevirt/pkg/hooks/v1alpha1"
	hooksV1alpha2 "kubevirt.io/kubevirt/pkg/hooks/v1alpha2"
)

const (
	onDefineDomainLoggingMessage = "Hook's OnDefineDomain callback method has been called"
	usage                        = `updater
  --version v1alpha1|v1alpha2
  --file /var/run/hooks/config/vm.xml`
)

var Dir = "/var/run/hooks/config"

type infoServer struct {
	Version string
}

func (s infoServer) Info(ctx context.Context, params *hooksInfo.InfoParams) (*hooksInfo.InfoResult, error) {
	log.Log.Info("Hook's Info method has been called")

	return &hooksInfo.InfoResult{
		Name: "update-xml",
		Versions: []string{
			s.Version,
		},
		HookPoints: []*hooksInfo.HookPoint{
			{
				Name:     hooksInfo.OnDefineDomainHookPointName,
				Priority: 0,
			},
		},
	}, nil
}

type v1alpha1Server struct {
	File string
}

type v1alpha2Server struct {
	File string
}

func (s v1alpha2Server) OnDefineDomain(ctx context.Context, params *hooksV1alpha2.OnDefineDomainParams) (*hooksV1alpha2.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	newDomainXML, err := MergeKubeVirtXMLWithProvidedXML(s.File, params.GetVmi())
	if err != nil {
		return nil, err
	}
	return &hooksV1alpha2.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}
func (s v1alpha2Server) PreCloudInitIso(_ context.Context, params *hooksV1alpha2.PreCloudInitIsoParams) (*hooksV1alpha2.PreCloudInitIsoResult, error) {
	return &hooksV1alpha2.PreCloudInitIsoResult{
		CloudInitData: params.GetCloudInitData(),
	}, nil
}

func (s v1alpha1Server) OnDefineDomain(ctx context.Context, params *hooksV1alpha1.OnDefineDomainParams) (*hooksV1alpha1.OnDefineDomainResult, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	newDomainXML, err := MergeKubeVirtXMLWithProvidedXML(s.File, params.GetVmi())
	if err != nil {
		return nil, err
	}
	return &hooksV1alpha1.OnDefineDomainResult{
		DomainXML: newDomainXML,
	}, nil
}

func MergeKubeVirtXMLWithProvidedXML(file string, vmiJSON []byte) ([]byte, error) {
	log.Log.Info(onDefineDomainLoggingMessage)
	vmiSpec := vmSchema.VirtualMachineInstance{}
	err := json.Unmarshal(vmiJSON, &vmiSpec)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to unmarshal given VMI spec: %s", vmiJSON)
		panic(err)
	}
	uid := string(vmiSpec.ObjectMeta.UID)
	log.Log.Infof("UID: %s", uid)

	return ioutil.ReadFile(fmt.Sprintf("%s/%s", Dir, file))
}

func main() {
	log.InitializeLogging("xml update")

	var version string
	pflag.StringVar(&version, "version", "", "hook version to use")

	var file string
	pflag.StringVar(&file, "config", "", "xml config file to use")
	pflag.Parse()

	socketPath := hooks.HookSocketsSharedDirectory + "/update.sock"
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Log.Reason(err).Errorf("Failed to initialized socket on path: %s", socket)
		log.Log.Error("Check whether given directory exists and socket name is not already taken by other file")
		panic(err)
	}
	defer os.Remove(socketPath)

	server := grpc.NewServer([]grpc.ServerOption{}...)

	if version == "" {
		panic(fmt.Errorf(usage))
	}
	if file == "" {
		panic(fmt.Errorf(usage))
	}

	hooksInfo.RegisterInfoServer(server, infoServer{Version: version})
	hooksV1alpha1.RegisterCallbacksServer(server, v1alpha1Server{File: file})
	hooksV1alpha2.RegisterCallbacksServer(server, v1alpha2Server{File: file})
	log.Log.Infof("Starting hook server exposing 'info' and 'v1alpha1' services on socket %s", socketPath)
	server.Serve(socket)
}
