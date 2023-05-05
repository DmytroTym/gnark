package cs

import (
	"crypto/sha256"
	"fmt"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/consensys/gnark/debug"
	"github.com/consensys/gnark/logger"
	"hash/fnv"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Bsb22CommitmentComputePlaceholder(mod *big.Int, input []*big.Int, output []*big.Int) error {
	if (len(os.Args) > 0 && (strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe"))) || debug.Debug {
		// usually we only run solver without prover during testing
		log := logger.Logger()
		log.Error().Msg("Augmented groth16 commitment hint not replaced. Proof will not be sound and verification will fail!")
		toHash := make([]byte, 0, (1+mod.BitLen()/8)*len(input))
		for _, in := range input {
			inBytes := in.Bytes()
			toHash = append(toHash, inBytes[:]...)
		}
		hsh := sha256.New().Sum(toHash)
		output[0].SetBytes(hsh)
		output[0].Mod(output[0], mod)

		return nil
	}
	return fmt.Errorf("placeholder function: to be replaced by commitment computation")
}

var maxNbCommitments int
var m sync.Mutex // TODO: Remove; this is only used in the frontend which is not required to be thread-safe

func RegisterBsb22CommitmentComputePlaceholder(index int) (hintName constraint.HintIds, err error) {
	hintName = constraint.HintIds{
		Name: "bsb22 commitment #" + strconv.Itoa(index),
	}
	// mimic solver.GetHintID
	hf := fnv.New32a()
	if _, err = hf.Write([]byte(hintName.Name)); err != nil {
		return
	}
	hintName.UUID = solver.HintID(hf.Sum32())

	m.Lock()
	if maxNbCommitments == index {
		maxNbCommitments++
		solver.RegisterNamedHint(Bsb22CommitmentComputePlaceholder, hintName.UUID)
	}
	m.Unlock()

	return
}
func init() {
	solver.RegisterHint(Bsb22CommitmentComputePlaceholder)
}
