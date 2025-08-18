package bank

import (
	"europm/internal/bank/model"
	"europm/internal/util"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CallAccount(c *gin.Context) {
	var resp model.MbResponse
	requestId := fmt.Sprintf("%d", time.Now().UnixNano())
	cfg, err := util.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Println("error config.yaml:", err)
		resp.Acc.ResponseCode = "01"
		resp.Acc.ResponseDesc = "Customer not found"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	customerAcc := c.Query("customerAcc")

	var dataRes = model.AccountInfo{
		CustomerAcc: customerAcc,
	}

	if acc, found := util.FindAccount(customerAcc, cfg); found {
		dataRes.ResponseCode = "00"
		dataRes.ResponseDesc = "Successful"
		dataRes.CustomerName = acc.CustomerName
		dataRes.Rate = acc.Rate
	} else {
		dataRes.ResponseCode = "01"
		dataRes.ResponseDesc = "Customer not found"
	}

	priv, err := util.LoadPrivateKey(util.GetConfig("mb.privateKey"))

	if err != nil {
		log.Printf("Error get private key: %v", err)
		resp.Acc.ResponseCode = "01"
		resp.Acc.ResponseDesc = "Customer not found"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	sig, err := util.SignWithRSA(customerAcc, priv)
	if err != nil {
		log.Printf("Error SignWithRSA: %v", err)
		resp.Acc.ResponseCode = "01"
		resp.Acc.ResponseDesc = "Customer not found"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp = model.MbResponse{
		RequestId: requestId,
		Acc:       dataRes,
		Signature: sig,
	}
	c.JSON(http.StatusOK, resp)
}

func OrderUpdateHandler(c *gin.Context) {
	var mbResp model.MBOrderUpdateResponse
	var req model.MBOrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		mbResp.OrderRes.ResponseCode = "01"
		mbResp.OrderRes.ResponseDesc = "Thất bại. MB sẽ hoàn tiền lại cho khách hàng realtime"
		c.JSON(http.StatusInternalServerError, mbResp)
		return
	}

	// Tạo chuỗi cần verify (theo tài liệu: referenceNumber + customerAcc + amount + transDate)
	verifyStr := req.Order.ReferenceNumber + req.Order.CustomerAcc + req.Order.Amount + req.Order.TransDate

	// Load public key MB
	pub, err := util.LoadPublicKey(util.GetConfig("mb.publicKey"))
	if err != nil {
		log.Printf("Cannot load MB public key: %v", err)
		mbResp.OrderRes.ResponseCode = "01"
		mbResp.OrderRes.ResponseDesc = "Thất bại. MB sẽ hoàn tiền lại cho khách hàng realtime"
		c.JSON(http.StatusInternalServerError, mbResp)
		return
	}
	// Verify chữ ký từ MB
	if err := util.VerifySignature(verifyStr, req.Signature, pub); err != nil {
		log.Printf("Invalid signature: %v", err)
		mbResp.OrderRes.ResponseCode = "01"
		mbResp.OrderRes.ResponseDesc = "Thất bại. MB sẽ hoàn tiền lại cho khách hàng realtime"
		c.JSON(http.StatusInternalServerError, mbResp)
		return
	}

	// ----- Xử lý nghiệp vụ order update tại đây -----
	// Ví dụ response:
	transactionId := "TXN123455234234"
	responseCode := "00"
	responseDesc := "Order updated successfully"

	// Tạo chuỗi để ký trả về (transactionId + responseCode)
	signResp := transactionId + responseCode
	priv, err := util.LoadPrivateKey(util.GetConfig("mb.privateKey"))
	if err != nil {
		log.Printf("Cannot load partner private key: %v", err)
		mbResp.OrderRes.ResponseCode = "01"
		mbResp.OrderRes.ResponseDesc = "Thất bại. MB sẽ hoàn tiền lại cho khách hàng realtime"
		c.JSON(http.StatusInternalServerError, mbResp)
		return
	}
	signature, err := util.SignWithRSA(signResp, priv)
	if err != nil {
		log.Printf("Sign response failed: %v", err)
		mbResp.OrderRes.ResponseCode = "01"
		mbResp.OrderRes.ResponseDesc = "Thất bại. MB sẽ hoàn tiền lại cho khách hàng realtime"
		c.JSON(http.StatusInternalServerError, mbResp)
		return
	}

	// Trả về JSON
	resp := model.MBOrderUpdateResponse{
		RequestId: req.RequestId,
		Signature: signature,
	}
	resp.OrderRes.TransactionId = transactionId
	resp.OrderRes.ResponseCode = responseCode
	resp.OrderRes.ResponseDesc = responseDesc

	c.JSON(200, resp)
}
