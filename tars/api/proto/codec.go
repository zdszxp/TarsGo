package tars_api
import( 
	proto "github.com/golang/protobuf/proto"
	"errors"
)

func MarshalResponse(msg proto.Message) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Header: map[string]*Pair{},
		Body : string(data),
	}
	
	messageName := proto.MessageName(msg)
	if len(messageName) <= 0 {
		return nil, errors.New("Failed to marshal response : check import proto path")
	}

	resp.Header["Message"] = &Pair{Key : "Message"}
	resp.Header["Message"].Values = append(resp.Header["Message"].Values, messageName)

	respData, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respData, nil
}