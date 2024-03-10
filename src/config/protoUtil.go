package config

import (
	"fmt"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"gitlab.zixel.cn/go/framework/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var log = logger.Get()

func StructpbToMap(s *structpb.Struct) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range s.Fields {
		switch v := value.Kind.(type) {
		case *structpb.Value_NullValue:
			result[key] = nil
		case *structpb.Value_NumberValue:
			result[key] = v.NumberValue
		case *structpb.Value_StringValue:
			result[key] = v.StringValue
		case *structpb.Value_BoolValue:
			result[key] = v.BoolValue
		case *structpb.Value_StructValue:
			result[key] = StructpbToMap(v.StructValue)
		case *structpb.Value_ListValue:
			result[key] = ListValueToSlice(v.ListValue)
		}
	}

	return result
}

func ListValueToSlice(l *structpb.ListValue) []interface{} {
	result := make([]interface{}, len(l.Values))

	for i, v := range l.Values {
		switch value := v.Kind.(type) {
		case *structpb.Value_NullValue:
			result[i] = nil
		case *structpb.Value_NumberValue:
			result[i] = value.NumberValue
		case *structpb.Value_StringValue:
			result[i] = value.StringValue
		case *structpb.Value_BoolValue:
			result[i] = value.BoolValue
		case *structpb.Value_StructValue:
			result[i] = StructpbToMap(value.StructValue)
		case *structpb.Value_ListValue:
			result[i] = ListValueToSlice(value.ListValue)
		}
	}

	return result
}
func MapStringToStructpb(data map[string]string) (*structpb.Struct, error) {
	fields := make(map[string]*structpb.Value)

	for key, value := range data {
		fields[key] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: value},
		}
	}

	return &structpb.Struct{
		Fields: fields,
	}, nil
}
func MapToStructpb(data map[string]interface{}) (*structpb.Struct, error) {
	fields := make(map[string]*structpb.Value)

	for key, value := range data {
		fieldValue, err := ValueToStructpbValue(value)
		if err != nil {
			return nil, err
		}
		fields[key] = fieldValue
	}

	return &structpb.Struct{
		Fields: fields,
	}, nil
}

func ValueToStructpbValue(value interface{}) (*structpb.Value, error) {
	switch v := value.(type) {
	case nil:
		return &structpb.Value{Kind: &structpb.Value_NullValue{}}, nil
	case bool:
		return &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: v}}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(v.(int64))}}, nil
	case float32, float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(v.(float64))}}, nil
	case string:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: v}}, nil
	case []interface{}:
		listValue := make([]*structpb.Value, len(v))
		for i, item := range v {
			itemValue, err := ValueToStructpbValue(item)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			listValue[i] = itemValue
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: listValue}}}, nil
	case map[string]interface{}:
		structValue, err := MapToStructpb(v)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: structValue}}, nil
	case primitive.A:
		listValue := make([]*structpb.Value, len(v))
		for i, item := range v {
			itemValue, err := ValueToStructpbValue(item)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			listValue[i] = itemValue
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{Values: listValue}}}, nil
	default:
		return nil, fmt.Errorf("unsupported value type: %T,Vaule:%+v", v, value)
	}
}
