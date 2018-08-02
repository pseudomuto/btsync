package main_test

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/pseudomuto/btsync/cmd/btsync"
)

func TestConfigure(t *testing.T) {
	tests := []struct {
		args       []string
		configFile []byte
		env        map[string]string
		values     map[string]string
	}{
		{
			values: map[string]string{
				"dir":             "config/partitions",
				"instance":        "",
				"project":         "",
				"service_account": "",
			},
		},
		{
			args: []string{
				"--dir", "some/dir",
				"--instance", "myinstance",
				"--project", "someproject",
				"--service_account", "/path/to/sa.json",
			},
			values: map[string]string{
				"dir":             "some/dir",
				"instance":        "myinstance",
				"project":         "someproject",
				"service_account": "/path/to/sa.json",
			},
		},
		{
			args: []string{
				"--project", "someproject",
			},
			env: map[string]string{
				"BTSYNC_DIR":             "some/dir",
				"BTSYNC_PROJECT":         "overridden by cli",
				"BTSYNC_SERVICE_ACCOUNT": "/path/to/sa.json",
			},
			values: map[string]string{
				"dir":             "some/dir",
				"instance":        "",
				"project":         "someproject",
				"service_account": "/path/to/sa.json",
			},
		},
		{
			args: []string{
				"--dir", "finalDir",
			},
			configFile: []byte(strings.Replace(`
			dir: config_dir
			instance: config_instance
			project: config_project
			service_account: config_service_account
			`,
				"\t",
				"",
				-1,
			)),
			env: map[string]string{
				"BTSYNC_DIR":             "some/dir",
				"BTSYNC_SERVICE_ACCOUNT": "/path/to/sa.json",
			},
			values: map[string]string{
				"dir":             "finalDir",
				"instance":        "config_instance",
				"project":         "config_project",
				"service_account": "/path/to/sa.json",
			},
		},
	}

	withTempEnv := func(env map[string]string, testFn func()) {
		if env != nil {
			for k, v := range env {
				os.Setenv(k, v)
			}

			defer func() {
				for k := range env {
					os.Unsetenv(k)
				}
			}()
		}

		testFn()
	}

	withConfigFile := func(data []byte, testFn func()) {
		if data != nil {
			assert.NoError(t, ioutil.WriteFile("btsync.yml", data, 0644))
			defer os.Remove("btsync.yml")
		}

		testFn()
	}

	for _, test := range tests {
		fs := flag.NewFlagSet("btsync", flag.ExitOnError)
		fs.String("dir", "config/partitions", "directory")
		fs.String("instance", "", "instance")
		fs.String("project", "", "project")
		fs.String("service_account", "", "sa")

		if test.args != nil {
			fs.Parse(test.args)
		}

		withTempEnv(test.env, func() {
			withConfigFile(test.configFile, func() {
				if test.configFile != nil {
					assert.NoError(t, main.Configure(fs))
				} else {
					main.Configure(fs)
				}

				for k, v := range test.values {
					assert.Equal(t, v, viper.Get(k))
				}
			})
		})
	}
}
