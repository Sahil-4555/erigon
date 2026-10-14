// Copyright 2026 The Erigon Authors
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

package stagedsync

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/erigontech/erigon/db/dbservices"
	"github.com/erigontech/erigon/db/kv"
	"github.com/erigontech/erigon/execution/chain"
	"github.com/erigontech/erigon/execution/stagedsync/stages"
	"github.com/erigontech/erigon/execution/types"
)

type headerByNumberReader struct {
	dbservices.FullBlockReader
	header *types.Header
	err    error
}

func (r headerByNumberReader) HeaderByNumber(context.Context, kv.Getter, uint64) (*types.Header, error) {
	return r.header, r.err
}

func TestSerialForkValidation(t *testing.T) {
	t.Parallel()

	amsterdam := uint64(100)
	tests := []struct {
		name            string
		mode            stages.Mode
		experimentalBAL bool
		amsterdamTime   *uint64
		header          *types.Header
		err             error
		want            bool
	}{
		{name: "no amsterdam scheduled", mode: stages.ModeForkValidation, want: true},
		{name: "pre-amsterdam block", mode: stages.ModeForkValidation, amsterdamTime: &amsterdam, header: &types.Header{Time: 99}, want: true},
		{name: "amsterdam block", mode: stages.ModeForkValidation, amsterdamTime: &amsterdam, header: &types.Header{Time: 100}},
		{name: "experimental bal", mode: stages.ModeForkValidation, experimentalBAL: true},
		{name: "applying blocks", mode: stages.ModeApplyingBlocks},
		{name: "header missing", mode: stages.ModeForkValidation, amsterdamTime: &amsterdam},
		{name: "header error", mode: stages.ModeForkValidation, amsterdamTime: &amsterdam, header: &types.Header{Time: 99}, err: errors.New("read")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StageState{State: &Sync{mode: tt.mode}}
			cfg := ExecuteBlockCfg{
				chainConfig:     &chain.Config{AmsterdamTime: tt.amsterdamTime},
				blockReader:     headerByNumberReader{header: tt.header, err: tt.err},
				experimentalBAL: tt.experimentalBAL,
			}
			require.Equal(t, tt.want, serialForkValidation(t.Context(), s, cfg, nil, 1))
		})
	}
}
