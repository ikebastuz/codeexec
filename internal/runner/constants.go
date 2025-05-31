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
const (
	ERR_TIMEOUT_EXCEEDED = "timeout exceeded"
)

var LangDefinitions = map[Lang]LangDefinition{
	LANG_GO: {
		image:          "golang:1.21-alpine",
		buildCommand:   []string{"go", "build", "-o", "main", "main.go"},
		execCommand:    []string{"./main"},
		sourceFileName: "main.go",
		SampleCode: `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`,
		nonDetermenisticKeywords: []string{"rand.", "time.Now"},
	},
	LANG_JAVASCRIPT: {
		image:                    "node:20-alpine",
		execCommand:              []string{"node"},
		sourceFileName:           "main.js",
		SampleCode:               `console.log("Hello, World!");`,
		nonDetermenisticKeywords: []string{"Math.random()", "Date.now()", "new Date()"},
	},
	LANG_PYTHON: {
		image:          "python:3.11-alpine",
		execCommand:    []string{"python3"},
		sourceFileName: "main.py",
		SampleCode:     `print("Hello, World!")`,
		nonDetermenisticKeywords: []string{"random.", "time.time()", "datetime.now"},
	},
	LANG_RUBY: {
		image:          "ruby:3.3-alpine",
		execCommand:    []string{"ruby"},
		sourceFileName: "main.rb",
		SampleCode:     "puts 'Hello, World!'",
		nonDetermenisticKeywords: []string{"rand", "Random.", "Time.now"},
	},
	LANG_RUST: {
		image:          "rust:1.77-alpine",
		buildCommand:   []string{"rustc", "main.rs", "-o", "main"},
		execCommand:    []string{"./main"},
		sourceFileName: "main.rs",
		SampleCode: `fn main() {
    println!("Hello, World!");
}`,
		nonDetermenisticKeywords: []string{"rand::", "Instant::now", "SystemTime::now"},
	},
	LANG_PHP: {
		image:          "php:8.3-cli-alpine",
		execCommand:    []string{"php"},
		sourceFileName: "main.php",
		SampleCode: `<?php
echo "Hello, World!";`,
		nonDetermenisticKeywords: []string{"rand()", "random_int", "time()", "date(", "microtime"},
	},
	LANG_JAVA: {
		image:          "openjdk:21-slim",
		buildCommand:   []string{"javac", "Main.java"},
		execCommand:    []string{"java", "-cp", "/app", "Main"},
		sourceFileName: "Main.java",
		SampleCode: `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`,
		nonDetermenisticKeywords: []string{"Math.random", "Random.", "System.currentTimeMillis", "LocalDateTime.now"},
	},
}

func GetLangs() []string {
	var languages []string
	for k := range LangDefinitions {
		languages = append(languages, string(k))
	}
	return languages
}
