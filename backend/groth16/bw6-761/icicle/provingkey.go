package icicle_bw6761

import (
	"unsafe"

	groth16_bw6761 "github.com/consensys/gnark/backend/groth16/bw6-761"
	cs "github.com/consensys/gnark/constraint/bw6-761"
)

type deviceInfo struct {
	G1Device struct {
		A, B, K, Z unsafe.Pointer
	}
	DomainDevice struct {
		Twiddles, TwiddlesInv     unsafe.Pointer
		CosetTable, CosetTableInv unsafe.Pointer
	}
	G2Device struct {
		B unsafe.Pointer
	}
	DenDevice             unsafe.Pointer
	InfinityPointIndicesK []int
}

type ProvingKey struct {
	groth16_bw6761.ProvingKey
	*deviceInfo
}

func Setup(r1cs *cs.R1CS, pk *ProvingKey, vk *groth16_bw6761.VerifyingKey) error {
	return groth16_bw6761.Setup(r1cs, &pk.ProvingKey, vk)
}

func DummySetup(r1cs *cs.R1CS, pk *ProvingKey) error {
	return groth16_bw6761.DummySetup(r1cs, &pk.ProvingKey)
}