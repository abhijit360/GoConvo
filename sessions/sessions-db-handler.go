package sessions

import (
	"log"
	"database/sql"
	"time"
	"fmt"
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
	subscribers string[] // array of subscriber websocket connections that we write to
	chatMetaData ChatMetaData
	db *sql.DB
}

type Mainnet struct{
	tracker map[uuid]*Session
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
		return nil, fmt.Errorf("database not initialized")
	}

	sql_query := findExpiryBasedOnChatId

    row := s.db.QueryRow(sql_query, s.chatId)
	var foundRow ChatMetaData
	err := row.Scan(&foundRow)
	if err != nil {
		if err == sql.ErrNoRows{
			return "", fmtErrorf("chatId not found")
		}
		return nil, err
	}

	return foundRow, nil
}

func CreateSession(current_time string) (*Session,error){
	current_time, err := time.Parse(current_time)
	if err != nil {
		return "",fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := current_time + (24 * time.Hour) // set expirty to be after 24 hours
	newSession = Session{
		expiry_date: expiry_time
	}
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	createdRow, err := db.Query(createNewSession,expiry_time)
	if err != nil {
		return nil, err
	}
	
	var chatMetaData ChatMetaData
	err := createdRow.Scan(&chatMetaData)
	if err != nil {
		return nil, fmtErrorf("error parsing row data into struct")
	}
	newSession.chatMetaData = chatMetaData

	m.tracker[chat_id] = *newSession
	return *newSession, nil
}

func (s *Session)UpdateSessionExpiryDate(current_time string) (string,error){
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	current_time, err := time.Parse(current_time)
	if err != nil {
		return "",fmt.Errorf("Received time could not be parsed")
	}
	expiry_time := current_time + (24 * time.Hour) // set expirty to be after 24 hours

	updatedRow, err := s.db.Query(updateExistingSession,s.chatMetaData.chat_id,expiry_time)
	
	if err != nil {
		return nil, err
	}
	
	var chatMetaData ChatMetaData
	updatedRow.Scan(&chatMetaData)
	s.chatMetaData = chatMetaData
	return chatMetaData.expiry_date, nil
}
