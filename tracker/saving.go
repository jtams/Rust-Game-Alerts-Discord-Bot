package tracker

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type MessengerSaveData struct {
	MessageID string `json:"messageID"`
	ChannelID string `json:"channelID"`
}

type SaveFileData struct {
	Tracker       *PlayerTracker     `json:"tracker"`
	MessengerData *MessengerSaveData `json:"messengerData"`
}

func SaveTrackerData(filename string, tracker *PlayerTracker, messenger *Messenger) error {
	msgData := MessengerSaveData{
		MessageID: messenger.Message.ID,
		ChannelID: messenger.Message.ChannelID,
	}

	saveData := SaveFileData{
		Tracker:       tracker,
		MessengerData: &msgData,
	}

	// tracker.Channel = make(chan time.Time)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(&saveData)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func LoadTrackerData(filename string, session *discordgo.Session) (error, *Messenger, *PlayerTracker) {
	saveData := &SaveFileData{}

	rawSaveData, err := os.ReadFile(filename)
	if err != nil {
		return err, nil, nil
	}

	if err := json.Unmarshal(rawSaveData, saveData); err != nil {
		return err, nil, nil
	}

	if saveData.Tracker == nil || saveData.MessengerData == nil {
		return errors.New("failed to loads tracker or messenger data from file, file may be corrupted"), nil, nil
	}

	messenger := NewMessageUpdater(session)
	messenger.Message, err = session.ChannelMessage(saveData.MessengerData.ChannelID, saveData.MessengerData.MessageID)
	if messenger.Message != nil && saveData.Tracker.Running {
		content := "Reestablishing connection..."
		msgEdit := &discordgo.MessageEdit{
			Channel: saveData.MessengerData.ChannelID,
			ID:      saveData.MessengerData.MessageID,
			Content: &content,
		}

		messenger.Session.ChannelMessageEditComplex(msgEdit)
	}

	if messenger.Message == nil && saveData.Tracker.Running {
		log.Println("Failed to load message, creating new message")
		channel, _ := session.Channel(saveData.MessengerData.ChannelID)
		if channel != nil {
			newMessage, err := session.ChannelMessageSend(saveData.MessengerData.ChannelID, "Loading...")
			if newMessage != nil && err == nil {
				messenger.Message = newMessage
			} else {
				saveData.Tracker.Running = false
				log.Println("Failed to send message")
			}
		} else {
			saveData.Tracker.Running = false
			log.Println("Failed to resolve channel")
		}
	}

	return nil, messenger, saveData.Tracker
}
