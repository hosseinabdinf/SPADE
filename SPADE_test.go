package SPADE

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/hosseinabdinf/SPADE/utils"
)

func TestSpade(t *testing.T) {
	for _, tc := range TestVector {
		fmt.Println(TestString("SPADE", tc))
		testSpade(t, tc.m, tc.n, tc.l, tc.v)
	}

	//To manually select test case
	//tc := TestVector[1]
	//fmt.Println(TestString("SPADE", tc))
	//testSpade(t, tc.n, tc.m, tc.l, tc.v)
}

func testSpade(t *testing.T, m int, n int, l int64, v int) {
	dummyData := utils.GenDummyData(m, n, l)
	// generate q, q = (2 ^ 128) + 1
	q := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)
	q.Add(q, big.NewInt(1))
	// generate g
	g := RandomElementInZMod(q)

	spade := NewSpade(q, g, n)
	var sks, pks, dks, res []*big.Int
	var ciphertexts [][]*big.Int

	t.Run("Setup", func(t *testing.T) {
		sks, pks = spade.Setup()
	})

	// initialize registration keys
	alphas := make([]*big.Int, m)
	regKeys := make([]*big.Int, m)

	// test only one user registration
	t.Run("Register", func(t *testing.T) {
		alphas[0] = RandomElementInZMod(q)
		regKeys[0] = spade.Register(alphas[0])
	})

	// do the registration for the rest of users
	for i := 1; i < m; i++ {
		alphas[i] = RandomElementInZMod(q)
		regKeys[i] = spade.Register(alphas[i])
	}

	t.Run("Encryption", func(t *testing.T) {
		ciphertexts = spade.Encrypt(pks, alphas[0], dummyData[0])
	})

	t.Run("KeyDerivation", func(t *testing.T) {
		dks = spade.KeyDerivation(0, v, sks, regKeys[0])
	})

	t.Run("Decryption", func(t *testing.T) {
		res = spade.Decrypt(dks, v, ciphertexts)
	})

	if len(res) != len(dummyData[0]) {
		t.Errorf("Decrypt failed: invalid length of decrypted message slice")
	}

	//fmt.Println("data: ", dummyData[0])
	//fmt.Println("res: ", res)
	utils.VerifyResults(dummyData[0], res, v)
}
