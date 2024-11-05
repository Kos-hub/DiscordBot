package commands

import (
	"github.com/bwmarrin/discordgo"
)

func HelpCommand(i *discordgo.InteractionCreate, args []string) error {

	embed := &discordgo.MessageEmbed{
		Title:       "Help - Comandi disponibili",
		Description: "Comandi chi poi utilizzari cu 'ssu bot",
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
