package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No GITHUB token present")
	}

	key := os.Getenv("WAKA_TIME_API_KEY")
	if key == "" {
		log.Fatal("Unauthorized: No WakaTime API key present")
	}

	owner := os.Getenv("INPUT_OWNER")
	if owner == "" {
		log.Fatal("option INPUT_OWNER missing")
	}

	repository := os.Getenv("INPUT_REPOSITORY")
	if repository == "" {
		log.Fatal("option INPUT_REPOSITORY missing")
	}

	commitMessage := os.Getenv("INPUT_COMMIT_MESSAGE")
	if commitMessage == "" {
		log.Fatal("option INPUT_COMMIT_MESSAGE missing")
	}

	showTitle, err := strconv.ParseBool(os.Getenv("INPUT_SHOW_TITLE"))
	if err != nil {
		log.Fatal("option INPUT_SHOW_TITLE missing")
	}

	tableHeight, err := strconv.Atoi(os.Getenv("INPUT_TABLE_HEIGHT"))
	if err != nil {
		log.Fatal("option INPUT_TABLE_HEIGHT missing")
	}

	cellPadding, err := strconv.Atoi(os.Getenv("INPUT_CELL_PADDING"))
	if err != nil {
		log.Fatal("option INPUT_CELL_PADDING missing")
	}

	blocks := os.Getenv("INPUT_BLOCKS")
	if blocks == "" {
		log.Fatal("option INPUT_BLOCKS missing")
	}
	blockRunes := []rune(blocks)

	blockLength, err := strconv.Atoi(os.Getenv("INPUT_BLOCK_LENGTH"))
	if err != nil {
		log.Fatal("option INPUT_BLOCK_LENGTH missing")
	}

	startSection := os.Getenv("INPUT_START_SECTION")
	if startSection == "" {
		log.Fatal("option INPUT_START_SECTION missing")
	}

	endSection := os.Getenv("INPUT_END_SECTION")
	if endSection == "" {
		log.Fatal("option INPUT_END_SECTION missing")
	}

	client := newGithubClient(token)

	oldReadme := getReadme(client, owner)
	statistics := getWakaTimeStats(key)

	wakaData := getRepresentData(
		statistics.Data.Languages,
		showTitle, tableHeight, cellPadding,
		blockLength, blockRunes,
		startSection, endSection,
	)
	newReadme := replaceContent(oldReadme, wakaData, startSection, endSection)

	fmt.Println(newReadme)
	if oldReadme != newReadme {
		updateReadme(client, owner, repository, newReadme, commitMessage)
	}
}
