package gocsi

import (
	"context"
	"fmt"
	"os/user"
	"syscall"
	"testing"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/dell/gocsi/mock/service"
	log "github.com/sirupsen/logrus"
)

func TestRun(t *testing.T) {
	originalOsExit := osExit

	calledExit := make(chan struct{})
	Exit = func(_ int) {
		calledExit <- struct{}{}
	}

	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	endpoint := fmt.Sprintf("unix://%s/csi.sock", wd)

	defer func() {
		osExit = originalOsExit
		os.Unsetenv(EnvVarDebug)
		os.Unsetenv(EnvVarLogLevel)
		os.Unsetenv(EnvVarEndpoint)
		os.Unsetenv(EnvVarEndpointPerms)
		os.Unsetenv(EnvVarCredsCreateVol)
		os.Unsetenv(EnvVarCredsDeleteVol)
		os.Unsetenv(EnvVarCredsCtrlrPubVol)
		os.Unsetenv(EnvVarCredsCtrlrUnpubVol)
		os.Unsetenv(EnvVarCredsNodeStgVol)
		os.Unsetenv(EnvVarCredsNodePubVol)
		os.Unsetenv(EnvVarDisableFieldLen)
		os.Unsetenv(EnvVarRequireStagingTargetPath)
		os.Unsetenv(EnvVarRequireVolContext)
		os.Unsetenv(EnvVarCreds)
		os.Unsetenv(EnvVarSpecValidation)
		os.Unsetenv(EnvVarLoggingDisableVolCtx)
		os.Unsetenv(EnvVarPluginInfo)
		os.Unsetenv(EnvVarSerialVolAccessTimeout)
		os.Unsetenv(EnvVarSpecReqValidation)
		os.Unsetenv(EnvVarSpecRepValidation)
		os.Unsetenv(EnvVarEndpointUser)
		os.Unsetenv(EnvVarEndpointGroup)
		os.Unsetenv(EnvVarSerialVolAccessEtcdEndpoints)
	}()

	os.Setenv(EnvVarDebug, "true")
	os.Setenv(EnvVarLogLevel, "debug")
	os.Setenv(EnvVarEndpoint, endpoint)
	os.Setenv(EnvVarEndpointPerms, "0777")
	os.Setenv(EnvVarCredsCreateVol, "true")
	os.Setenv(EnvVarCredsDeleteVol, "true")
	os.Setenv(EnvVarCredsCtrlrPubVol, "true")
	os.Setenv(EnvVarCredsCtrlrUnpubVol, "true")
	os.Setenv(EnvVarCredsNodeStgVol, "true")
	os.Setenv(EnvVarCredsNodePubVol, "true")
	os.Setenv(EnvVarDisableFieldLen, "true")
	os.Setenv(EnvVarRequireStagingTargetPath, "true")
	os.Setenv(EnvVarRequireVolContext, "true")
	os.Setenv(EnvVarCreds, "true")
	os.Setenv(EnvVarSpecValidation, "false")
	os.Setenv(EnvVarLoggingDisableVolCtx, "true")
	os.Setenv(EnvVarPluginInfo, "true")
	os.Setenv(EnvVarSerialVolAccessTimeout, "10s")
	os.Setenv(EnvVarSpecReqValidation, "true")
	os.Setenv(EnvVarSpecRepValidation, "true")
	os.Setenv(EnvVarEndpointUser, user.Name)
	os.Setenv(EnvVarEndpointGroup, user.Gid)
	os.Setenv(EnvVarSerialVolAccessEtcdEndpoints, "http://127.0.0.1:2379")

	svc := service.NewServer()
	go Run(context.Background(), "Dell CSM Driver", "A Dell Container Storage Interface (CSI) Plugin", "", newMockStoragePlugin(svc, svc, svc))
	time.Sleep(5 * time.Second)

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		t.Fatalf("failed to send SIGINT: %v", err)
	}
	<-calledOsExit
}
