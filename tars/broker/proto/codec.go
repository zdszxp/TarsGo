package tars_broker
import( 
	proto "github.com/golang/protobuf/proto"
	"errors"
)

func MarshalEvent(msg proto.Message) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	eventName := proto.MessageName(msg)
	if len(eventName) <= 0 {
		return nil, errors.New("Failed to marshal event : check import proto path")
	}

	event := &Event{
		Name: eventName,
		Data : data,
	}
	
	respData, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}

	return respData, nil
}