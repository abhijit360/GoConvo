package sessions

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

const database = "WebsocketManager.db"
const findExpiryBasedOnChatId = "Select expiry_date from chats where chat_id == ?"
const createNewSession = "INSERT INTO chats (chat_id, expiry_date) VALUES (hex(randomblob(16)), ?)"
const updateExistingSession = "UPDATE chats SET expiry_date = ? where chat_id == ?"

type ChatMetaData struct {
	Chat_id     string
	Expiry_date string
}

type Session struct {
	sync.RWMutex
	Subscribers  map[*websocket.Conn]bool // array of subscriber websocket connections that we write to
	Broadcast    chan []byte
	ChatMetaData ChatMetaData
	done         chan struct{}
	DB           *sql.DB
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS chats (
		chat_id TEXT PRIMARY KEY,
		expiry_date TEXT
	);`)
	if err != nil {
		log.Fatalln("failed to create table", err)
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
		Subscribers: make(map[*websocket.Conn]bool),
		Broadcast:   make(chan []byte),
		done:        make(chan struct{}),
		DB:          db,
	}
	parsedTime, err := time.Parse(time.RFC3339, current_time)
	if err != nil {
		return &newSession, fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := parsedTime.Add(24 * time.Hour) // fixed time addition

	// fmt.Printf("parsedTime: %v | ExpiryTime: %v | DB: %v",parsedTime, expiry_time, newSession.DB)
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var chatMetaData ChatMetaData 
	err = db.QueryRow(createNewSession, expiry_time).Scan(&chatMetaData.Chat_id, &chatMetaData.Expiry_date)
	if err != nil {
		fmt.Printf("error executing the query")
		return nil, err
	}
	newSession.ChatMetaData = chatMetaData

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
	s.Lock()
	delete(s.Subscribers, c)
	isEmpty := len(s.Subscribers) == 0
	s.Unlock()

	if isEmpty {
		s.cleanup()
	}
}

func (s *Session) AddSession(c *websocket.Conn) {
	s.Lock()
	s.Subscribers[c] = true
	s.Unlock()
}

func GetSession(chat_id string) (*Session, bool) {
	foundSession := Session{}
	if session, ok := sessionManager.tracker[chat_id]; !ok {
		return &foundSession, ok
	} else {
		return session, ok
	}
}

func (s *Session) HandleBroadcast() {
	for {
		select {
		case message := <-s.Broadcast:
			s.RLock()
			for conn, active := range s.Subscribers {
				if !active {
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					s.RemoveSession(conn)
				}
			}
			s.RUnlock()
		case <-s.done:
			return
		}
	}
}

func (s *Session) GetChatMetaData() ChatMetaData {
	return s.ChatMetaData
}

func (s *Session) cleanup() {
	// Signal goroutines to stop
	close(s.done)

	// Close broadcast channel
	close(s.Broadcast)

	// Clean up all connections
	s.Lock()
	for conn := range s.Subscribers {
		conn.Close()
	}
	s.Subscribers = nil
	s.Unlock()

	// Remove from session manager
	delete(sessionManager.tracker, s.ChatMetaData.Chat_id)

}
