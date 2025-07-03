package models

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/hosseinabdinf/SPADE"
	pb "github.com/hosseinabdinf/SPADE/spadeProto"
	"github.com/hosseinabdinf/SPADE/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Analyst struct {
	id    int
	q     *big.Int
	g     *big.Int
	spade *SPADE.SPADE
	mpk   []*big.Int
}

func NewAnalyst(q, g *big.Int, mpk []*big.Int) *Analyst {
	return &Analyst{
		id:    1,
		q:     q,
		g:     g,
		spade: nil,
		mpk:   mpk,
	}
}

// StartAnalyst accept the configuration as in input and use SPADE to partially decrypt
// and get the results from a user's cipher
func StartAnalyst(config *Config, userID, queryValue int64) (int64, []*big.Int) {
	start := time.Now()

	pbHandler := NewPBHandler()

	log.Println(">>> Analyst starts connecting to the server..")
	addr := fmt.Sprintf("localhost:%d", utils.Port)
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxMsgSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// proto buffer init
	a := pb.NewCuratorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), config.TimeOut)
	defer cancel()

	// SPADE calls for Analysts here :)
	q, g, mpk, err := pbHandler.ReadPublicParams(a.GetPublicParams(ctx, &pb.PublicParamsReq{}))
	if err != nil {
		log.Fatalf("could not fetch the public parameters: %v", err)
	}

	// create a new Analyst
	// create an instance of SPADE with same public params of server
	analyst := NewAnalyst(q, g, mpk)
	spd := SPADE.NewSpade(q, g, config.MaxVecSize)
	analyst.spade = spd

	utils.PrintBigIntHex("q", q)
	utils.PrintBigIntHex("g", g)

	// send a query for value(1) and user-id(1) to the server
	req := &pb.AnalystReq{
		Id:    userID,
		Value: queryValue,
	}

	// get the unmarshal values
	dkv, cts, err := pbHandler.ReadDecryptionKey(a.Query(ctx, req))
	if err != nil {
		log.Fatalf("could not send the request: %v", err)
	}

	// partially decrypt the ciphertext vector using dkv to get the result for query
	results := spd.Decrypt(dkv, int(req.Value), cts)

	end := time.Now()
	elapsed := end.Sub(start)
	log.Printf("Analyst finished in %s", elapsed)
	return req.Value, results
}
