package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

func Execute() (content string, instance string, raw bool) {
	text := &survey.Multiline{
		Message: "enter text to upload",
	}

	if err := survey.AskOne(text, &content, survey.WithValidator(survey.Required)); err == terminal.InterruptErr {
		log.Fatalln("interrupted")
	}

	url := &survey.Select{
		Message: "choose an instance:",
		Options: []string{"https://p.x4.pm", "https://hastebin.cc", "other"},
	}
	survey.AskOne(url, &instance)

	if instance == "other" {
		custom := &survey.Input{
			Message: "enter custom instance:",
		}
		if err := survey.AskOne(custom, &instance, survey.WithValidator(survey.Required)); err == terminal.InterruptErr {
			log.Fatalln("interrupted")
		}
	}

	prompt := &survey.Confirm{
		Message: "return raw file?",
	}
	survey.AskOne(prompt, &raw)

	return
}
