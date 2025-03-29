package sessions

import (
	"fmt"
)

type message struct {
	id string
	chat_id string
	message string
}

type chatLog struct{
	messages []message
}


const getChatlog = "select * from messages where chat_id = ?"
const updateChatLog = "Insert INTO messages values (?,?)"


func getChatLog(chat_id string) (chatLog,error){
	messages := chatLog{
		messages: make([]message, 0),
	}
	if db == nil {
		return messages, fmt.Errorf("error connecting to database")
	}
	rows, err := db.Query(getChatlog,chat_id)
	if err != nil {
		return messages, fmt.Errorf("eror carrying out query to database")
	}


	for rows.Next() {
		var m message
		err := rows.Scan(&m.id, &m.chat_id, &m.message)
		if err != nil {
			return messages, fmt.Errorf("returned row does not match expected structure %v",m)
		}
		messages.messages = append(messages.messages, m)
	}
	defer rows.Close()
	
	return messages, nil
}

func updateLog(chat_id string, messages []string) error {
	for _, m := range messages {
		_, err := db.Exec(updateChatLog,chat_id,m)
		if err != nil {
			return fmt.Errorf("error updating log: %v", err)
		}		
	}
}