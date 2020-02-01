package main

const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "

func main() {
	// fmt.Println(Hello("World"))
}

//Hello a func return greeting text
func Hello(name, language string) string {
	if name == "" {
		name = "World"
	}

	return getGreetingPrefix(language) + name
}

//Hello a func return greeting text
func getGreetingPrefix(language string) (prefix string) {

	switch language {
	case "Spanish":
		prefix = spanishHelloPrefix
	case "French":
		prefix = frenchHelloPrefix
	default:
		prefix = englishHelloPrefix
	}

	return
}
