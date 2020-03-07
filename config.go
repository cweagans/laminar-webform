package main

// Config is the global config object for laminar-webform.
type Config struct {
	General General
	Forms   map[string]Form
}

// General contains high level application information.
type General struct {
	Title      string
	LaminarURL string `toml:"laminar_url"`
	Debug      bool
}

// Form defines a form to be presented to a user.
type Form struct {
	Title       string
	Description string
	Job         string
	Fields      []FormField
}

// FormField defins an individual element of a form.
type FormField struct {
	Title       string
	Description string
	Name        string
	Type        string   // Supported types: select, text, or longtext.
	Options     []string // Required for type=select
	Filter      string   // Required for type=text. Should contain a regex that matches only on valid contents (usually including ^ at the beginning and $ at the end to make sure you're matching the entire line of text)
}
