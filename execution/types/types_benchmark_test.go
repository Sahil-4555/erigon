

package types

import (
	"testing"

	"github.com/erigontech/erigon/common"
	"github.com/erigontech/erigon/common/u256"
	"github.com/holiman/uint256"
)

// BenchmarkSigHash_DynamicFeeTx benchmarks the sigHash method for DynamicFeeTx
// with all fields populated.
func BenchmarkSigHash_DynamicFeeTx(b *testing.B) {
	to := common.HexToAddress("0x000000000000000000000000000000000000dead")

	tx := &DynamicFeeTransaction{
		CommonTx: CommonTx{
			Nonce:    3,
			To:       &to,
			Value:    uint256.NewInt(10),
			GasLimit: 25000,
			Data:     common.FromHex("5544"),
			V:        *uint256.NewInt(27),
			R:        *uint256.NewInt(123456789),
			S:        *uint256.NewInt(987654321),
		},
		ChainID: uint256.NewInt(1),
		TipCap:  uint256.NewInt(1_000_000_000), // 1 gwei
		FeeCap:  uint256.NewInt(2_000_000_000), // 2 gwei
		AccessList: AccessList{
			AccessTuple{
				Address: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				StorageKeys: []common.Hash{
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
			},
			AccessTuple{
				Address: common.HexToAddress("0x0000000000000000000000000000000000000002"),
				StorageKeys: []common.Hash{
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003"),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tx.SigningHash(tx.ChainID.ToBig())
	}
}

// BenchmarkSigHash_AccessListTx benchmarks the sigHash method for AccessListTx
// with all fields populated.
func BenchmarkSigHash_AccessListTx(b *testing.B) {
	to := common.HexToAddress("0x000000000000000000000000000000000000dead")

	tx := &AccessListTx{
		LegacyTx: LegacyTx{
			CommonTx: CommonTx{
				Nonce:    3,
				To:       &to,
				Value:    uint256.NewInt(10),
				GasLimit: 25000,
				Data:     common.FromHex("5544"),
				V:        *uint256.NewInt(27),
				R:        *uint256.NewInt(123456789),
				S:        *uint256.NewInt(987654321),
			},
			GasPrice: uint256.NewInt(2_000_000_000), // 2 gwei
		},
		ChainID: uint256.NewInt(1),
		AccessList: AccessList{
			AccessTuple{
				Address: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				StorageKeys: []common.Hash{
					common.HexToHash("0x01"),
					common.HexToHash("0x02"),
				},
			},
			AccessTuple{
				Address: common.HexToAddress("0x0000000000000000000000000000000000000002"),
				StorageKeys: []common.Hash{
					common.HexToHash("0x03"),
				},
			},
		},
	}

	chainID := tx.ChainID.ToBig()

	for b.Loop() {
		_ = tx.SigningHash(chainID)
	}
}

// BenchmarkSigHash_LegacyTx benchmarks the sigHash method for LegacyTx
// with all fields populated.
func BenchmarkSigHash_LegacyTx(b *testing.B) {
	to := common.HexToAddress("0x000000000000000000000000000000000000dead")

	tx := &LegacyTx{
		CommonTx: CommonTx{
			Nonce:    42,
			To:       &to,
			Value:    uint256.NewInt(1_000_000_000_000_000_000), // 1 ether
			GasLimit: 21000,
			Data:     []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
			V:        *uint256.NewInt(27),
			R:        *uint256.NewInt(123456789),
			S:        *uint256.NewInt(987654321),
		},
		GasPrice: uint256.NewInt(2_000_000_000), // 2 gwei
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tx.SigningHash(u256.Num1.ToBig())
	}
}

// BenchmarkSigHash_BlobTx benchmarks the sigHash method for BlobTx
// with all fields populated.
func BenchmarkSigHash_BlobTx(b *testing.B) {
	to := common.HexToAddress("0x000000000000000000000000000000000000dead")

	tx := &BlobTx{
		DynamicFeeTransaction: DynamicFeeTransaction{
			CommonTx: CommonTx{
				Nonce:    42,
				To:       &to,
				Value:    uint256.NewInt(1_000_000_000_000_000_000), // 1 ether
				GasLimit: 21000,
				Data:     []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
				V:        *uint256.NewInt(27),
				R:        *uint256.NewInt(123456789),
				S:        *uint256.NewInt(987654321),
			},
			ChainID: &u256.Num1,
			TipCap:  uint256.NewInt(1_000_000_000), // 1 gwei
			FeeCap:  uint256.NewInt(2_000_000_000), // 2 gwei
			AccessList: AccessList{
				{
					Address: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					StorageKeys: []common.Hash{
						common.HexToHash("0x01"),
						common.HexToHash("0x02"),
					},
				},
				{
					Address: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					StorageKeys: []common.Hash{
						common.HexToHash("0x03"),
					},
				},
			},
		},
		MaxFeePerBlobGas: uint256.NewInt(100_000_000),
		BlobVersionedHashes: []common.Hash{
			common.HexToHash("0x01"),
			common.HexToHash("0x02"),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tx.SigningHash(tx.ChainID.ToBig())
	}
}
