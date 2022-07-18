package model

var HTTP_SUCCESS = "ok"
var HTTP_CLIENT_ERROR = "Error client"
var HTTP_SERVER_ERROR = "Error server"

type Reponse struct {
	Code    string      `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
