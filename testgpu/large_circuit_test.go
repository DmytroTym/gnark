package testgpu

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/logger"
	"github.com/consensys/gnark/test"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

const TEST_SIZE = 2

type LargeCircuitCommitment struct {
	P, Q [TEST_SIZE]frontend.Variable
	N    frontend.Variable `gnark:",public"`
}

func TestLargeCircuitInGpu(t *testing.T) {
	logger.Set(zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).With().Timestamp().Logger())
	assert := test.NewAssert(t)

	field := ecc.BN254.ScalarField()

	innerCcs, err := frontend.Compile(field, r1cs.NewBuilder, &LargeCircuitCommitment{})
	assert.NoError(err)
	innerPK, innerVK, err := groth16.Setup(innerCcs)
	assert.NoError(err)

	var p, q [TEST_SIZE]frontend.Variable
	for i := 0; i < TEST_SIZE; i++ {
		p[i] = 3
		q[i] = 5
	}
	// inner proof
	innerAssignment := &LargeCircuitCommitment{
		N: 15,
	}
	innerWitness, err := frontend.NewWitness(innerAssignment, field)
	assert.NoError(err)
	innerProof, err := groth16.Prove(innerCcs, innerPK, innerWitness)
	assert.NoError(err)
	innerPubWitness, err := innerWitness.Public()
	assert.NoError(err)
	err = groth16.Verify(innerProof, innerVK, innerPubWitness)
	assert.NoError(err)
}

func (c *LargeCircuitCommitment) Define(api frontend.API) error {
	for i := 0; i < TEST_SIZE; i++ {
		res := api.Mul(c.P[i], c.Q[i])
		api.AssertIsEqual(res, c.N)
	}

	commitment, err := api.Compiler().(frontend.Committer).Commit(c.P[0], c.Q[0], c.N)
	if err != nil {
		return err
	}

	api.AssertIsDifferent(commitment, 0)

	return nil
}