package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"

	pb "github.com/hosseinabdinf/SPADE/spadeProto"
	_ "github.com/mattn/go-sqlite3"
)

type PBHandler interface {
	ReadPublicParams(res *pb.PublicParamsResp, err error) (*big.Int, *big.Int, []*big.Int, error)
	ReadDecryptionKey(response *pb.AnalystResp, err error) ([]*big.Int, [][]*big.Int, error)
}

type pbHandler struct{}

func NewPBHandler() PBHandler {
	return &pbHandler{}
}

// ReadPublicParams convert the bytes data from server into required data type for @Client
func (pbh pbHandler) ReadPublicParams(response *pb.PublicParamsResp, err error) (*big.Int, *big.Int, []*big.Int, error) {
	if err != nil {
		return nil, nil, nil, err
	}

	// Unmarshal q, g
	q := new(big.Int)
	q.SetBytes(response.Q)

	g := new(big.Int)
	g.SetBytes(response.G)

	// Unmarshal mpk (slice of big.Int)
	mpk := make([]*big.Int, 0, len(response.Mpk))
	for _, mpkBytes := range response.Mpk {
		temp := new(big.Int)
		temp.SetBytes(mpkBytes)
		mpk = append(mpk, temp)
	}

	return q, g, mpk, nil
}

// ReadDecryptionKey convert the bytes data from server into required data type for @Analyst
func (pbh pbHandler) ReadDecryptionKey(response *pb.AnalystResp, err error) ([]*big.Int, [][]*big.Int, error) {
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal dkv
	dkv := make([]*big.Int, 0, len(response.Dkv))
	for _, dkvBytes := range response.Dkv {
		temp := new(big.Int)
		temp.SetBytes(dkvBytes)
		dkv = append(dkv, temp)
	}

	// Unmarshal ciphertext
	// Note: we encoded the ct = [n][2]*big.Int into ctBytes = [n*2][t]byte,
	// where t=len(ct_element.Bytes()), i.e. here will be ctBytes = [n*2][16]byte.
	// therefore, we have to unmarshal it in the opposite way to get [n][2]*big.Int
	cts := make([][]*big.Int, 0, len(response.Ciphertext)/2)
	for i := 0; i < len(response.Ciphertext); i += 2 {
		c0, c1 := new(big.Int), new(big.Int)
		c0.SetBytes(response.Ciphertext[i])
		c1.SetBytes(response.Ciphertext[i+1])
		cts = append(cts, []*big.Int{c0, c1})
	}

	return dkv, cts, nil
}

// DBHandler is an interface to work with sqlite3 database API
type DBHandler interface {
	CreateUsersCipherTable() error
	InsertUsersCipher(data *pb.UserReq) error
	GetUserReqById(userId int64) (*pb.UserReq, error)
}

type dbHandler struct {
	DbName string
	TbName string
}

func (d dbHandler) GetUserReqById(userId int64) (*pb.UserReq, error) {
	// Open the database
	db, err := sql.Open("sqlite3", d.DbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the query
	query := "SELECT regKey, ciphertext FROM " + d.TbName + " WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the user ID
	var regKey []byte
	var ctx []byte // we stored it as []bytes
	err = stmt.QueryRow(userId).Scan(&regKey, &ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No user found with the ID
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal the retrieved data back to UserReq
	var userReq pb.UserReq
	//err = proto.Unmarshal(rawData, &userReq)
	//if err != nil {
	//	return nil, err
	//}

	// Decode the retrieved ciphertext into [][]bytes
	var cipherText [][]byte
	err = json.Unmarshal(ctx, &cipherText)
	if err != nil {
		panic(err)
	}

	userReq.Id = userId
	userReq.RegKey = regKey
	userReq.Ciphertext = cipherText

	// Return the user request
	return &userReq, nil
}

func (d dbHandler) InsertUsersCipher(data *pb.UserReq) error {
	// Marshal user data to JSON bytes
	//row, err := proto.Marshal(data)
	//if err != nil {
	//	return err
	//}

	// open the database
	db, err := sql.Open("sqlite3", d.DbName)
	if err != nil {
		return err
	}
	defer db.Close()

	// we need to convert ciphertext from [][]bytes to []bytes to be able to store it using sqlite
	ctx, err := json.Marshal(data.Ciphertext)
	if err != nil {
		panic(err)
	}

	// Insert the data into the table
	insertQuery := "INSERT INTO " + d.TbName + " (id, regKey, ciphertext) VALUES (?, ?, ?)"
	_, err = db.Exec(insertQuery, data.Id, data.RegKey, ctx)
	if err != nil {
		return err
	}

	// OK
	return nil
}

func (d dbHandler) CreateUsersCipherTable() error {
	// open the database
	db, err := sql.Open("sqlite3", d.DbName)
	if err != nil {
		return err
	}
	defer db.Close()

	// create the table if it doesn't exist (schema migration)
	//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + d.TbName + ` (
	//	id INTEGER PRIMARY KEY,
	//	regKey BLOB,
	//	ciphertext BLOB
	//)`)
	// create the table if it doesn't exist (schema migration)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + d.TbName + ` (
		id INTEGER PRIMARY KEY,
		regKey BLOB,
		ciphertext BLOB
	)`)
	if err != nil {
		return err
	}

	// OK
	return nil
}

func NewDBHandler(database string, table string) DBHandler {
	return &dbHandler{
		DbName: database,
		TbName: table,
	}
}
