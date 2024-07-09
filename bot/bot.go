package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BotToken string
var Done chan bool

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func Run() {
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	discord.AddHandler(newMessage)

	discord.Open()
	defer discord.Close()

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func waitForStartTime(discord *discordgo.Session, channelID string, frequency time.Duration, startTimestamp time.Time, message string) {
	diff := time.Until(startTimestamp)
	for {
		select {
		case <-Done:
			fmt.Println("Done!")
			return
		case <-time.After(diff):
			startAlert(discord, channelID, frequency, startTimestamp, message)
			return
		}
	}
}

func startAlert(discord *discordgo.Session, channelID string, frequency time.Duration, startTimestamp time.Time, message string) {
	role := fmt.Sprintf("<@&%s>", "1260019771882082334")
	alert := fmt.Sprintf("Attention %s, %s", role, message)
	diff := time.Until(startTimestamp)
	time.Sleep(diff)
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()
	for {
		select {
		case <-Done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			discord.ChannelMessageSend(channelID, alert)
		}
	}
}

func handleStartAlertRequest(discord *discordgo.Session, message *discordgo.MessageCreate) {
	var re = regexp.MustCompile(`(?m)^\!setAlert:\d+["h","m","s"]:\d+:\w+`)
	matched := re.Match([]byte(message.Content))
	if matched {
		inputs := strings.Split(message.Content, ":")
		frequency, err := time.ParseDuration(inputs[1])
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Invalid frequency format. Please use units h,m,s. For example, 8h")
		}

		timestamp, err := strconv.ParseInt(inputs[2], 10, 64)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Invalid timestamp. Please use a unix timestamp for your alert start time.")
		}
		fmt.Println("setting alert")
		waitForStartTime(discord, message.ChannelID, frequency, time.Unix(timestamp, 0), inputs[3])
	} else {
		discord.ChannelMessageSend(message.ChannelID, "Invalid format. Please set alert with format ```!setAlert:frequency:startTimestamp:message```\n For example, ```!setAlert:8h:1260006791375224943:hello```")
	}
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "!setAlert"):
		handleStartAlertRequest(discord, message)
	case strings.Contains(message.Content, "!stopAlert"):
		Done <- true
		discord.ChannelMessageSend(message.ChannelID, "Done")
	case strings.Contains(message.Content, "!alertHelp"):
		discord.ChannelMessageSend(message.ChannelID, "You can set a recurring alert with ```!setAlert:frequency:startTimestamp:message```\n For example, ```!setAlert:8h:1260006791375224943:hello```")
		discord.ChannelMessageSend(message.ChannelID, "To stop all alerts, use ```!stopAlert```")
	}

}
