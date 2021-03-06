package insert_one

import (
	"context"
	"encoding/json"
	"log"

	"github.com/sirius1024/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"

	n "server/common/names"
	"server/db/pkg/endpoint"
	"server/db/pkg/service"
	"server/db/pkg/types"
	kitendpoint "server/kit/endpoint"
	kithandler "server/kit/handler"
)

var (
	Request = n.RDBCvatTaskInsertOne
	Queue   = n.QDatabase
)

func Send(
	ctx context.Context,
	conn *rabbitmq.Connection,
	req RequestData,
) chan kitendpoint.Response {
	return kithandler.SendRequest(
		ctx,
		conn,
		Queue,
		request{
			Request: Request,
			Data:    req,
		},
		encodeRequest,
		decodeResponse,
		true,
	)
}

func Handle(
	eps endpoint.Endpoints,
	conn *rabbitmq.Connection,
	msg amqp.Delivery,
) {
	kithandler.HandleRequest(
		eps.CvatTaskInsertOne,
		conn,
		msg,
		decodeRequest,
		encodeResponse,
	)

}

type request struct {
	Request string      `json:"request"`
	Data    RequestData `json:"data"`
}

type RequestData = service.CvatTaskInsertOneRequestData

func encodeRequest(_ context.Context, pub *amqp.Publishing, req interface{}) (err error) {
	b, err := json.Marshal(req.(request))
	if err != nil {
		log.Println("cvat_task.insert_one.encodeRequest.Marshal", err)
	}
	pub.Body = b
	return
}

func decodeRequest(_ context.Context, deliv *amqp.Delivery) (interface{}, error) {
	var req request
	err := json.Unmarshal(deliv.Body, &req)
	if err != nil {
		log.Println("cvat_task.insert_one.decodeRequest.Unmarshal", err)
	}
	return req.Data, err
}

type ResponseData = types.CvatTask

func decodeResponse(_ context.Context, deliv *amqp.Delivery) (interface{}, error) {
	var res kitendpoint.Response
	var resData ResponseData
	err := json.Unmarshal(deliv.Body, &res)
	if err != nil {
		log.Println("cvat_task.insert_one.decodeResponse.Unmarshal(deliv.Body, &res)", err)
	}
	b, err := json.Marshal(res.Data)
	if err != nil {
		log.Println("cvat_task.insert_one.decodeResponse.Marshal(res.Data)", err)
	}
	err = json.Unmarshal(b, &resData)
	if err != nil {
		log.Println("cvat_task.insert_one.decodeResponse.Unmarshal(b, &resData)", err)
	}
	res.Data = resData
	return res, err
}

func encodeResponse(_ context.Context, pub *amqp.Publishing, resp interface{}) error {
	b, err := json.Marshal(resp.(kitendpoint.Response))
	if err != nil {
		log.Println("cvat_task.insert_one.encodeResponse.Marshal", err)
	}
	pub.Body = b
	return err
}
