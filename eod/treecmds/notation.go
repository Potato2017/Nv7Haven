package treecmds

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *TreeCmds) NotationCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(db)

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	db.RLock()
	msg, suc := tree.AddElem(el.ID)
	db.RUnlock()
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	txt := tree.String()
	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	if len(txt) <= 2000 {
		id := rsp.Message("Sent notation in DMs!")
		data.SetMsgElem(id, el.ID)
		rsp.DM(txt)
		return
	}
	id := rsp.Message("The notation was too long! Sending it as a file in DMs!")

	data.SetMsgElem(id, el.ID)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Notation for **%s**:", el.Name),
		Files: []*discordgo.File{
			{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}

func (b *TreeCmds) CatNotationCmd(catName string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(db)

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}

	db.RLock()
	for elem := range cat.Elements {
		msg, suc := tree.AddElem(elem)
		if !suc {
			db.RUnlock()
			rsp.ErrorMessage(msg)
			return
		}
	}
	db.RUnlock()

	txt := tree.String()

	if len(txt) <= 2000 {
		rsp.Message("Sent notation in DMs!")

		rsp.DM(txt)
		return
	}
	rsp.Message("The notation was too long! Sending it as a file in DMs!")

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Notation for category **%s**:", cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
