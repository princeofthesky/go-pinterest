package model

var HTTP_SUCCESS = 0
var HTTP_CLIENT_ERROR = 1
var HTTP_SERVER_ERROR = 2

type Reponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
