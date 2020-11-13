package signalk

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Value ...
type Value struct {
	Path  string  `json:"path"`
	Valor float64 `json:"value"`
}

// Update ...
type Update struct {
	source    string
	Values    []Value `json:"values"`
	timeStamp string
}

// WsSignalk ..
type WsSignalk struct {
	Updates []Update `json:"updates"`
	Context string   `json:"context"`
}

func Connect_signalk() (*websocket.Conn, error) {
	var ErrNotFound = errors.New("not found")
	u := url.URL{Scheme: "ws",
		Host: "localhost:3000",
		Path: "/signalk/v1/stream",
	}
	log.Info(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// handle error
		log.Error(err)
		return nil, ErrNotFound
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		// handle error
		log.Error(err)
		return nil, ErrNotFound
	}
	var v interface{}
	json.Unmarshal(message, &v)
	data := v.(map[string]interface{})
	for k, v := range data {
		switch v := v.(type) {
		case string:
			log.Info(k, v, "(string)")
		case float64:
			log.Info(k, v, "(float64)")
		case []interface{}:
			log.Info(k, "(array):")
			for i, u := range v {
				log.Info("    ", i, u)
			}
		default:
			log.Info(k, v, "(unknown)")
		}
	}
	return c, nil
}

func Subscribe(jStr []byte, subsc string) (float64, error) {
	var ErrNotFound = errors.New("not found")
	var WsK WsSignalk
	err := json.Unmarshal(jStr, &WsK)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(jStr))
	for _, s := range WsK.Updates {
		//fmt.Println(i, s)
		for _, s1 := range s.Values {
			//fmt.Println("subscribed", s1)
			log.Info("subscribed path: ", s1.Path)
			log.Info("subscribed valor: ", s1.Valor)
			if s1.Path == subsc {
				return s1.Valor, nil
			}
		}
	}
	return 0, ErrNotFound
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	f_log, _ := os.Create("signalk.log")
	//check(err)
	log.SetOutput(f_log)

	// Only log the warning severity or above.
	log.SetLevel(log.ErrorLevel)
}
