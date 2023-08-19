package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"email-wizard/data/utils"

	pb "email-wizard/data/grpc"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50100, "The server port")
)

type server struct {
	pb.UnimplementedDatabaseHelperServer
}

func prepare_error_response(err_msg string) (*pb.Response, error) {
	message := make(map[string]interface{})
	message["err_msg"] = err_msg
	message_str, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: string(message_str)}, nil
}

func prepare_normal_response(message map[string]interface{}) (*pb.Response, error) {
	message_str, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: string(message_str)}, nil
}

func load_json_value(json_str string) (map[string]interface{}, error) {
	var value map[string]interface{}
	err := json.Unmarshal([]byte(json_str), &value)
	return value, err
}

func (s *server) AddRow(ctx context.Context, in *pb.AddRowRequest) (*pb.Response, error) {
	table := in.GetTable()
	row_str := in.GetRow()
	message := make(map[string]interface{})
	row, err := load_json_value(row_str)
	if err != nil {
		err_msg := fmt.Sprintf("fail to convert row from JSON, got %s", row_str)
		return prepare_error_response(err_msg)
	}
	pk_values, err := utils.AddRow(row, table)
	if err != nil {
		err_msg := err.Error()
		return prepare_error_response(err_msg)
	}
	message["pk"] = pk_values
	message["err_msg"] = ""
	return prepare_normal_response(message)
}

func (s *server) UpdateValue(ctx context.Context, in *pb.UpdateValueRequest) (*pb.Response, error) {
	table := in.GetTable()
	col_values, err := load_json_value(in.GetValueMap())
	var column string
	var value interface{}
	if err != nil {
		err_msg := fmt.Sprintf("fail to convert from JSON, got %s", in.GetValueMap())
		return prepare_error_response(err_msg)
	}
	for col, val := range col_values {
		column = col
		value = val
		break
	}
	condition, err := load_json_value(in.GetCondition())
	if err != nil {
		err_msg := fmt.Sprintf("fail to convert from JSON, got %s", in.GetCondition())
		return prepare_error_response(err_msg)
	}
	message := make(map[string]interface{})
	err = utils.UpdateValue(column, value, condition, table)
	if err != nil {
		return prepare_error_response(err.Error())
	}
	message["err_msg"] = ""
	return prepare_normal_response(message)
}

func (s *server) DeleteRows(ctx context.Context, in *pb.DeleteRowRequest) (*pb.Response, error) {
	table := in.GetTable()
	condition, err := load_json_value(in.GetCondition())
	if err != nil {
		err_msg := fmt.Sprintf("fail to convert from JSON, got %s", in.GetCondition())
		return prepare_error_response(err_msg)
	}
	message := make(map[string]interface{})
	err = utils.DeleteRows(condition, table)
	if err != nil {
		return prepare_error_response(err.Error())
	}
	message["err_msg"] = ""
	return prepare_normal_response(message)
}

func (s *server) Query(ctx context.Context, in *pb.QueryRequest) (*pb.Response, error) {
	table := in.GetTable()
	columns := in.GetColumns()
	message := make(map[string]interface{})
	condition, err := load_json_value(in.GetCondition())
	if err != nil {
		return prepare_error_response(err.Error())
	}
	values, err := utils.Query(columns, condition, table)
	if err != nil {
		err_msg := err.Error()
		return prepare_error_response(err_msg)
	}
	message["values"] = values
	message["err_msg"] = ""
	return prepare_normal_response(message)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDatabaseHelperServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
