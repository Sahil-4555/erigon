// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package shuffling_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/erigontech/erigon/cl/clparams"
	"github.com/erigontech/erigon/cl/phase1/core/state"
	"github.com/erigontech/erigon/cl/phase1/core/state/raw"
	"github.com/erigontech/erigon/cl/phase1/core/state/shuffling"
	"github.com/erigontech/erigon/cl/utils"
	"github.com/erigontech/erigon/cl/utils/eth2shuffle"
)

func BenchmarkLambdaShuffledIndex(b *testing.B) {
	keccakOptimized := utils.OptimizedSha256NotThreadSafe()
	eth2ShuffleHash := func(data []byte) []byte {
		hashed := keccakOptimized(data)
		return hashed[:]
	}
	seed := [32]byte{2, 35, 6}

	for b.Loop() {
		eth2shuffle.PermuteIndex(eth2ShuffleHash, uint8(clparams.MainnetBeaconConfig.ShuffleRoundCount), 10, 1000, seed)
	}
}

// Faster by ~40%, the effects of it will be felt mostly on computation of the proposer index.
func BenchmarkErigonShuffledIndex(b *testing.B) {
	s := state.New(&clparams.MainnetBeaconConfig)
	keccakOptimized := utils.OptimizedSha256NotThreadSafe()

	seed := [32]byte{2, 35, 6}
	preInputs := shuffling.ComputeShuffledIndexPreInputs(s.BeaconConfig(), seed)

	for b.Loop() {
		shuffling.ComputeShuffledIndex(s.BeaconConfig(), 10, 1000, seed, preInputs, keccakOptimized)
	}
}

func TestShuffling(t *testing.T) {
	s := raw.GetTestState()
	idx, err := shuffling.ComputeProposerIndex(s, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, [32]byte{1})
	require.NoError(t, err)
	require.Equal(t, uint64(2), idx)
}

// TestProposerDrift reproduces the scenario from #21582: two states exist for
// the same validator set — a historical snapshot with original balances, and a
// head state where balances have drifted. Using the head state for a historical
// epoch produces different proposers than using the correct historical state.
func TestProposerDrift(t *testing.T) {
	historicalState := raw.GetTestState()
	headState := raw.GetTestState()
	indices := []uint64{1, 2, 3, 4, 5, 6, 7, 8}
	seed := [32]byte{1}

	// Baseline: both states produce the same proposer
	expected, err := shuffling.ComputeProposerIndex(historicalState, indices, seed)
	require.NoError(t, err)

	headResult, err := shuffling.ComputeProposerIndex(headState, indices, seed)
	require.NoError(t, err)
	require.Equal(t, expected, headResult, "identical states must produce identical proposers")

	// Simulate head advancing: mutate effective balance on the head state only
	headValidator, err := headState.ValidatorForValidatorIndex(int(expected))
	require.NoError(t, err)
	headValidator.SetEffectiveBalance(0)

	// Head state now produces a different (drifted) proposer — this is the bug
	drifted, err := shuffling.ComputeProposerIndex(headState, indices, seed)
	require.NoError(t, err)
	require.NotEqual(t, expected, drifted, "head state with changed balances must produce a different proposer")

	// Historical state is unchanged — using it gives the correct, stable result
	stable, err := shuffling.ComputeProposerIndex(historicalState, indices, seed)
	require.NoError(t, err)
	require.Equal(t, expected, stable, "historical state must produce the same proposer regardless of head changes")
}

