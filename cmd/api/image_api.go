package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pinterest/db"
	"go-pinterest/model"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	folderNFT = "/data/nft/images"
)

func main() {
	db.Init()
	defer db.Close()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	///users/:addr/nfts?last= & length=
	r.GET("/users/:addr/nfts", nftByUsers)
	r.GET("/nfts/sell", nftSell)
	r.GET("/nft/:tokenid/info", nftInfo)
	r.GET("/users/:addr/nonce", GetNonce)
	r.GET("/manga/metadata/:id", MetaData)

	r.POST("/users/sell_nft", SellNFT)
	r.POST("/manga/upload", UploadFile)
	r.Run(":9090")

}

func nftByUsers(c *gin.Context) {
	addrText := c.Param("addr")
	if !common.IsHexAddress(addrText) {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addrText  ", len(addrText))
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address is not hex address when upload file"})
		return
	}
	addr := common.HexToAddress(addrText)
	if addr.Hex() == "0x0000000000000000000000000000000000000000" {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addr  ", addr.Hex())
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address empty not allow when upload file"})
		return
	}
	lastText, exit := c.GetQuery("last")
	if !exit {
		lastText = "0"
	}
	lengthText, exit := c.GetQuery("length")
	if !exit {
		lengthText = "20"
	}
	last, _ := strconv.ParseInt(lastText, 10, 64)
	length, _ := strconv.Atoi(lengthText)
	data, err := db.GetNFTByAddress(addr, last, length)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addr  ", addr.Hex())
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error address empty not allow when upload file"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: data})
}

func nftSell(c *gin.Context) {
	lastText, exit := c.GetQuery("last")
	if !exit {
		lastText = "0"
	}
	lengthText, exit := c.GetQuery("length")
	if !exit {
		lengthText = "20"
	}
	last, _ := strconv.ParseInt(lastText, 10, 64)
	length, _ := strconv.Atoi(lengthText)
	data, err := db.GetNFTSell(last, length)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error address empty not allow when upload file"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: data})
}

func nftInfo(c *gin.Context) {
	id := c.Param("tokenid")
	tokenId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: fmt.Sprintf("Error when parse token id :%d ", tokenId)})
		return
	}
	dataId, err := db.GetDataIdFromTokenId(tokenId)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: fmt.Sprintf("Token id :%d not found", tokenId)})
		return
	}
	nftInfo, err := db.GetNFTInfo(dataId)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: fmt.Sprintf("Erorr when get token info  :%d  , err : %s ", tokenId, err)})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: nftInfo})
}

func UploadFile(c *gin.Context) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when read multipart form "})
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("form.File[\"file\"]    ", len(files))
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when read multipart form with form file = \"file\" not found"})
		return
	}
	file := files[0]
	address := form.Value["address"]
	if len(address) == 0 {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("form.Value[\"address\"]    ", len(address))
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address not found when upload file"})
		return
	}
	addrText := address[0]
	var fileData, _ = file.Open()
	fileBytes, err := ioutil.ReadAll(fileData)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println(err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error when read data multipart form "})
		return
	}
	harsher := md5.New()
	harsher.Write(fileBytes)
	fileName := hex.EncodeToString(harsher.Sum(nil))
	if !common.IsHexAddress(addrText) {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addrText  ", len(addrText))
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address is not hex address when upload file"})
		return
	}
	addr := common.HexToAddress(addrText)
	if addr.Hex() == "0x0000000000000000000000000000000000000000" {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addr  ", addr.Hex())
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address empty not allow when upload file"})
		return
	}
	nftWaiting := model.NFTWaiting{Owner: addr, Image: fileName}
	exit, err := db.CheckNFTWaitingExit(nftWaiting)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("CheckNFTWaitingExit  ", err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when check file exits"})
		return
	}
	if exit {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("CheckNFTWaitingExit  ", exit)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error NFT file exits , upload again not allow"})
		return
	}

	fileFullDir := folderNFT + "/" + fileName + ".png"
	f, err := os.OpenFile(fileFullDir, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("fileFullDir  ", fileFullDir, err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when save uploaded file "})
		return
	}
	// write this byte array to our temporary file
	_, err = f.Write(fileBytes)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("fileFullDir  ", fileFullDir, err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when write uploaded file "})
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("fileFullDir  ", fileFullDir, err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when write uploaded file "})
		return
	}

	check, err := db.AddNFTWaiting(nftWaiting)
	if err != nil {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("AddNFTWaiting  ", check, err)
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when add nft waiting upload "})
		return
	}
	// Upload the file to specific dst.
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Message: "Successfully Uploaded File"})
}

func MetaData(c *gin.Context) {
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	idText := c.Param("id")
	fmt.Println("MetaData idText", idText)
	tokenId, err := strconv.ParseInt(idText, 10, 64)

	if err != nil {
		return
	}
	fmt.Println("MetaData tokenId", tokenId)
	dataId, err := db.GetDataIdFromTokenId(tokenId)
	if err != nil {
		return
	}
	fmt.Println("MetaData dataId", dataId)
	nftInfo, err := db.GetNFTInfo(dataId)
	if err != nil {
		return
	}

	fmt.Println("MetaData nftInfo", nftInfo)
	nftMetaData := model.NFTMetaData{}
	nftMetaData.Image = "http://nft.skymeta.pro/images/" + nftInfo.Image + ".png"
	nftMetaData.Name = "1111111111"
	nftMetaData.Description = " Created by Mana App"
	// return that we have successfully uploaded our file!
	c.JSON(http.StatusOK, nftMetaData)
}

func SellNFT(c *gin.Context) {
	fmt.Println("SellNFT RequestURI", c.Request.RequestURI)
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("read data sell nft with err : ", err)
		return
	}

	nftOrder := model.NFTOrder{}
	err = json.Unmarshal(data, &nftOrder)
	if err != nil {
		fmt.Println("un marshal data sell nft with err : ", err)
		return
	}
	fmt.Println("nftOrder", nftOrder)
	data, _ = json.Marshal(nftOrder.NFT)
	sig := common.Hex2Bytes(nftOrder.Hash)
	publicKey, err := crypto.Ecrecover(nftOrder.NFT.Hash().Bytes(), sig)
	if err != nil {
		fmt.Println("un marshal data sell nft with err : ", err)
		return
	}
	var publicAddress common.Address
	copy(publicAddress[:], crypto.Keccak256(publicKey[1:])[12:])
	fmt.Println(publicAddress.Hex())
	c.String(http.StatusOK, publicAddress.Hex())
	nftInfo, err := db.GetNFTInfo(nftOrder.NFT.TokenId)
	if err != nil {
		fmt.Println("get nft info with err : ", err, "token Id ", nftOrder.NFT.TokenId)
		return
	}
	owner := nftInfo.Owner

	if bytes.Compare(owner.Bytes(), publicAddress.Bytes()) != 0 {
		fmt.Println("check nft owner with err : ", err, "got", publicAddress.Hex(), "wanted", owner)
	}
	nonce, err := db.GetNonceFromAddress(publicAddress)
	if err != nil {
		fmt.Println("get nonce addr with err : ", err, publicAddress.Hex())
		return
	}
	if nonce != nftOrder.NFT.Nonce {
		fmt.Println("error when check nonce addr with err : ", err, "got", nftOrder.NFT.Nonce, "wanted", nonce)
		return
	}

	db.UpdateNFTPriceInfo(nftOrder.NFT)
	db.SetNonceFromAddress(publicAddress, nonce)

	// return that we have successfully uploaded our file!
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Message: "Successfully set up Price NFT"})
}

func GetNonce(c *gin.Context) {
	addrText := c.Param("addr")
	if !common.IsHexAddress(addrText) {
		fmt.Println(c.Request.RequestURI)
		fmt.Println("addrText  ", len(addrText))
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_CLIENT_ERROR, Message: "Error address is not hex address when upload file"})
		return
	}
	fmt.Println("GetNonce RequestURI", c.Request.RequestURI)
	fmt.Println("MetaData idText", addrText)
	addr := common.HexToAddress(addrText)
	nonce, err := db.GetNonceFromAddress(addr)
	if err != nil {
		c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SERVER_ERROR, Message: "Error when data from db"})
		return
	}
	c.JSON(http.StatusOK, model.Reponse{Code: model.HTTP_SUCCESS, Data: nonce})
}
