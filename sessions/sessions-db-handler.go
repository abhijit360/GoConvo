package sessions

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

const database = "WebsocketManager.db"
const findExpiryBasedOnChatId = "Select expiry_date from chats where chat_id == ?"
const createNewSession = "INSERT INTO chats (expiry_date) VALUES (?)"
const updateExistingSession = "UPDATE chats SET expiry_date = ? where chat_id == ?"

type ChatMetaData struct {
	Chat_id     string
	Expiry_date string
}

type Session struct {
	Subscribers  map[*websocket.Conn]bool // array of subscriber websocket connections that we write to
	Broadcast    chan []byte
	ChatMetaData ChatMetaData
	done         chan struct{}
	DB           *sql.DB
}

func (s *Session) HandleBroadcast() {
	panic("unimplemented")
}

type Mainnet struct {
	tracker map[string]*Session
}

var db *sql.DB
var sessionManager Mainnet

func init() {
	var err error
	db, err = sql.Open("sqlite3", database)
	if err != nil {
		log.Fatalf("failed to connect to database %v", err)
	}
}

func (s *Session) GetSession() (string, error) {
	if s.DB == nil {
		return "", fmt.Errorf("database not initialized")
	}

	sql_query := findExpiryBasedOnChatId

	row := s.DB.QueryRow(sql_query, s.ChatMetaData.Chat_id)
	var foundRow ChatMetaData
	err := row.Scan(&foundRow.Chat_id, &foundRow.Expiry_date)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("chatId not found")
		}
		return "", err
	}

	return foundRow.Expiry_date, nil
}

func CreateSession(current_time string) (*Session, error) {
	newSession := Session{
		Subscribers: make(map[*websocket.Conn]bool, 0),
	}
	parsedTime, err := time.Parse(time.RFC3339, current_time)
	if err != nil {
		return &newSession, fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := parsedTime.Add(24 * time.Hour) // fixed time addition

	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	createdRow, err := db.Query(createNewSession, expiry_time)
	if err != nil {
		return nil, err
	}

	var chatMetaData ChatMetaData
	err = createdRow.Scan(&chatMetaData.Chat_id, &chatMetaData.Expiry_date)
	if err != nil {
		return nil, fmt.Errorf("error parsing row data into struct")
	}
	newSession.ChatMetaData = chatMetaData
	newSession.DB = db
	newSession.Broadcast = make(chan []byte)
	newSession.done = make(chan struct{})

	sessionManager.tracker[chatMetaData.Chat_id] = &newSession
	return &newSession, nil
}

func (s *Session) UpdateSessionExpiryDate(current_time string) (string, error) {
	if s.DB == nil {
		return "", fmt.Errorf("database not initialized")
	}
	parsedTime, err := time.Parse(time.RFC3339, current_time)
	if err != nil {
		return "", fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := parsedTime.Add(24 * time.Hour) // fixed time addition

	updatedRow, err := s.DB.Query(updateExistingSession, s.ChatMetaData.Chat_id, expiry_time)

	if err != nil {
		return "", err
	}

	var chatMetaData ChatMetaData
	updatedRow.Scan(&chatMetaData)
	s.ChatMetaData = chatMetaData
	return chatMetaData.Expiry_date, nil
}

func (s *Session) RemoveSession(c *websocket.Conn) {
	delete(s.Subscribers, c)
	if len(s.Subscribers) == 0 {
		// reached a case where we have no more subscribers listening to this socket
		close(s.done)
		close(s.Broadcast)
		delete(sessionManager.tracker, s.ChatMetaData.Chat_id)
	}
}

func (s *Session) AddSession(c *websocket.Conn) {
	s.Subscribers[c] = true
}

func GetSession(chat_id string) (*Session, bool) {
	foundSession := Session{}
	if session, ok := sessionManager.tracker[chat_id]; !ok {
		return &foundSession, ok
	} else {
		return session, ok
	}
}

func (s *Session) handleBroadcast() {
	for {
		select {
		case message := <-s.Broadcast:
			for conn, active := range s.Subscribers {
				if !active {
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					s.RemoveSession(conn)
				}
			}
		case <-s.done:
			return
		}
	}
}

func (s *Session) GetChatMetaData() ChatMetaData {
	return s.ChatMetaData
}
