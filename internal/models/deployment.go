package models

import "fmt"

type Deployment int

const (
	CodeDavinci002      Deployment = iota // code-davinci-002
	TextDavinci003                        // text-davinci-003
	Gpt35Turbo0301                        // gpt-3.5-turbo-0301
	Gpt35Turbo                            // gpt-3.5-turbo
	Gpt35Turbo0301Azure                   // gpt-35-turbo-0301
	Gpt4_0314                             // gpt-4-0314
	Gpt4_32k_0314                         // gpt-4-32k-0314
)

func (d Deployment) String() string {
	return [...]string{
		"code-davinci-002",
		"text-davinci-003",
		"gpt-3.5-turbo-0301",
		"gpt-3.5-turbo",
		"gpt-35-turbo-0301",
		"gpt-4-0314",
		"gpt-4-32k-0314",
	}[d]
}

func (d Deployment) MaxTokens() int {
	return [...]int{
		8001, // code-davinci-002
		4097, // text-davinci-003
		4096, // gpt-3.5-turbo-0301
		4096, // gpt-3.5-turbo
		4096, // gpt-35-turbo-0301
		8192, // gpt-4-0314
		8192, // gpt-4-32k-0314
	}[d]
}

func (d Deployment) IsChat() bool {
	return [...]bool{
		false, // code-davinci-002
		false, // text-davinci-003
		true,  // gpt-3.5-turbo-0301
		true,  // gpt-3.5-turbo
		true,  // gpt-35-turbo-0301
		true,  // gpt-4-0314
		true,  // gpt-4-32k-0314
	}[d]
}

func DeploymentFromName(name string) (Deployment, error) {
	switch name {
	case "code-davinci-002":
		return CodeDavinci002, nil
	case "text-davinci-003":
		return TextDavinci003, nil
	case "gpt-3.5-turbo-0301":
		return Gpt35Turbo0301, nil
	case "gpt-3.5-turbo":
		return Gpt35Turbo, nil
	case "gpt-35-turbo-0301":
		return Gpt35Turbo0301Azure, nil
	case "gpt-4-0314":
		return Gpt4_0314, nil
	case "gpt-4-32k-0314":
		return Gpt4_32k_0314, nil
	}

	return -1, fmt.Errorf("The specified deployment does not exist, please choose one of these options: code-davinci-002, text-davinci-003, gpt-3.5-turbo-0301, gpt-3.5-turbo, gpt-35-turbo-0301, gpt-4-0314, gpt-4-32k-0314")
}
