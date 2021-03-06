// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package configuration

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/gae/service/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ufspb "infra/unifiedfleet/api/v1/models"
	ufsds "infra/unifiedfleet/app/model/datastore"
	"infra/unifiedfleet/app/util"
)

// MachineLSEPrototypeKind is the datastore entity kind for MachineLSEPrototypes.
const MachineLSEPrototypeKind string = "MachineLSEPrototype"

// MachineLSEPrototypeEntity is a datastore entity that tracks a platform.
type MachineLSEPrototypeEntity struct {
	_kind string   `gae:"$kind,MachineLSEPrototype"`
	ID    string   `gae:"$id"`
	Tags  []string `gae:"tags"`
	// ufspb.MachineLSEPrototype cannot be directly used as it contains pointer.
	MachineLSEPrototype []byte `gae:",noindex"`
}

// GetProto returns the unmarshaled MachineLSEPrototype.
func (e *MachineLSEPrototypeEntity) GetProto() (proto.Message, error) {
	var p ufspb.MachineLSEPrototype
	if err := proto.Unmarshal(e.MachineLSEPrototype, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func newMachineLSEPrototypeEntity(ctx context.Context, pm proto.Message) (ufsds.FleetEntity, error) {
	p := pm.(*ufspb.MachineLSEPrototype)
	if p.GetName() == "" {
		return nil, errors.Reason("Empty MachineLSEPrototype ID").Err()
	}
	machineLSEPrototype, err := proto.Marshal(p)
	if err != nil {
		return nil, errors.Annotate(err, "fail to marshal MachineLSEPrototype %s", p).Err()
	}
	return &MachineLSEPrototypeEntity{
		ID:                  p.GetName(),
		MachineLSEPrototype: machineLSEPrototype,
		Tags:                p.GetTags(),
	}, nil
}

// CreateMachineLSEPrototype creates a new machineLSEPrototype in datastore.
func CreateMachineLSEPrototype(ctx context.Context, machineLSEPrototype *ufspb.MachineLSEPrototype) (*ufspb.MachineLSEPrototype, error) {
	return putMachineLSEPrototype(ctx, machineLSEPrototype, false)
}

// UpdateMachineLSEPrototype updates machineLSEPrototype in datastore.
func UpdateMachineLSEPrototype(ctx context.Context, machineLSEPrototype *ufspb.MachineLSEPrototype) (*ufspb.MachineLSEPrototype, error) {
	return putMachineLSEPrototype(ctx, machineLSEPrototype, true)
}

// BatchUpdateMachineLSEPrototypes updates machineLSEPrototype in datastore.
//
// This is a non-atomic operation and doesnt check if the object already exists before
// update. Must be used within a Transaction where objects are checked before update.
// Will lead to partial updates if not used in a transaction.
func BatchUpdateMachineLSEPrototypes(ctx context.Context, machineLSEPrototypes []*ufspb.MachineLSEPrototype) ([]*ufspb.MachineLSEPrototype, error) {
	protos := make([]proto.Message, len(machineLSEPrototypes))
	updateTime := ptypes.TimestampNow()
	for i, p := range machineLSEPrototypes {
		p.UpdateTime = updateTime
		protos[i] = p
	}
	_, err := ufsds.PutAll(ctx, protos, newMachineLSEPrototypeEntity, true)
	if err == nil {
		return machineLSEPrototypes, err
	}
	return nil, err
}

// GetMachineLSEPrototype returns machineLSEPrototype for the given id from datastore.
func GetMachineLSEPrototype(ctx context.Context, id string) (*ufspb.MachineLSEPrototype, error) {
	pm, err := ufsds.Get(ctx, &ufspb.MachineLSEPrototype{Name: id}, newMachineLSEPrototypeEntity)
	if err == nil {
		return pm.(*ufspb.MachineLSEPrototype), err
	}
	return nil, err
}

func getMachineLSEPrototypeID(pm proto.Message) string {
	p := pm.(*ufspb.MachineLSEPrototype)
	return p.GetName()
}

// BatchGetMachineLSEPrototypes returns a batch of machine lse prototypes from datastore.
func BatchGetMachineLSEPrototypes(ctx context.Context, ids []string) ([]*ufspb.MachineLSEPrototype, error) {
	protos := make([]proto.Message, len(ids))
	for i, n := range ids {
		protos[i] = &ufspb.MachineLSEPrototype{Name: n}
	}
	pms, err := ufsds.BatchGet(ctx, protos, newMachineLSEPrototypeEntity, getMachineLSEPrototypeID)
	if err != nil {
		return nil, err
	}
	res := make([]*ufspb.MachineLSEPrototype, len(pms))
	for i, pm := range pms {
		res[i] = pm.(*ufspb.MachineLSEPrototype)
	}
	return res, nil
}

// ListMachineLSEPrototypes lists the machineLSEPrototypes
//
// Does a query over MachineLSEPrototype entities. Returns up to pageSize entities, plus non-nil cursor (if
// there are more results). pageSize must be positive.
func ListMachineLSEPrototypes(ctx context.Context, pageSize int32, pageToken string, filterMap map[string][]interface{}, keysOnly bool) (res []*ufspb.MachineLSEPrototype, nextPageToken string, err error) {
	// Passing -1 for query limit fetches all the entities from the datastore
	q, err := ufsds.ListQuery(ctx, MachineLSEPrototypeKind, -1, pageToken, filterMap, keysOnly)
	if err != nil {
		return nil, "", err
	}
	var nextCur datastore.Cursor
	err = datastore.Run(ctx, q, func(ent *MachineLSEPrototypeEntity, cb datastore.CursorCB) error {
		if keysOnly {
			machineLSEPrototype := &ufspb.MachineLSEPrototype{
				Name: ent.ID,
			}
			res = append(res, machineLSEPrototype)
		} else {
			pm, err := ent.GetProto()
			if err != nil {
				logging.Errorf(ctx, "Failed to UnMarshal: %s", err)
				return nil
			}
			res = append(res, pm.(*ufspb.MachineLSEPrototype))
		}
		if len(res) >= int(pageSize) {
			if nextCur, err = cb(); err != nil {
				return err
			}
			return datastore.Stop
		}
		return nil
	})
	if err != nil {
		logging.Errorf(ctx, "Failed to List MachineLSEPrototypes %s", err)
		return nil, "", status.Errorf(codes.Internal, ufsds.InternalError)
	}
	if nextCur != nil {
		nextPageToken = nextCur.String()
	}
	return
}

// DeleteMachineLSEPrototype deletes the machineLSEPrototype in datastore
func DeleteMachineLSEPrototype(ctx context.Context, id string) error {
	return ufsds.Delete(ctx, &ufspb.MachineLSEPrototype{Name: id}, newMachineLSEPrototypeEntity)
}

func putMachineLSEPrototype(ctx context.Context, machineLSEPrototype *ufspb.MachineLSEPrototype, update bool) (*ufspb.MachineLSEPrototype, error) {
	machineLSEPrototype.UpdateTime = ptypes.TimestampNow()
	pm, err := ufsds.Put(ctx, machineLSEPrototype, newMachineLSEPrototypeEntity, update)
	if err == nil {
		return pm.(*ufspb.MachineLSEPrototype), err
	}
	return nil, err
}

// ImportMachineLSEPrototypes creates or updates a batch of machine lse prototypes in datastore
func ImportMachineLSEPrototypes(ctx context.Context, lps []*ufspb.MachineLSEPrototype) (*ufsds.OpResults, error) {
	protos := make([]proto.Message, len(lps))
	utime := ptypes.TimestampNow()
	for i, m := range lps {
		m.UpdateTime = utime
		protos[i] = m
	}
	return ufsds.Insert(ctx, protos, newMachineLSEPrototypeEntity, true, true)
}

// GetMachineLSEPrototypeIndexedFieldName returns the index name
func GetMachineLSEPrototypeIndexedFieldName(input string) (string, error) {
	var field string
	input = strings.TrimSpace(input)
	switch strings.ToLower(input) {
	case util.TagFilterName:
		field = "tags"
	default:
		return "", status.Errorf(codes.InvalidArgument, "Invalid field name %s - field name for machineprototype are tag", input)
	}
	return field, nil
}
