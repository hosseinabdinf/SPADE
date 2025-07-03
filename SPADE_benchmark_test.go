package SPADE

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/hosseinabdinf/SPADE/utils"
)

func BenchmarkSpade(b *testing.B) {
	for _, tc := range TestVector {
		fmt.Println(TestString("SPADE", tc))
		benchmarkSpade(b, tc.m, tc.n, tc.l, tc.v)
	}
	//To manually select test case
	//tc := TestVector[2]
	//fmt.Println(TestString("SPADE", tc))
	//benchmarkSpade(b, tc.m, tc.n, tc.l, tc.v)
}

func benchmarkSpade(b *testing.B, m int, n int, l int64, v int) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode.")
	}

	dummyData := utils.GenDummyData(m, n, l)

	// generate q, q = (2 ^ 128) + 1
	q := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
	q.Add(q, big.NewInt(1))
	// generate g
	g := RandomElementInZMod(q)

	spade := NewSpade(q, g, n)
	var sks, pks, dks, res []*big.Int
	var ciphertexts [][]*big.Int

	b.Run("Setup", func(b *testing.B) {
		b.ResetTimer()
		sks, pks = spade.Setup()
	})

	// create dummy registration keys
	alphas := make([]*big.Int, m)
	regKeys := make([]*big.Int, m)

	b.Run("Register", func(b *testing.B) {
		for i := 0; i < m; i++ {
			alphas[i] = RandomElementInZMod(q)
			regKeys[i] = spade.Register(alphas[i])
		}
	})

	b.Run("Encryption", func(b *testing.B) {
		b.ResetTimer()
		ciphertexts = spade.Encrypt(pks, alphas[0], dummyData[0])
	})

	b.Run("KeyDerivation", func(b *testing.B) {
		b.ResetTimer()
		dks = spade.KeyDerivation(0, v, sks, regKeys[0])
	})

	b.Run("Decryption", func(b *testing.B) {
		b.ResetTimer()
		res = spade.Decrypt(dks, v, ciphertexts)
	})

	if len(res) != len(dummyData[0]) {
		b.Errorf("Decrypt failed: invalid length of decrypted message slice")
	}
}
