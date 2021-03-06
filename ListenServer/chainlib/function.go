package chainlib

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

type CommonFunc interface {
	SetTickAccountAddr(account string)

	Tick()

	ReTry(trx Transaction) bool

	Close()
}

//TransactionInfo ...
type TransactionInfo struct {
	ID  string `json:"id"`
	Trx struct {
		Trx struct {
			Actions []ActionInfo `json:"actions"`
		} `json:"trx"`
	} `json:"trx"`
}

//ActionInfo ...
type ActionInfo struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Data    struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Quantity string `json:"quantity"`
		Category string `json:"category"`
		Memo     string `json:"memo"`
	} `json:"data"`
}

//TransferInfo ...
type TransferInfo struct {
	Hash     string `json:"hash"`
	From     string `form:"from" json:"from" binding:"required"`
	To       string `form:"to" json:"to" bdinding:"required"`
	Quantity string `form:"quantity" json:"quantity" bdinding:"required"`
	Memo     string `form:"memo" json:"memo"`
}

//ChargeInfo ...
type ChargeInfo struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	BlockNum int64  `json:"blocknum"`
	Quantity string `json:"quantity"`
	Category string `json:"category"`
	Memo     string `json:"memo"`
}

// RandomSha256 ...
func RandomSha256() string {
	randStr := strconv.FormatUint(rand.Uint64(), 10)
	hash := sha256.New()
	hash.Write([]byte(randStr))
	return hex.EncodeToString(hash.Sum(nil))
}

// ParseTrxID ...
// string	返回TrxID
// error	返回json解析错误
func ParseTrxID(inStr string) (string, error) {
	js, err := simplejson.NewJson([]byte(inStr))
	if err != nil {
		return "", fmt.Errorf("simplejson error: %v\n", inStr)
	}
	return js.Get("transaction_id").MustString(), nil
}

// ExecShell ...
// string	返回标准输出
// error	返回标准错误
func ExecShell(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := string(stdout.Bytes())
	errStr := string(stderr.Bytes())
	if err != nil {
		return errStr, fmt.Errorf("Bash command error: %v\n", errStr)
	}
	return outStr, nil
}

// ClWalletUnlock ...
func ClWalletUnlock(password string) (string, error) {
	command := "cldatx wallet unlock " + " --password " + password
	return ExecShell(command)
}

// ClPushAction ...
func ClPushAction(account string, action string, data string) (string, error) {
	//Ensure that your wallet is unlocked before using it!
	command := "cldatx push action " + account + " " + action + " " + data + " -j " + " -f " + " -p " + account
	fmt.Println(command)
	return ExecShell(command)
}

// ClPushTransfer ...
// string	返回TrxID
func ClPushTransfer(account string, action string, trans TransferInfo) (string, error) {
	//Ensure that your wallet is unlocked before using it!
	js, _ := json.Marshal(trans)
	transStr := "'" + string(js) + "'"
	outStr, err := ClPushAction(account, action, transStr)
	if err != nil {
		return "", err
	}
	TrxID, err := ParseTrxID(outStr)
	if err != nil {
		return "", err
	}
	return TrxID, nil
}

// ClPushCharge ...
// string	返回TrxID
func ClPushCharge(account string, action string, charge ChargeInfo) (string, error) {
	//Ensure that your wallet is unlocked before using it!
	js, _ := json.Marshal(charge)
	chargeStr := "'" + string(js) + "'"
	outStr, err := ClPushAction(account, action, chargeStr)
	if err != nil {
		return "", err
	}
	TrxID, err := ParseTrxID(outStr)
	if err != nil {
		return "", err
	}
	return TrxID, nil
}

// ClGetTrxBlockNum ...
// int64	返回交易的block_num
func ClGetTrxBlockNum(trxID string) (int64, error) {
	command := "cldatx get transaction " + trxID
	outStr, err := ExecShell(command)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson([]byte(outStr))
	if err != nil {
		return 0, fmt.Errorf("simplejson error: %v\n", outStr)
	}
	return js.Get("block_num").MustInt64(), nil
}

// ClGetTrxInfo ...
func ClGetTrxInfo(trxID string) (*TransactionInfo, error) {
	command := "cldatx get transaction " + trxID
	outStr, err := ExecShell(command)
	if err != nil {
		return nil, err
	}
	var trx TransactionInfo
	if err := json.Unmarshal([]byte(outStr), &trx); err != nil {
		return nil, fmt.Errorf("ClGetTrxInfo Unmarshal error: %v\n", outStr)
	}
	return &trx, nil
}

// ClGetLIBNum ...
// int64	返回LIBNUM
func ClGetLIBNum() (int64, error) {
	command := "cldatx get info"
	outStr, err := ExecShell(command)
	if err != nil {
		return 0, err
	}
	js, err := simplejson.NewJson([]byte(outStr))
	if err != nil {
		return 0, fmt.Errorf("simplejson error: %v\n", outStr)
	}
	return js.Get("last_irreversible_block_num").MustInt64(), nil
}

// ClWaitIrreversible ...
func ClWaitIrreversible(blockNum int64) error {
	for {
		libNum, err := ClGetLIBNum()
		if err != nil { // 报错退出循环
			return err
		}
		if libNum > blockNum { // 不可逆退出循环
			return nil
		}
		time.Sleep(time.Duration(1) * time.Second) // 否则等待1秒
	}
}
