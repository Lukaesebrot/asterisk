package commands

import (
	"reflect"

	"github.com/Lukaesebrot/asterisk/utils"
	"github.com/Lukaesebrot/dgc"
	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
)

// Debug handles the debug command
func Debug(ctx *dgc.Ctx) {
	// Validate the arguments
	codeblock := ctx.Arguments.AsCodeblock()
	if codeblock == nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, utils.GenerateInvalidUsageEmbed(ctx.Command.Usage))
		return
	}

	// Create the interpreter
	interpreter := interp.New(interp.Options{})

	// Inject the custom variables
	custom := make(map[string]map[string]reflect.Value)
	custom["main"] = map[string]reflect.Value{
		"ctx": reflect.ValueOf(ctx),
	}
	interpreter.Use(stdlib.Symbols)
	interpreter.Use(custom)
	_, err := interpreter.Eval("import (\n. \"main\"\n\"fmt\"\n)")
	if err != nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, utils.GenerateErrorEmbed(err.Error()))
		return
	}

	// Evaluate the given string and output the result
	_, err = interpreter.Eval(codeblock.Content)
	if err != nil {
		ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, utils.GenerateErrorEmbed(err.Error()))
		return
	}
	ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID, utils.GenerateSuccessEmbed("Evaluation succeeded."))
}