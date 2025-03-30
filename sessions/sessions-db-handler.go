package sessions

import (
	"context"
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
	chat_id string
	expiry_date string
}

type Session struct{
	subscribers map[*websocket.Conn]bool // array of subscriber websocket connections that we write to
	broadcast chan []byte
	chatMetaData ChatMetaData
	db *sql.DB
}

type Mainnet struct{
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

func (s *Session)GetSession() (string,error){
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	sql_query := findExpiryBasedOnChatId

    row := s.db.QueryRow(sql_query, s.chatMetaData.chat_id)
	var foundRow ChatMetaData
	err := row.Scan(&foundRow.chat_id,&foundRow.expiry_date)
	if err != nil {
		if err == sql.ErrNoRows{
			return "", fmt.Errorf("chatId not found")
		}
		return "", err
	}

	return foundRow.expiry_date, nil
}

func CreateSession(current_time string) (*Session,error){
	newSession := Session{
		subscribers: make(map[*websocket.Conn]bool,0),
	}
	parsedTime, err := time.Parse(time.RFC3339, current_time)
	if err != nil {
		return &newSession, fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := parsedTime.Add(24 * time.Hour) // fixed time addition
	
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	createdRow, err := db.Query(createNewSession,expiry_time)
	if err != nil {
		return nil, err
	}
	
	var chatMetaData ChatMetaData
	err = createdRow.Scan(&chatMetaData.chat_id, &chatMetaData.expiry_date)
	if err != nil {
		return nil, fmt.Errorf("error parsing row data into struct")
	}
	newSession.chatMetaData = chatMetaData
	newSession.db = db
	newSession.broadcast = make(chan []byte)

	sessionManager.tracker[chatMetaData.chat_id] = &newSession
	return &newSession, nil
}

func (s *Session)UpdateSessionExpiryDate(current_time string) (string,error){
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}
	parsedTime, err := time.Parse(time.RFC3339, current_time)
	if err != nil {
		return "", fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := parsedTime.Add(24 * time.Hour) // fixed time addition

	updatedRow, err := s.db.Query(updateExistingSession,s.chatMetaData.chat_id,expiry_time)
	
	if err != nil {
		return "", err
	}
	
	var chatMetaData ChatMetaData
	updatedRow.Scan(&chatMetaData)
	s.chatMetaData = chatMetaData
	return chatMetaData.expiry_date, nil
}

func (s *Session)RemoveSession(c *websocket.Conn){
	delete(s.subscribers,c)
}

func (s *Session)AddSession(c *websocket.Conn){
	s.subscribers[c] = true
}


func GetSession(chat_id string) (*Session, bool){
	foundSession := Session{}
	if session, ok := sessionManager.tracker[chat_id]; !ok {
		return &foundSession, ok
	}else {
		return session, ok
	}
}

func (s *Session) handleBroadcast(){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for message := range s.broadcast {
		for conn, active := range s.subscribers{
			if !active{
				continue
			}
			if err := conn.Write(ctx,websocket.MessageText,message); err != nil {
				s.RemoveSession(conn)
			}
		}
	}
}