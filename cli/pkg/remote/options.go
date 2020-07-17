package remote

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"replicate.ai/cli/pkg/files"
)

// TODO: password login

// Options for connecting to an SSH server
type Options struct {
	Host           string
	Port           int
	Username       string
	PrivateKeys    []string
	ConnectTimeout time.Duration
}

// ParseHost creates an options struct from a string in the form [username@]hostname[:port]
func ParseHost(hostWithUsernameAndPort string) (options *Options, err error) {
	re := regexp.MustCompile("^(?:([^@]+)@)?([^:]+)(?::([0-9]+))?$")
	matches := re.FindStringSubmatch(hostWithUsernameAndPort)
	if len(matches) == 0 {
		return nil, fmt.Errorf("Invalid host. The host must be in the format [username@]hostname[:port]")
	}

	options = &Options{}
	options.Username = matches[1]
	options.Host = matches[2]
	if matches[3] != "" {
		options.Port, err = strconv.Atoi(matches[3])
		if err != nil {
			return options, err
		}
	}
	return options, nil
}

// SSHArgs returns SSH arguments/flags, except for host, and username
func (o *Options) SSHArgs() []string {
	args := []string{}

	if o.PrivateKeys != nil {
		for _, key := range o.PrivateKeys {
			args = append(args, "-i", key)
		}
	}
	for _, keyPath := range DefaultPrivateKeys() {
		absPath, err := files.ExpandUser(keyPath)
		if err != nil {
			panic("default private key has an error")
		}
		args = append(args, "-i", absPath)
	}

	if o.Port != 0 {
		args = append(args, "-p", fmt.Sprintf("%d", o.Port))
	}

	// TODO: host key verification
	args = append(args,
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "StrictHostKeyChecking=no",
	)

	// TODO: lower log level?
	args = append(args, "-o", "LogLevel=ERROR")

	return args
}