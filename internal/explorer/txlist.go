package explorer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

	"github.com/chawin-a/wallet-monitor/internal/entity"
	"github.com/chawin-a/wallet-monitor/internal/utils"
	"github.com/ethereum/go-ethereum/common"
)

type TxListResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Result  []TxResponse `json:"result"`
}

type TxResponse struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	TxReceiptStatus  string `json:"txreceipt_status"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
}

func (txList *TxListResponse) ToEntity() *TxListEntityResponse {
	txEntityList := []*entity.Transaction{}
	for _, txResponse := range txList.Result {
		txEntityList = append(txEntityList, &entity.Transaction{
			BlockNumber:      utils.Must(strconv.ParseUint(txResponse.BlockNumber, 10, 64)),
			BlockHash:        common.HexToHash(txResponse.BlockHash),
			From:             common.HexToAddress(txResponse.From),
			To:               common.HexToAddress(txResponse.To),
			Gas:              utils.MustOk(new(big.Int).SetString(txResponse.Gas, 10)),
			GasPrice:         utils.MustOk(new(big.Int).SetString(txResponse.GasPrice, 10)),
			Nonce:            utils.Must(strconv.ParseUint(txResponse.Nonce, 10, 64)),
			Hash:             common.HexToHash(txResponse.Hash),
			Input:            common.FromHex(txResponse.Input),
			TransactionIndex: utils.Must(strconv.ParseUint(txResponse.TransactionIndex, 10, 64)),
			Value:            utils.MustOk(new(big.Int).SetString(txResponse.Value, 10)),
			Status:           utils.Must(strconv.ParseUint(txResponse.TxReceiptStatus, 10, 64)),
		})
	}
	return &TxListEntityResponse{
		Status:  txList.Status,
		Message: txList.Message,
		Result:  txEntityList,
	}
}

type TxListEntityResponse struct {
	Status  string
	Message string
	Result  []*entity.Transaction
}

func (e *Explorer) TxList(
	ctx context.Context,
	walletAddress common.Address,
	startBlock uint64,
	endBlock uint64,
) (*TxListEntityResponse, error) {
	endpoint, err := url.Parse(e.ApiURL + "/api")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("module", "account")
	params.Add("action", "txlist")
	params.Add("address", walletAddress.String())
	params.Add("startblock", fmt.Sprint(startBlock))
	params.Add("endblock", fmt.Sprint(endBlock))
	params.Add("sort", "asc")
	params.Add("apikey", e.ApiKey)
	endpoint.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	response.Body.Close()
	responseObject := &TxListResponse{}
	if err := json.Unmarshal(body, responseObject); err != nil {
		return nil, err
	}
	if responseObject.Status != "1" {
		return nil, errors.New(responseObject.Message)
	}
	return responseObject.ToEntity(), nil
}
