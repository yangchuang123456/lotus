package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/filecoin-project/filecoin-ffi"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/sigs"
	_ "github.com/filecoin-project/lotus/lib/sigs/bls"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test_verifySignature(t *testing.T) {
	addr, err := address.NewFromString("t3qmprzvmkv2x6f4insp76utefvs6y6v4dn4w27c24h6cu577lqfhpl6eyjngrkpj7bk4txr43vcoav4fo7upq")
	assert.NoError(t, err)

	msg,err:=ioutil.ReadFile("message.json")
	assert.NoError(t,err)

	log.Println("the message is:",hex.EncodeToString(msg))
	
	message := types.Message{}
	err = json.Unmarshal(msg,&message)
	assert.NoError(t,err)
	log.Println("the message id is:",message.Cid())


	signature, err := hex.DecodeString("7b2254797065223a322c2244617461223a226b44434f757057352b4b48574350545a6a43583573585645324f693256746d577a56526f4c746c637a3571754c6956396e57456b484d3943504d336a2b774942446167324b344e4770465761623861482f536e4657376b6c616c636b575a53594d4261736c644771324c6f584e63525047446258515243366330454276363232227d")

/*	decode :=make([]string,0)
	err=json.Unmarshal(signature,&decode)
	assert.NoError(t, err)
	log.Println("the decode is:",decode)*/

	assert.NoError(t, err)
	sig := crypto.Signature{}
	err = json.Unmarshal(signature, &sig)
	assert.NoError(t, err)
	err = sigs.Verify(&sig, addr, msg)
	assert.NoError(t, err)
}

func InitLotusWrite(t *testing.T)(api.FullNode, jsonrpc.ClientCloser) {
	requestHeader := http.Header{}
	requestHeader.Set("Content-Type", "application/json")
	requestHeader.Set("Authorization", fmt.Sprintf("%v%v", "Bearer ", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiXX0.T9Ck5HeFxjq-EMGrBt3iLytJrcCiLVKqlKs26b5f8cs"))
	cli, stopper, err := client.NewFullNodeRPC("ws://192.168.131.2:1234/rpc/v0", requestHeader)
	assert.NoError(t,err)
	return cli,stopper
}

func Test_SignMessageAndPushMPool(t *testing.T){
	//api,close := InitLotusWrite(t)
}


