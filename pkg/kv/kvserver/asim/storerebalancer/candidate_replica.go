// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package storerebalancer

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/kv/kvpb"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/allocator"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/asim/state"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvflowcontrol/rac2"
	"github.com/cockroachdb/cockroach/pkg/raft"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/errors"
)

// simulatorReplica is a replica that is being tracked as a potential candidate
// for rebalancing activities. It maintains a set of methods that enable
// querying its state and processing a rebalancing action if taken.
type simulatorReplica struct {
	rng   state.Range
	repl  state.Replica
	usage allocator.RangeUsageInfo
	state state.State
}

func newSimulatorReplica(repl state.Replica, s state.State) *simulatorReplica {
	rng, ok := s.Range(repl.Range())
	if !ok {
		return nil
	}
	sr := &simulatorReplica{
		rng:   rng,
		repl:  repl,
		usage: s.RangeUsageInfo(repl.Range(), repl.StoreID()),
		state: s,
	}
	return sr
}

// OwnsValidLease returns whether this replica is the current valid
// leaseholder.
func (sr *simulatorReplica) OwnsValidLease(context.Context, hlc.ClockTimestamp) bool {
	return sr.repl.HoldsLease()
}

// NodeID returns the Replica's NodeID.
func (sr *simulatorReplica) NodeID() roachpb.NodeID {
	return roachpb.NodeID(sr.repl.NodeID())
}

// StoreID returns the Replica's StoreID.
func (sr *simulatorReplica) StoreID() roachpb.StoreID {
	return roachpb.StoreID(sr.repl.StoreID())
}

// GetRangeID returns the Range ID.
func (sr *simulatorReplica) GetRangeID() roachpb.RangeID {
	return roachpb.RangeID(sr.repl.Range())
}

// RaftStatus returns the current raft status of the replica. It returns
// nil if the Raft group has not been initialized yet.
func (sr *simulatorReplica) RaftStatus() *raft.Status {
	return sr.state.RaftStatus(sr.rng.RangeID(), sr.repl.StoreID())
}

// GetCompactedIndex returns the compacted index of the raft log.
func (sr *simulatorReplica) GetCompactedIndex() kvpb.RaftIndex {
	// TODO(kvoli): We always return 1 here as RaftStatus is unimplemented.
	// When it is implmeneted, this may become variable.
	return 1
}

// LoadSpanConfig returns the authoritative range descriptor as well
// as the span config for the replica.
func (sr *simulatorReplica) LoadSpanConfig(ctx context.Context) (*roachpb.SpanConfig, error) {
	return sr.rng.SpanConfig(), nil
}

// Desc returns the authoritative range descriptor, acquiring a replica lock in
// the process.
func (sr *simulatorReplica) Desc() *roachpb.RangeDescriptor {
	return sr.rng.Descriptor()
}

// RangeUsageInfo returns usage information (sizes and traffic) needed by
// the allocator to make rebalancing decisions for a given range.
func (sr *simulatorReplica) RangeUsageInfo() allocator.RangeUsageInfo {
	return sr.usage
}

// AdminTransferLease transfers the LeaderLease to another replica.
func (sr *simulatorReplica) AdminTransferLease(
	ctx context.Context, target roachpb.StoreID, bypassSafetyChecks bool,
) error {
	if !sr.state.ValidTransfer(sr.repl.Range(), state.StoreID(target)) {
		return errors.Errorf(
			"unable to transfer lease for r%d to store %d, invalid transfer.",
			sr.repl.Range(), target)
	}

	if ok := sr.state.TransferLease(sr.repl.Range(), state.StoreID(target)); !ok {
		return errors.Errorf(
			"unable to transfer lease for r%d to store %d, application failed.",
			sr.repl.Range(), target)
	}

	return nil
}

func (sr *simulatorReplica) SendStreamStats(stats *rac2.RangeSendStreamStats) {}

// Replica returns the underlying kvserver replica, however when called from
// the simulator it only returns nil.
func (sr *simulatorReplica) Repl() *kvserver.Replica {
	return nil
}

// String implements the string interface.
func (sr *simulatorReplica) String() string {
	return sr.repl.Descriptor().String()
}

// GetStateRaftStatusFn returns a function that given a candidate replica, will
// return the raft status associated with it.
func GetStateRaftStatusFn(s state.State) func(replica kvserver.CandidateReplica) *raft.Status {
	return func(replica kvserver.CandidateReplica) *raft.Status {
		return s.RaftStatus(state.RangeID(replica.GetRangeID()), state.StoreID(replica.StoreID()))
	}
}
