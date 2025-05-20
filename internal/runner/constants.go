package runner

const (
	LANG_GO         Lang = "go"
	LANG_JAVASCRIPT Lang = "javascript"
	LANG_PYTHON     Lang = "python"
	LANG_RUBY       Lang = "ruby"
	LANG_RUST       Lang = "rust"
	LANG_PHP        Lang = "php"
	LANG_JAVA       Lang = "java"
)

const WORKDIR = "/app"

var LangDefinitions = map[Lang]LangDefinition{
	LANG_GO: {
		image:    "golang:1.21-alpine",
		command:  []string{"go", "run"},
		fileName: "main.go",
		SampleCode: `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`,
	},
	LANG_JAVASCRIPT: {
		image:      "node:20-alpine",
		command:    []string{"node"},
		fileName:   "main.js",
		SampleCode: `console.log("Hello, World!");`,
	},
	LANG_PYTHON: {
		image:      "python:3.11-alpine",
		command:    []string{"python3"},
		fileName:   "main.py",
		SampleCode: `print("Hello, World!")`,
	},
	LANG_RUBY: {
		image:      "ruby:3.3-alpine",
		command:    []string{"ruby"},
		fileName:   "main.rb",
		SampleCode: "puts 'Hello, World!'",
	},
	LANG_RUST: {
		image:    "rust:1.77-alpine",
		command:  []string{"sh", "-c", "rustc /app/main.rs -o /app/main && /app/main"},
		fileName: "main.rs",
		SampleCode: `fn main() {
    println!("Hello, World!");
}`,
	},
	LANG_PHP: {
		image:    "php:8.3-cli-alpine",
		command:  []string{"php"},
		fileName: "main.php",
		SampleCode: `<?php
echo "Hello, World!";`,
	},
	LANG_JAVA: {
		image:    "openjdk:21-slim",
		command:  []string{"sh", "-c", "javac /app/Main.java && java -cp /app Main"},
		fileName: "Main.java",
		SampleCode: `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`,
	},
}

func GetLangs() []string {
	var languages []string
	for k := range LangDefinitions {
		languages = append(languages, string(k))
	}
	return languages
}
