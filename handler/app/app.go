package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/appengine/datastore"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Response resp
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Entity struct {
	Value string
}

func TestHandle(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Infof(c, "server start")

	k := datastore.NewKey(c, "testkind", "testStringID", 0, nil)
	e := new(Entity)
	e.Value = r.URL.Path

	if _, err := datastore.Put(c, k, e); err != nil {
		json, _ := json.Marshal(Response{Status: "NG", Message: err.Error()})
		fmt.Fprintf(w, string(json))
		return
	}

	log.Infof(c, "put end")

	json, _ := json.Marshal(Response{Status: "OK", Message: "hello world"})
	fmt.Fprintf(w, string(json))
}

func TestPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Infof(c, "test put start")
	k := datastore.NewKey(c, "testkind", "testStringID", 0, nil)
	e := new(Entity)
	e.Value = r.URL.Path

	if _, err := datastore.Put(c, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type gaeLog struct {
	path       string
	request_id string
	msg        []interface{}
}

// WarnwGAE ...
func WarnwGAE(gaectx context.Context, msg ...interface{}) {
	log.Warningf(gaectx, "%#v", mergeGAE(msg))
}

func mergeGAE(msg []interface{}) gaeLog {
	return gaeLog{path: "testpath", request_id: "testid", msg: msg}
}

func Logger(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	log.Infof(c, "test log info")

	lgr, _ := zap.NewProduction()
	lgr.Info("test zap info")
	lgr.Error("test zap error")

	log.Errorf(c, "test log error")

	a := OtherAPI{}

	WarnwGAE(c, "test warn log", a)

	fmt.Print("test fmt")
	w.WriteHeader(http.StatusOK)
}

func TestMirror(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	other, err := getother(r)
	if err != nil {
		log.Errorf(c, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	json, err := createResponse(200, "OK", other)
	if err != nil {
		log.Errorf(c, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, json)
}

// OtherAPI ...
type OtherAPI struct {
	Extension string `json:"Extension"`
	Word      string `json:"word"`
	Phrase    string `json:"phrase"`
	Time      int    `json:"time"`
	Ani       string `json:"ani"`
	Dnis      string `json:"dnis"`
	Omnisid   string `json:"omnisid"`
	Speaker   string `json:"speaker"`
	Callid    string `json:"callid"`
}

// HTTPResponse レスポンス定義
type HTTPResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func getother(r *http.Request) (OtherAPI, error) {
	var req OtherAPI
	//リクエストBODYを取得
	err := json.Unmarshal([]byte(readBody(r)), &req)
	if err != nil {
		return req, err
	}

	return req, nil
}

func readBody(r *http.Request) (body string) {
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body = bufbody.String()
	return body
}

func createResponse(httpstatus int, mes string, v interface{}) (JSON string, err error) {
	var res HTTPResponse

	res.Status = httpstatus
	res.Message = mes
	res.Data = v

	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)

	enc.SetEscapeHTML(false)
	enc.Encode(res)
	if err != nil {
		err = errors.New(fmt.Sprintln("json Marshal error: ", err))
		return "", err
	}

	return buffer.String(), nil
}

func CheckRequest(w http.ResponseWriter, r *http.Request) {
	for a, h := range r.Header {
		fmt.Println("header:", a, h)
	}

	fmt.Fprintf(w, "OK")
}
