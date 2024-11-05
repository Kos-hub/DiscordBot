package commands

import (
	"github.com/bwmarrin/discordgo"
)

func HelpCommand(i *discordgo.InteractionCreate) error {

	embed := &discordgo.MessageEmbed{
		Title:       "Help - Available commands",
		Description: "Commands you can use with this bot:",
		Color:       0x00ffcc,
		Fields:      []*discordgo.MessageEmbedField{},
	}
	//builder := new(strings.Builder)
	for _, c := range b.Commands {
		//builder.WriteString(fmt.Sprintf("/%s - %s\n", c.Name, c.Description))
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "/" + c.Name,
			Value:  c.Description,
			Inline: false,
		})
	}

	_, err := b.Session.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		return err
	}
	return nil
}
