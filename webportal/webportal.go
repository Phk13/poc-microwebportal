package webportal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/phk13/poc-micro/databaselayer"
	"github.com/phk13/poc-micro/webportal/api"
	"github.com/phk13/poc-micro/webportal/template"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type DataFeedMessage struct {
	Heartrate     int
	Bloodpressure int
}

const (
	INPUTNAME     = "inputname"
	SIGNINSESSION = "signinsession"
	USERNAME      = "username"
)

// RunWebPortal starts running the web portal on address addr
func RunWebPortal(dbtype uint8, addr, dbconnection, frontend, secret string) error {
	var cookieStore = sessions.NewCookieStore([]byte(secret))
	rand.Seed(time.Now().UTC().UnixNano())
	r := mux.NewRouter()
	db, err := databaselayer.GetDatabaseHandler(dbtype, dbconnection)
	if err != nil {
		return err
	}
	api.RunApiOnRouter(r, db)

	r.Path("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		session, err := cookieStore.Get(req, SIGNINSESSION)
		if err != nil {
			template.HandleSignIn(w)
			return
		}
		log.Println(session.Values)
		val, ok := session.Values[USERNAME]
		if !ok {
			template.HandleSignIn(w)
			return
		}
		name, ok := val.(string)
		if !ok {
			template.HandleSignIn(w)
			return
		}
		template.Homepage("Dino portal", fmt.Sprintf("Hello %s, welcome to the Dino portal, where you can find metrics and information", name), w)
	})

	r.PathPrefix("/signin/").Methods("POST").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			return
		}
		namelist := req.Form[INPUTNAME]
		session, err := cookieStore.Get(req, SIGNINSESSION)
		if err != nil || len(namelist) == 0 {
			return
		}
		session.Values[USERNAME] = namelist[0]
		session.Save(req, w)
		template.Homepage("Dino portal", fmt.Sprintf("Hello %s, welcome to the Dino portal, where you can find metrics and information", namelist[0]), w)
	})

	r.PathPrefix("/metrics/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		animals, err := db.GetAvailableDinos()
		if err != nil {
			return
		}
		template.HandleMetrics(animals, w)
	})

	r.PathPrefix("/info/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		animals, err := db.GetAvailableDinos()
		if err != nil {
			return
		}
		template.HandleInfo(animals, w)
	})

	fileserver := http.FileServer(http.Dir(frontend))
	r.Path("/datafeed").HandlerFunc(dataFeedHandler)

	r.PathPrefix("/").Handler(fileserver)
	return http.ListenAndServe(addr, r)
}

func dataFeedHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Could not establish websocket connection, error", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Could not read message from websocket, error", err, messageType)
			return
		}

		go func(dino string) {
			for {
				time.Sleep(1 * time.Second)
				msg := getDinoData(dino)
				databytes, err := json.Marshal(msg)
				if err != nil {
					log.Println("Could not convert data to JSON, error", err)
					return
				}

				err = conn.WriteMessage(websocket.TextMessage, databytes)
				if err == websocket.ErrCloseSent {
					return
				}
				if err != nil {
					log.Println("Could not write message to websocket, error", err)
					return
				}
			}
		}(string(p))
	}
}

func getDinoData(dinoName string) *DataFeedMessage {
	return &DataFeedMessage{rand.Intn(50) + 120, rand.Intn(300) + 1000}
}
