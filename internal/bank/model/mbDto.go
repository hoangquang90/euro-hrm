package model

type Config struct {
	Accounts []Account `yaml:"accounts"`
}

type Account struct {
	CustomerAcc  string `yaml:"customerAcc"`
	CustomerName string `yaml:"customerName"`
	Rate         string `yaml:"rate"`
}

type MbResponse struct {
	RequestId string      `json:"requestId"`
	Acc       AccountInfo `json:"data"`
	Signature string      `json:"signature"`
}

type AccountInfo struct {
	CustomerAcc  string `json:"customerAcc"`
	CustomerName string `json:"customerName"`
	Rate         string `json:"rate"`
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}
type AdditionalData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MBOrderUpdateRequest struct {
	RequestId string         `json:"requestId"`
	Signature string         `json:"signature"`
	Order     MBOrderRequest `json:"data"`
}
type MBOrderRequest struct {
	ReferenceNumber string           `json:"referenceNumber"`
	Amount          string           `json:"amount"`
	CustomerAcc     string           `json:"customerAcc"`
	TransDate       string           `json:"transDate"`
	BillNumber      string           `json:"billNumber"`
	EndPointUrl     string           `json:"endPointUrl,omitempty"`
	UserName        string           `json:"userName,omitempty"`
	Rate            string           `json:"rate,omitempty"`
	CustomerName    string           `json:"customerName,omitempty"`
	AdditionalData  []AdditionalData `json:"additionalData,omitempty"`
}

type MBOrderUpdateResponse struct {
	RequestId string          `json:"requestId"`
	OrderRes  MBOrderResponse `json:"data"`
	Signature string          `json:"signature"`
}
type MBOrderResponse struct {
	TransactionId string `json:"transactionId"`
	ResponseCode  string `json:"responseCode"`
	ResponseDesc  string `json:"responseDesc"`
}
