package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/libcontainerd"
	"github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/pkg/system"
)

var defaultDaemonConfigFile = os.Getenv("programdata") + string(os.PathSeparator) + "docker" + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "daemon.json"

// currentUserIsOwner checks whether the current user is the owner of the given
// file.
func currentUserIsOwner(f string) bool {
	return false
}

// setDefaultUmask doesn't do anything on windows
func setDefaultUmask() error {
	return nil
}

func getDaemonConfDir() string {
	return os.Getenv("PROGRAMDATA") + `\docker\config`
}

// notifySystem sends a message to the host when the server is ready to be used
func notifySystem() {
}

// setupConfigReloadTrap configures a Win32 event to reload the configuration.
func setupConfigReloadTrap(configFile string, flags *mflag.FlagSet, reload func(*daemon.Config)) {
	go func() {
		sa := syscall.SecurityAttributes{
			Length: 0,
		}
		ev := "Global\\docker-daemon-config-" + fmt.Sprint(os.Getpid())
		if h, _ := system.CreateEvent(&sa, false, false, ev); h != 0 {
			logrus.Debugf("Config reload - waiting signal at %s", ev)
			for {
				syscall.WaitForSingleObject(h, syscall.INFINITE)
				if err := daemon.ReloadConfiguration(configFile, flags, reload); err != nil {
					logrus.Error(err)
				}
			}
		}
	}()
}

func (cli *DaemonCli) getPlatformRemoteOptions() []libcontainerd.RemoteOption {
	return nil
}

// getLibcontainerdRoot gets the root directory for libcontainerd to store its
// state. The Windows libcontainerd implementation does not need to write a spec
// or state to disk, so this is a no-op.
func (cli *DaemonCli) getLibcontainerdRoot() string {
	return ""
}

func allocateDaemonPort(addr string) error {
	return nil
}
