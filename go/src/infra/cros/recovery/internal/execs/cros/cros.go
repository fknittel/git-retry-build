// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cros

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.chromium.org/luci/common/errors"

	"infra/cros/recovery/internal/execs"
	"infra/cros/recovery/internal/log"
	"infra/cros/recovery/internal/retry"
)

const (
	// Values presented as the string of the hex without 0x to match
	// representation in sysfs (idVendor/idProduct).

	// Servo's DUT side HUB vendor id
	SERVO_DUT_HUB_VID = "04b4"
	// Servo's DUT side HUB product id
	SERVO_DUT_HUB_PID = "6502"
	// Servo's DUT side NIC vendor id
	SERVO_DUT_NIC_VID = "0bda"
	// Servo's DUT side NIC product id
	SERVO_DUT_NIC_PID = "8153"
)

const (
	// Time to wait a rebooting ChromeOS, in seconds.
	NormalBootingTime = 150
	// Command to extract release builder path from device.
	extactReleaseBuilderPathCommand = "cat /etc/lsb-release | grep CHROMEOS_RELEASE_BUILDER_PATH"
)

// releaseBuildPath reads release build path from lsb-release.
func releaseBuildPath(ctx context.Context, run execs.Runner) (string, error) {
	// lsb-release is set of key=value so we need extract right part from it.
	//  Example: CHROMEOS_RELEASE_BUILDER_PATH=board-release/R99-9999.99.99
	output, err := run(ctx, time.Minute, extactReleaseBuilderPathCommand)
	if err != nil {
		return "", errors.Annotate(err, "release build path").Err()
	}
	log.Debug(ctx, "Read value: %q.", output)
	p, err := regexp.Compile("CHROMEOS_RELEASE_BUILDER_PATH=([\\w\\W]*)")
	if err != nil {
		return "", errors.Annotate(err, "release build path").Err()
	}
	parts := p.FindStringSubmatch(output)
	if len(parts) < 2 {
		return "", errors.Reason("release build path: fail to read value from %s", output).Err()
	}
	return strings.TrimSpace(parts[1]), nil
}

const (
	extactReleaseBoardCommand = "cat /etc/lsb-release | grep CHROMEOS_RELEASE_BOARD"
	releaseBoardRegexp        = `CHROMEOS_RELEASE_BOARD=(\S+)`
)

// ReleaseBoard reads release board info from lsb-release.
func ReleaseBoard(ctx context.Context, r execs.Runner) (string, error) {
	output, err := r(ctx, time.Minute, extactReleaseBoardCommand)
	if err != nil {
		return "", errors.Annotate(err, "release board").Err()
	}
	compiledRegexp, err := regexp.Compile(releaseBoardRegexp)
	if err != nil {
		return "", errors.Annotate(err, "release board").Err()
	}
	matches := compiledRegexp.FindStringSubmatch(output)
	if len(matches) != 2 {
		return "", errors.Reason("release board: cannot find chromeos release board information").Err()
	}
	board := matches[1]
	log.Debug(ctx, "Release board: %q.", board)
	return board, nil
}

// uptime returns uptime of resource.
func uptime(ctx context.Context, run execs.Runner) (*time.Duration, error) {
	// Received value represent two parts where the first value represents the total number
	// of seconds the system has been up and the second value is the sum of how much time
	// each core has spent idle, in seconds. We are looking
	//  E.g.: 683503.88 1003324.85
	// Consequently, the second value may be greater than the overall system uptime on systems with multiple cores.
	out, err := run(ctx, time.Minute, "cat /proc/uptime")
	if err != nil {
		return nil, errors.Annotate(err, "uptime").Err()
	}
	log.Debug(ctx, "Read value: %q.", out)
	p, err := regexp.Compile("([\\d.]{6,})")
	if err != nil {
		return nil, errors.Annotate(err, "uptime").Err()
	}
	parts := p.FindStringSubmatch(out)
	if len(parts) < 2 {
		return nil, errors.Reason("uptime: fail to read value from %s", out).Err()
	}
	// Direct parse to duration.
	// Example: 683503.88s -> 189h51m43.88s
	dur, err := time.ParseDuration(fmt.Sprintf("%ss", parts[1]))
	return &dur, errors.Annotate(err, "get uptime").Err()
}

// IsPingable checks whether the resource is pingable
// TODO: Migrate usage from components.
func IsPingable(ctx context.Context, info *execs.ExecInfo, resourceName string, count int) error {
	return info.RunArgs.Access.Ping(ctx, resourceName, count)
}

// IsNotPingable checks whether the resource is not pingable
func IsNotPingable(ctx context.Context, info *execs.ExecInfo, resourceName string, count int) error {
	if err := info.RunArgs.Access.Ping(ctx, resourceName, count); err != nil {
		log.Debug(ctx, "Resource %s is not pingble, but expected.", resourceName)
		return nil
	}
	return errors.Reason("not pingable: is pingable").Err()
}

const (
	pingAttemptInteval = 5 * time.Second
	sshAttemptInteval  = 10 * time.Second
)

// WaitUntilPingable waiting resource to be pingable.
// TODO: Migrate usage from components.
func WaitUntilPingable(ctx context.Context, info *execs.ExecInfo, resourceName string, waitTime time.Duration, count int) error {
	log.Debug(ctx, "Start ping %q for the next %s.", resourceName, waitTime)
	return retry.WithTimeout(ctx, pingAttemptInteval, waitTime, func() error {
		return IsPingable(ctx, info, resourceName, count)
	}, "wait to ping")
}

// WaitUntilNotPingable waiting resource to be not pingable.
func WaitUntilNotPingable(ctx context.Context, info *execs.ExecInfo, resourceName string, waitTime time.Duration, count int) error {
	return retry.WithTimeout(ctx, pingAttemptInteval, waitTime, func() error {
		return IsNotPingable(ctx, info, resourceName, count)
	}, "wait to be not pingable")
}

// IsSSHable checks whether the resource is sshable
// TODO: Migrate usage from components.
func IsSSHable(ctx context.Context, run execs.Runner) error {
	_, err := run(ctx, time.Minute, "true")
	return errors.Annotate(err, "is sshable").Err()
}

// WaitUntilSSHable waiting resource to be sshable.
// TODO: Migrate usage from components.
func WaitUntilSSHable(ctx context.Context, run execs.Runner, waitTime time.Duration) error {
	log.Debug(ctx, "Start SSH check for the next %s.", waitTime)
	return retry.WithTimeout(ctx, sshAttemptInteval, waitTime, func() error {
		return IsSSHable(ctx, run)
	}, "wait to ssh access")
}

// hasOnlySingleLine determines if the given string is only one single line.
func hasOnlySingleLine(ctx context.Context, s string) bool {
	if s == "" {
		log.Debug(ctx, "The string is empty")
		return false
	}
	lines := strings.Split(s, "\n")
	if len(lines) != 1 {
		log.Debug(ctx, "Found %d lines in the string.", len(lines))
		return false
	}
	return true
}

const (
	// findFilePathByContentCmdGlob find the file path by the content.
	// ex: grep -l xxx $(find /xxx/xxxx -maxdepth 1 -name xxx)
	findFilePathByContentCmdGlob = "grep -l %s $(find %s -maxdepth 1 -name %s)"
)

// FindSingleUsbDeviceFSDir find the common parent directory where the unique device with VID and PID is enumerated by file system.
//
//   1) Get path to the unique idVendor file with VID
//   2) Get path to the unique idProduct file with PID
//   3) Get directions of both files and compare them
//
// @param basePath: Path to the directory where to look for the device.
// @param vid: Vendor ID of the looking device.
// @param pid: Product ID of the looking device.
//
// @returns: path to the folder of the device.
func FindSingleUsbDeviceFSDir(ctx context.Context, r execs.Runner, basePath string, vid string, pid string) (string, error) {
	if basePath == "" {
		return "", errors.Reason("find single usb device file system directory: basePath is not provided").Err()
	}
	basePath += "*/"
	// find vid path:
	vidPath, err := r(ctx, time.Minute, fmt.Sprintf(findFilePathByContentCmdGlob, vid, basePath, "idVendor"))
	if err != nil {
		return "", errors.Annotate(err, "find single usb device file system directory").Err()
	} else if !hasOnlySingleLine(ctx, vidPath) {
		return "", errors.Reason("find single usb device file system directory: found more then one device with required VID: %s", vid).Err()
	}
	// find pid path:
	pidPath, err := r(ctx, time.Minute, fmt.Sprintf(findFilePathByContentCmdGlob, pid, basePath, "idProduct"))
	if err != nil {
		return "", errors.Annotate(err, "find single usb device file system directory").Err()
	} else if !hasOnlySingleLine(ctx, pidPath) {
		return "", errors.Reason("find single usb device file system directory: found more then one device with required PID: %s", pid).Err()
	}
	// If both files locates matched then we found our device.
	commDirCmd := fmt.Sprintf("LC_ALL=C comm -12 <(dirname %s) <(dirname %s)", vidPath, pidPath)
	commDir, err := r(ctx, time.Minute, commDirCmd)
	if err != nil {
		return "", errors.Annotate(err, "find single usb device file system directory").Err()
	} else if commDir == "" || commDir == "." {
		return "", errors.Reason("find single usb device file system directory: directory not found").Err()
	}
	return commDir, nil
}

const (
	// macAddressFileUnderNetFolderOfThePathGlob find NIC address from the nic path.
	// start finding the file name that contains both the /net/ and /address/ under the nic path folder.
	macAddressFileUnderNetFolderOfThePathGlob = "find %s/ | grep /net/ | grep /address"
	// Regex string to validate that MAC address is valid.
	// example of a correct format MAC address: f4:f5:e8:50:d1:cf
	macAddressVerifyRegexp = `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`
)

// ServoNICMacAddress read servo NIC mac address visible from DUT side.
//
// @param nic_path: Path to network device on the host
func ServoNICMacAddress(ctx context.Context, r execs.Runner, nicPath string) (string, error) {
	findNICAddressFileCmd := fmt.Sprintf(macAddressFileUnderNetFolderOfThePathGlob, nicPath)
	nicAddressFile, err := r(ctx, time.Minute, findNICAddressFileCmd)
	if err != nil {
		return "", errors.Annotate(err, "servo nic mac address").Err()
	} else if !hasOnlySingleLine(ctx, nicAddressFile) {
		return "", errors.Reason("servo nic mac address: found more then one nic address file").Err()
	}
	log.Info(ctx, "Found servo NIC address file: %q", nicAddressFile)
	macAddress, err := r(ctx, time.Minute, fmt.Sprintf("cat %s", nicAddressFile))
	if err != nil {
		return "", errors.Annotate(err, "servo nic mac address").Err()
	}
	macAddressRegexp, err := regexp.Compile(macAddressVerifyRegexp)
	if err != nil {
		return "", errors.Annotate(err, "servo nic mac address: regular expression for correct mac address cannot compile").Err()
	}
	if !macAddressRegexp.MatchString(macAddress) {
		log.Info(ctx, "Incorrect format of the servo nic mac address: %s", macAddress)
		return "", errors.Reason("servo nic mac address: incorrect format mac address found").Err()
	}
	log.Info(ctx, "Servo NIC MAC address visible from DUT: %s", macAddress)
	return macAddress, nil
}

const (
	// bootIDFile is the file path to the file that contains the boot id information.
	bootIDFilePath = "/proc/sys/kernel/random/boot_id"
	// noIDMessage is the default boot id file content if the device does not have a boot id.
	noIDMessage = "no boot_id available"
)

// BootID gets a unique ID associated with the current boot.
//
// @returns: A string unique to this boot if there is no error.
func BootID(ctx context.Context, run execs.Runner) (string, error) {
	bootId, err := run(ctx, 60*time.Second, fmt.Sprintf("cat %s", bootIDFilePath))
	if err != nil {
		return "", errors.Annotate(err, "boot id").Err()
	}
	if bootId == noIDMessage {
		log.Debug(ctx, "Boot ID: old boot ID not found, will be assumed empty.")
		return "", nil
	}
	return bootId, nil
}

const (
	// defaultPingRetryCount is the default ping retry count.
	defaultPingRetryCount = 2
	// waitDownRebootTime is the time the program will wait for the device to be down.
	waitDownRebootTime = 120 * time.Second
	// waitUpRebootTime is the time the program will wait for the device to be up after reboot.
	waitUpRebootTime = 240 * time.Second
)

// WaitForRestart will first wait the device to go down and then wait
// for the device to come up.
func WaitForRestart(ctx context.Context, info *execs.ExecInfo) error {
	// wait for it to be down.
	if waitDownErr := WaitUntilNotPingable(ctx, info, info.RunArgs.ResourceName, waitDownRebootTime, defaultPingRetryCount); waitDownErr != nil {
		log.Debug(ctx, "Wait For Restart: device shutdown failed.")
		return errors.Annotate(waitDownErr, "wait for restart").Err()
	}
	// wait down for servo device is successful, then wait for device
	// up.
	if waitUpErr := WaitUntilPingable(ctx, info, info.RunArgs.ResourceName, waitUpRebootTime, defaultPingRetryCount); waitUpErr != nil {
		return errors.Annotate(waitUpErr, "wait for restart").Err()
	}
	log.Info(ctx, "Device is up.")
	return nil
}

// TpmStatus is a data structure to represent the parse-version of the
// TPM Status.
type TpmStatus struct {
	statusMap map[string]string
	success   bool
}

// NewTpmStatus retrieves the TPM status for the DUT and returns the
// status values as a map.
func NewTpmStatus(ctx context.Context, run execs.Runner, timeout time.Duration) *TpmStatus {
	status, _ := run(ctx, timeout, "tpm_manager_client", "status", "--nonsensitive")
	log.Debug(ctx, "New Tpm Status :%q", status)
	statusItems := strings.Split(status, "\n")
	var ts = &TpmStatus{
		statusMap: make(map[string]string),
		// The uppercase on this string is deliberate.
		success: strings.Contains(strings.ToUpper(status), "STATUS_SUCCESS"),
	}
	// Following the logic in Labpack, if the TPM status string
	// contains 2 lines or fewer, we will return an empty map for the
	// TPM status values.
	if len(statusItems) > 2 {
		statusItems = statusItems[1 : len(statusItems)-1]
		for _, statusLine := range statusItems {
			item := strings.Split(statusLine, ":")[:]
			if item[0] == "" {
				continue
			}
			if len(item) == 1 {
				item = append(item, "")
			}
			for i, j := range item {
				item[i] = strings.TrimSpace(j)
			}
			ts.statusMap[item[0]] = item[1]
			// The labpack (Python) implementation checks whether the
			// string item[1] contains true of false in the string
			// form, and then explicitly converts that boolean
			// values. We do not attempt that here since the key and
			// value types for maps are strongly typed in Go-lang.
		}
	}
	return ts
}

// hasSuccess checks whether the TpmStatus includes success indicator
// or not.
func (tpmStatus *TpmStatus) hasSuccess() bool {
	return tpmStatus.success
}

// isOwned checks whether TPM has been cleared or not.
func (tpmStatus *TpmStatus) isOwned() (bool, error) {
	if len(tpmStatus.statusMap) == 0 {
		return false, errors.Reason("tpm status is owned: not initialized").Err()
	}
	return tpmStatus.statusMap["is_owned"] == "true", nil
}

// SimpleReboot executes a simple reboot command using a command
// runner for a DUT.
func SimpleReboot(ctx context.Context, run execs.Runner, timeout time.Duration, info *execs.ExecInfo) error {
	rebootCmd := "reboot"
	log.Debug(ctx, "Simple Rebooter : %s", rebootCmd)
	out, _ := run(ctx, timeout, rebootCmd)
	log.Debug(ctx, "Stdout: %s", out)
	if restartErr := WaitForRestart(ctx, info); restartErr != nil {
		return errors.Annotate(restartErr, "simple reboot").Err()
	}
	return nil
}
