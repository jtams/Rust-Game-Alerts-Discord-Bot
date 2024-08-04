package tracker

import (
	"encoding/json"
	"errors"
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

// Save the tracker data and messenger info to a file
func SaveTrackerData(filename string, tracker *PlayerTracker, messenger *Messenger) error {
	// Save the message ID and channel ID
	// This is so we can reestablish the connection to the original message
	msgData := MessengerSaveData{
		MessageID: messenger.Message.ID,
		ChannelID: messenger.Message.ChannelID,
	}

	saveData := SaveFileData{
		Tracker:       tracker,
		MessengerData: &msgData,
	}

	// Overwrite or create the file
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

// Load the tracker data and messenger info from a file
func LoadTrackerData(filename string, session *discordgo.Session) (error, *Messenger, *PlayerTracker) {
	saveData := &SaveFileData{}

	// Read file data
	rawSaveData, err := os.ReadFile(filename)
	if err != nil {
		return err, nil, nil
	}

	// Unmarshal the data
	if err := json.Unmarshal(rawSaveData, saveData); err != nil {
		return err, nil, nil
	}

	// If the tracker or messenger data wasn't able to be parsed.
	if saveData.Tracker == nil || saveData.MessengerData == nil {
		return errors.New("failed to loads tracker or messenger data from file, file may be corrupted"), nil, nil
	}

	// If the tracker was running when it was shutdown, prepare to resume.
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

	// Couldn't load the original messsage, create a new one
	if messenger.Message == nil && saveData.Tracker.Running {
		logger.Error("Failed to load message, creating new message")
		channel, _ := session.Channel(saveData.MessengerData.ChannelID)
		if channel != nil {
			newMessage, err := session.ChannelMessageSend(saveData.MessengerData.ChannelID, "Loading...")
			if newMessage != nil && err == nil {
				messenger.Message = newMessage
			} else {
				saveData.Tracker.Running = false
				logger.Error("Failed to send message", "error", err)
			}
		} else {
			saveData.Tracker.Running = false
			logger.Error("Failed to resolve channel")
		}
	}

	return nil, messenger, saveData.Tracker
}
