// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package util

import (
	"encoding/binary"
	"fmt"
	"net"

	"go.chromium.org/luci/common/errors"

	ufspb "infra/unifiedfleet/api/v1/proto"

	crimsonconfig "go.chromium.org/luci/machine-db/api/config/v1"
)

// ParseVlan parses vlan to a list of IPs
func ParseVlan(vlan *crimsonconfig.VLAN) ([]*ufspb.IP, int, error) {
	ip, subnet, err := net.ParseCIDR(vlan.CidrBlock)
	if err != nil {
		return nil, 0, errors.Reason("invalid CIDR block %q for vlan %d", vlan.CidrBlock, vlan.GetId()).Err()
	}
	ipv4 := ip.Mask(subnet.Mask).To4()
	if ipv4 == nil {
		return nil, 0, errors.Reason("invalid IPv4 CIDR block %q for vlan %d", vlan.CidrBlock, vlan.GetId()).Err()
	}
	ones, _ := subnet.Mask.Size()
	length := 1 << uint32(32-ones)
	ips := make([]*ufspb.IP, length)
	startIP := binary.BigEndian.Uint32(ipv4)
	for i := 0; i < length; i++ {
		ips[i] = &ufspb.IP{
			Id:   getIPName(vlan.GetId(), startIP),
			Ipv4: startIP,
		}
		startIP++
	}
	return ips, length, nil
}

// FormatIP initialize an IP object
func FormatIP(vlan int64, ipAddress string, occupied bool) *ufspb.IP {
	ipv4, err := IPv4StrToInt(ipAddress)
	if err != nil {
		return nil
	}
	return &ufspb.IP{
		Id:       getIPName(vlan, ipv4),
		Ipv4:     ipv4,
		Vlan:     GetBrowserLabName(Int64ToStr(vlan)),
		Occupied: occupied,
	}
}

func getIPName(vlan int64, ipv4 uint32) string {
	return fmt.Sprintf("vlan-%d/%d", vlan, ipv4)
}

// IPv4StrToInt returns an uint32 address from the given ip address string.
func IPv4StrToInt(ipAddress string) (uint32, error) {
	ip := net.ParseIP(ipAddress)
	if ip != nil {
		ip = ip.To4()
	}
	if ip == nil {
		return 0, errors.Reason("invalid IPv4 address %q", ipAddress).Err()
	}
	return binary.BigEndian.Uint32(ip), nil
}

// IPv4IntToStr returns a string ip address
func IPv4IntToStr(ipAddress uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipAddress)
	return ip.String()
}

// parseCidrBlock returns a tuple of (cidr_block, capacity of this block)
func parseCidrBlock(subnet, mask string) (string, int) {
	maskIP := net.ParseIP(mask)
	maskAddr := maskIP.To4()
	ones, sz := net.IPv4Mask(maskAddr[0], maskAddr[1], maskAddr[2], maskAddr[3]).Size()
	return fmt.Sprintf("%s/%d", subnet, ones), 1 << uint32(sz-ones)
}
