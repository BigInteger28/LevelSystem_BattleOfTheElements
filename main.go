package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	Plaats     int
	Naam       string
	Level      int
	Rating     int
	ColorName  string
	Tier       int
	Commentaar string
	Foreground string
	Color      string
}

var leagues = []struct {
	Name       string
	Background string
	Foreground string
}{
	{"White", "#FFFFFF", "black"},
	{"Brown", "#8B4513", "white"},
	{"Grey", "#C0C0C0", "black"},
	{"Taupe", "#B9A281", "white"},
	{"Olive", "#808000", "white"},
	{"Black", "#000000", "white"},
	{"Yellow", "#FFFF00", "black"},
	{"Green", "#008000", "white"},
	{"Lime", "#00FF00", "black"},
	{"Mint", "#98FF98", "black"},
	{"Teal Green", "#00827F", "white"},
	{"Blue", "#0000FF", "white"},
	{"Dark Blue", "#00008B", "white"},
	{"Red", "#FF0000", "white"},
	{"Crimson", "#DC143C", "white"},
	{"Dark Red", "#C30000", "white"},
	{"Ochre", "#FFC619", "black"},
	{"Orange", "#FF8C00", "black"},
	{"Bright Lavender", "#BF94E4", "black"},
	{"Purple", "#800080", "white"},
	{"Turquoise", "#40e0d0", "black"},
	{"Cyan", "#00FFFF", "black"},
	{"Pink", "#FFB3DE", "black"},
	{"Magenta", "#FF64B4", "white"},
	{"Salmon", "#FA8072", "black"},
}

func getColorAndForeground(level int) (string, string) {
	tierIndex := (level - 1) % 25
	if tierIndex >= len(leagues) {
		tierIndex = len(leagues) - 1
	}
	return leagues[tierIndex].Name, leagues[tierIndex].Foreground
}

func getTier(level int) int {
	return ((level - 1) / 25) + 1
}

func getColorBackground(level int) string {
	tierIndex := (level - 1) % 25
	if tierIndex >= len(leagues) {
		tierIndex = len(leagues) - 1
	}
	return leagues[tierIndex].Background
}

func main() {
	// Open the output file
	file, err := os.Open("bote.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	var entries []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line by tabs
		parts := strings.Split(line, "   ")
		if len(parts) < 3 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}

		level, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error parsing level:", err, "in line:", line)
			continue
		}

		rating, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Error parsing rating:", err, "in line:", line)
			continue
		}

		comment := ""
		if len(parts) == 4 {
			comment = parts[3]
		}
		colorName, foreground := getColorAndForeground(level)
		colorBackground := getColorBackground(level)

		entries = append(entries, Entry{
			Naam:       parts[0],
			Level:      level,
			Rating:     rating,
			ColorName:  colorName,
			Tier:       getTier(level),
			Commentaar: comment,
			Foreground: foreground,
			Color:      colorBackground,
		})
	}

	// Sort entries by Level, then by elo, with names starting with '---' at the bottom of their level
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Level == entries[j].Level {
			if strings.HasPrefix(entries[i].Naam, "---") && !strings.HasPrefix(entries[j].Naam, "---") {
				return false
			}
			if !strings.HasPrefix(entries[i].Naam, "---") && strings.HasPrefix(entries[j].Naam, "---") {
				return true
			}
			return entries[i].Rating > entries[j].Rating
		}
		return entries[i].Level > entries[j].Level
	})

	// Assign correct place values
	var pos int = 1
	for i := range entries {
		if !strings.HasPrefix(entries[i].Naam, "---") {
			entries[i].Plaats = pos
			pos++
		}
	}

	// Generate HTML
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))
	outputFile, err := os.Create("index.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, entries)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("HTML report generated successfully.")
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Level Report</title>
	<style>
		table { width: 100%; border-collapse: collapse; }
		th, td { padding: 8px; text-align: left; border: 1px solid #ddd; text-align: center; }
		th { background-color: #f2f2f2; }
	</style>
</head>
<body>
	<h1>Level Report</h1>
	<table>
		<tr>
			<th>Plaats</th>
			<th>Naam</th>
			<th>Level</th>
			<th>Color</th>
			<th>Tier</th>
			<th>Rating</th>
			<th>Commentaar</th>
		</tr>
		{{range .}}
		<tr style="background-color: {{.Color}}; color: {{.Foreground}}">
			<td>{{if .Plaats}}{{.Plaats}}{{end}}</td>
			<td>{{.Naam}}</td>
			<td>{{.Level}}</td>
			<td>{{.ColorName}}</td>
			<td>{{.Tier}}</td>
			<td>{{.Rating}}</td>
			<td>{{.Commentaar}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`
