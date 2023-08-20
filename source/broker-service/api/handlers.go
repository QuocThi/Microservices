package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"broker/event"
	"broker/internal/logs"
	api_proto "broker/internal/randoms"
	"broker/pkg/utils"
)

// Broker is a test handler, just to make sure we can hit the broker from a web client
func (app *Server) Broker(w http.ResponseWriter, r *http.Request) {
	payload := utils.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Server) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.HandleError(w, app.Log, "marshal broker request failed", err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// app.logItemViaRPC(w, requestPayload.Log)
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	case "random":
		app.callRandom(w, requestPayload.Random)
		// app.callRandomRPC(w, requestPayload.Random)
	default:
		utils.HandleError(w, app.Log, "invalid request", err)
	}
}

func (app *Server) callRandom(w http.ResponseWriter, entry RandomPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal random request failed", err)
		return
	}

	randomServiceURL := "http://random-service/random"

	request, err := http.NewRequest("POST", randomServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.HandleError(w, app.Log, "prepare new random request failed", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.HandleError(w, app.Log, "make random request failed", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, fmt.Errorf("response status not accepted"))
		return
	}

	var jsonRes utils.JsonResponse
	json.NewDecoder(response.Body).Decode(&jsonRes)
	utils.WriteJSON(w, http.StatusAccepted, jsonRes)
}

// logItem logs an item by making an HTTP Post request with a JSON payload, to the logger microservice
func (app *Server) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal log request failed", err)
		return
	}

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.HandleError(w, app.Log, "prepare log request failed", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		utils.HandleError(w, app.Log, "make log request failed", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, err)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	err = utils.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.Log.Error(err, "failed to write log response")
	}
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Server) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal authenticate request failed", err)
		return
	}

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		utils.HandleError(w, app.Log, "make authenticate request failed", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.HandleError(w, app.Log, "make authenticate request failed", err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		utils.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService utils.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		utils.HandleError(w, app.Log, "marshal random request failed", err)
		return
	}

	if jsonFromService.Error {
		utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	err = utils.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.Log.Error(err, "failed to write authenticate response")
	}
}

// sendMail sends email by calling the mail microservice
func (app *Server) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal mail request failed", err)
		return
	}

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		utils.HandleError(w, app.Log, "prepare mail request failed", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.HandleError(w, app.Log, "make mail request failed", err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		utils.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	err = utils.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.Log.Error(err, "failed write send mail response")
	}
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Server) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		utils.HandleError(w, app.Log, "push message to queue failed", err)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Server) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

// logItemViaRPC logs an item by making an RPC call to the logger microservice
func (app *Server) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal rpc log request failed", err)
		return
	}

	rpcPayload := RPCPayload(l)

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		utils.HandleError(w, app.Log, "call rpc log request failed", err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: result,
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

type RandomRPCPayload struct {
	Data string
}

type RandomRPCResponse struct {
	Method string
	Data   string
}

// send request to random-service via RPC
func (app *Server) callRandomRPC(w http.ResponseWriter, p RandomPayload) {
	conn, err := rpc.Dial("tcp", "random-service:5002")
	if err != nil {
		utils.HandleError(w, app.Log, "marshal random rpc request failed", err)
		return
	}
	var rpcResponse RandomRPCResponse
	err = conn.Call("RPCServer.RandomRPC", RandomRPCPayload(p), &rpcResponse)
	if err != nil {
		utils.HandleError(w, app.Log, "call random rpc request failed", err)
		return
	}
	res := utils.JsonResponse{
		Error:   false,
		Message: "Call Random RPC successed",
		Data:    rpcResponse,
	}

	utils.WriteJSON(w, http.StatusAccepted, res)
}

func (app *Server) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.HandleError(w, app.Log, "marshal grpc log request failed", err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		utils.HandleError(w, app.Log, "dial grpc log request failed", err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		utils.HandleError(w, app.Log, "invoke WriteLog request failed", err)
		return
	}

	var payload utils.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Server) CallRandomGRPC(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload

	err := utils.ReadJSON(w, r, &request)
	if err != nil {
		utils.HandleError(w, app.Log, "marshal grpc random request failed", err)
	}

	conn, err := grpc.Dial("random-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		utils.HandleError(w, app.Log, "dial grpc random request failed", err)
		return
	}
	defer conn.Close()

	c := api_proto.NewRandomServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.RandomGPRC(ctx, &api_proto.RandomRequest{
		Data: request.Random.Data,
	})
	if err != nil {
		utils.HandleError(w, app.Log, "invoke grpc random request failed", err)
		return
	}

	response := utils.JsonResponse{
		Error:   res.Result,
		Message: fmt.Sprintf("Called %s", res.Method),
		Data:    res.Data,
	}

	utils.WriteJSON(w, http.StatusAccepted, response)
}
