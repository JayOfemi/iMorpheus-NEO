package pow

import "blockChainMorp/types"

type IProofOfWork interface {

	// InitBlock inits a block in the proof-of-work
	InitBlock(b *types.Block)

	// Run performs a proof-of-work
	Run() (int, []byte)

	// Validate validates block's PoW
	Validate() bool
}