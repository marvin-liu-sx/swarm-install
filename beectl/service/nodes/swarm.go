package nodes

import (
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
)

const (
	Status_Online  = 1
	Status_Offline = 0
)

// TODO:: 可以直接调用Bee官方SDK

func (this *BeeNode) GetPeers() int {
	reply, err := this.debugApiReply(fmt.Sprintf("%s/peers", this.DebugApiPort))
	if err != nil {
		return 0
	}
	Result := gjson.Parse(string(reply))
	raw := Result.Get("peers").Array()
	return len(raw)
}

func (this *BeeNode) GetStatus() string {
	reply, err := this.debugApiReply("health")
	if err != nil {
		return ""
	}
	Result := gjson.Parse(string(reply))
	raw := Result.Get("status").String()
	return raw
}

func (this *BeeNode) GetSettlements() string {
	reply, err := this.debugApiReply("settlements")
	if err != nil {
		return "0"
	}
	Result := gjson.Parse(string(reply))
	raw := Result.Get("totalReceived").String()
	return raw
}

func (this *BeeNode) GetUncashedAmount() int {
	reply, err := this.debugApiReply("chequebook/cheque")
	if err != nil {
		return 0
	}
	Result := gjson.Parse(string(reply))
	raws := Result.Get("lastcheques").Array()
	var (
		total int
	)
	for _, raw := range raws {
		amount := raw.Get("uncashedAmount").String()
		atoi, err := strconv.Atoi(amount)
		if err != nil {
			return 0
		}
		total += atoi
	}
	return total
}

func (this *BeeNode) GetTotalCheque() int {
	reply, err := this.debugApiReply("chequebook/cheque")
	if err != nil {
		return 0
	}
	Result := gjson.Parse(string(reply))
	raws := Result.Get("lastcheques").Array()

	return len(raws)
}

func (this *BeeNode) GetBalance() (totalBalance, availableBalance string) {
	reply, err := this.debugApiReply("chequebook/balance")
	if err != nil {
		return "", ""
	}
	Result := gjson.Parse(string(reply))

	return Result.Get("totalBalance").String(), Result.Get("availableBalance").String()
}

func (this *BeeNode) Cashout() ([]string, error) {
	reply, err := this.debugApiReply("chequebook/cheque")
	if err != nil {
		return nil, err
	}
	Result := gjson.Parse(string(reply))
	raws := Result.Get("lastcheques").Array()
	var rst []string
	for _, raw := range raws {
		hash := this.cashout(raw.Get("peer").String())
		rst = append(rst, hash)
		println(hash)
	}

	return rst, nil
}

func (this *BeeNode) cashout(peer string) string {
	reply, err := this.debugApiPost(fmt.Sprintf("chequebook/cashout/%s", peer))
	if err != nil {
		return "0"
	}
	Result := gjson.Parse(string(reply))
	raw := Result.Get("transactionHash").String()

	return raw
}
