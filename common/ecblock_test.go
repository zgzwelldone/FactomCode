package common_test

import (
	"crypto/rand"
	"testing"

	"github.com/FactomProject/FactomCode/common"
	ed "github.com/FactomProject/ed25519"
	"github.com/davecgh/go-spew/spew"
)

func TestECBlockMarshal(t *testing.T) {
	ecb1 := common.NewECBlock()

	// build a CommitChain for testing
	cc := common.NewCommitChain()
	cc.Version = 0
	cc.MilliTime = &[6]byte{1, 1, 1, 1, 1, 1}
	cc.ChainIDHash.SetBytes(byteof(0xaa))
	cc.Weld.SetBytes(byteof(0xbb))
	cc.EntryHash.SetBytes(byteof(0xcc))
	cc.Credits = 11

	// make a key and sign the msg
	if pub, privkey, err := ed.GenerateKey(rand.Reader); err != nil {
		t.Error(err)
	} else {
		cc.ECPubKey = pub
		cc.Sig = ed.Sign(privkey, cc.CommitMsg())
	}

	// create an IncreaseBalance for testing
	pub := new([32]byte)
	copy(pub[:], byteof(0xaa))
	facTX := common.NewHash()
	facTX.SetBytes(byteof(0xbb))
	ib := common.MakeIncreaseBalance(pub, facTX, 12)

	// create a ECBlock for testing
	ecb1.Header.ECChainID.SetBytes(byteof(0x11))
	ecb1.Header.BodyHash.SetBytes(byteof(0x22))
	ecb1.Header.PrevHeaderHash.SetBytes(byteof(0x33))
	ecb1.Header.PrevLedgerKeyMR.SetBytes(byteof(0x44))
	ecb1.Header.DBHeight = 10
	ecb1.Header.HeaderExpansionArea = byteof(0x55)
	ecb1.Header.ObjectCount = 0

	// add the CommitChain to the ECBlock
	ecb1.AddEntry(cc)

	// add the IncreaseBalance
	ecb1.AddEntry(ib)

	ecb2 := common.NewECBlock()
	if p, err := ecb1.MarshalBinary(); err != nil {
		t.Error(err)
	} else {
		if err := ecb2.UnmarshalBinary(p); err != nil {
			t.Error(err)
		}
		t.Log(spew.Sdump(ecb2))
		if q, err := ecb2.MarshalBinary(); err != nil {
			t.Error(err)
		} else if string(p) != string(q) {
			t.Errorf("ecb1 = %x\n", p)
			t.Errorf("ecb2 = %x\n", q)
		}
	}
}

func byteof(b byte) []byte {
	r := make([]byte, 0, 32)
	for i := 0; i < 32; i++ {
		r = append(r, b)
	}
	return r
}
