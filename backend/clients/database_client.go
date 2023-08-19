// Package main implements a client for GetEmail service.
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "email-wizard/backend/clients/database_grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func load_json_value(json_str string) (map[string]interface{}, error) {
	var value map[string]interface{}
	err := json.Unmarshal([]byte(json_str), &value)
	return value, err
}

func prepare_grpc_client(timeout float32) (pb.DatabaseHelperClient, *grpc.ClientConn, 
		context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial("localhost:50100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	c := pb.NewDatabaseHelperClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	return c, conn, ctx, cancel, nil
}

func AddRow(row map[string]interface{}, table string) (map[string]interface{}, error) {
	c, conn, ctx, cancel, err := prepare_grpc_client(1)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer cancel()

	row_json, err := json.Marshal(row)
	if err != nil {
		return nil, err
	}

	r, err := c.AddRow(ctx, &pb.AddRowRequest{
		Row: string(row_json),
		Table: table},
	)
	if err != nil {
		return nil, err
	}

	response, err := load_json_value(r.GetMessage())
	if err != nil {
		return nil, err
	}
	err_msg, ok := response["err_msg"].(string)
	if !ok {
		return nil, fmt.Errorf("fail to read err_msg from %v", response)
	}
	if err_msg != "" {
		return nil, fmt.Errorf("MicroserviceError: %s", err_msg)
	}
	pk_values, ok := response["pk"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fail to read PK values from %v", response)
	}
	return pk_values, nil
}

func UpdateValue(column string, value interface{}, condition map[string]interface{}, table string) error {
	c, conn, ctx, cancel, err := prepare_grpc_client(1)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer cancel()

	col_value := map[string]interface{} {
		column: value,
	}
	col_value_json, err := json.Marshal(col_value)
	if err != nil {
		return err
	}
	condition_json, err := json.Marshal(condition)
	if err != nil {
		return nil
	}
	r, err := c.UpdateValue(ctx, &pb.UpdateValueRequest{
		ValueMap: string(col_value_json),
		Condition: string(condition_json),
		Table: table,
	})
	if err != nil {
		return err
	}
	response, err := load_json_value(r.GetMessage())
	if err != nil {
		return err
	}
	err_msg, ok := response["err_msg"].(string)
	if !ok {
		return fmt.Errorf("fail to read err_msg from %v", response)
	}
	if err_msg != "" {
		return fmt.Errorf("MicroserviceError: %s", err_msg)
	}
	return nil
}

func DeleteRows(condition map[string]interface{}, table string) error {
	c, conn, ctx, cancel, err := prepare_grpc_client(1)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer cancel()
	
	condition_json, err := json.Marshal(condition)
	if err != nil {
		return nil
	}
	r, err := c.DeleteRows(ctx, &pb.DeleteRowsRequest{
		Condition: string(condition_json),
		Table: table,
	})
	if err != nil {
		return err
	}
	response, err := load_json_value(r.GetMessage())
	if err != nil {
		return err
	}
	err_msg, ok := response["err_msg"].(string)
	if !ok {
		return fmt.Errorf("fail to read err_msg from %v", response)
	}
	if err_msg != "" {
		return fmt.Errorf("MicroserviceError: %s", err_msg)
	}
	return nil
}

func Query(columns []string, condition map[string]interface{}, table string) ([]map[string]interface{}, error) {
	c, conn, ctx, cancel, err := prepare_grpc_client(1)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer cancel()
	
	condition_json, err := json.Marshal(condition)
	if err != nil {
		return nil, err
	}
	r, err := c.Query(ctx, &pb.QueryRequest{
		Columns: columns,
		Condition: string(condition_json),
		Table: table,
	})
	if err != nil {
		return nil, err
	}
	response, err := load_json_value(r.GetMessage())
	if err != nil {
		return nil, err
	}
	err_msg, ok := response["err_msg"].(string)
	if !ok {
		return nil, fmt.Errorf("fail to read err_msg from %v", response)
	}
	if err_msg != "" {
		return nil, fmt.Errorf("MicroserviceError: %s", err_msg)
	}
	values_raw, ok := response["values"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("fail to read values from %v", response)
	}
	values := make([]map[string]interface{}, len(values_raw))
	for i := 0; i < len(values_raw); i++ {
		if values[i], ok = values_raw[i].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("fail to convert values to map[string]interface{}")
		}
	}
	return values, nil
}

func Reset() error {
	c, conn, ctx, cancel, err := prepare_grpc_client(1)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer cancel()

	r, err := c.ResetDB(ctx, &pb.EmptyRequest{})
	if err != nil {
		return err
	}
	response, err := load_json_value(r.GetMessage())
	if err != nil {
		return err
	}
	err_msg, ok := response["err_msg"].(string)
	if !ok {
		return fmt.Errorf("fail to read err_msg from %v", response)
	}
	if err_msg != "" {
		return fmt.Errorf("MicroserviceError: %s", err_msg)
	}
	return nil
}