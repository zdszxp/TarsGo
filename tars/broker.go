package tars

import (
	"log"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/broker/redis"
)

func InitBroker() (err error) {
	//broker.DefaultBroker = redis.NewBroker(broker.Addrs("192.168.10.158", "6387"))
	broker.DefaultBroker = redis.NewBroker(broker.Addrs("user:foobared@192.168.10.158:6379/0"))
	
	err = broker.Init()
	if err != nil {
		log.Fatal("Broker Init error: %v", err)
	}else{
		log.Print("Broker Init successfully")
	}

	if err = broker.Connect(); err != nil {
		log.Fatal("Broker Connect error: %v", err)
	}else{
		log.Print("Broker Connect successfully")
	}

	return err
}
