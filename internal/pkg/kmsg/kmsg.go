// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package kmsg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	"github.com/talos-systems/talos/pkg/constants"
)

// Setup configures the log package to write to the kernel ring buffer via
// /dev/kmsg.
func Setup(prefix string, withLogFile bool) error {
	kmsg, err := os.OpenFile("/dev/kmsg", os.O_RDWR|unix.O_CLOEXEC|unix.O_NONBLOCK|unix.O_NOCTTY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open /dev/kmsg: %w", err)
	}

	var writer io.Writer = kmsg

	if withLogFile {
		if err := os.MkdirAll(constants.DefaultLogPath, 0700); err != nil {
			return err
		}

		logPath := filepath.Join(constants.DefaultLogPath, "machined.log")

		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", logPath, err)
		}

		writer = io.MultiWriter(kmsg, f)
	}

	log.SetOutput(writer)
	log.SetPrefix(prefix + " ")
	log.SetFlags(0)

	return nil
}
