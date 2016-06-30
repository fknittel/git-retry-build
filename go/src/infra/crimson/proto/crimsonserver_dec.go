// Code generated by svcdec; DO NOT EDIT

package crimson

import (
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
)

type DecoratedCrimson struct {
	// Service is the service to decorate.
	Service CrimsonServer
	// Prelude is called in each method before forwarding the call to Service.
	// If Prelude returns an error, it is returned without forwarding the call.
	Prelude func(c context.Context, methodName string, req proto.Message) (context.Context, error)
}

func (s *DecoratedCrimson) CreateIPRange(c context.Context, req *IPRange) (*IPRangeStatus, error) {
	c, err := s.Prelude(c, "CreateIPRange", req)
	if err != nil {
		return nil, err
	}
	return s.Service.CreateIPRange(c, req)
}

func (s *DecoratedCrimson) ReadIPRange(c context.Context, req *IPRangeQuery) (*IPRanges, error) {
	c, err := s.Prelude(c, "ReadIPRange", req)
	if err != nil {
		return nil, err
	}
	return s.Service.ReadIPRange(c, req)
}

func (s *DecoratedCrimson) CreateHost(c context.Context, req *HostList) (*HostStatus, error) {
	c, err := s.Prelude(c, "CreateHost", req)
	if err != nil {
		return nil, err
	}
	return s.Service.CreateHost(c, req)
}

func (s *DecoratedCrimson) ReadHost(c context.Context, req *HostQuery) (*HostList, error) {
	c, err := s.Prelude(c, "ReadHost", req)
	if err != nil {
		return nil, err
	}
	return s.Service.ReadHost(c, req)
}
