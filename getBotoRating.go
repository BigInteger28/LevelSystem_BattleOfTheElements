package main

import (
	"fmt"
	"strconv"
)

func getUserInput(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scanln(&input)
	return input
}

func getLevel(rating uint64) uint64 {
	var ratingEachLevel uint64 = 75
	var ratingLevel2 uint64 = 500
	if rating < ratingLevel2 {
		return 1
	} else {
		return ((rating - ratingLevel2) / ratingEachLevel) + 2
	}
}

func getRating(level uint64) uint64 {
  if level <= 1 {
    return 0
  }
  return 500+(75*(level-2))
}

func main() {
  var level, rating, totlevel uint64
  var colors = []string{
		"White", "Grey", "Yellow", "Ochre Yellow", "Salmon", "Orange", "Lime", "Mint", "Green", "Teal Green", "Cyan", "Blue", "Dark_Blue", "Pink", "Magenta", "Bright Lavender", "Purple", "Indigo", "Olive", "Taupe", "Brown", "Red", "Crimson", "Dark_Red", "Black",
	}
  for {
    input := getUserInput("\nLevel/Rating/Tabel tot level (..l / ..r / ..t): ")
    if input[len(input)-1:] == "l" {
	level, _ = strconv.ParseUint(input[0:len(input)-1], 10, 64)
      	rating = getRating(level)
      	fmt.Println("Level", input[:len(input)-1], " Tier", ((level-1)/25)+1, colors[(level-1)%25], " --> ", rating, "Rating")
    } else if input[len(input)-1:] == "r" {
	rating, _ = strconv.ParseUint(input[0:len(input)-1], 10, 64)
      	level = getLevel(rating)
      	fmt.Println(input[:len(input)-1], "Rating --> Level", level, "Tier", ((level-1)/25)+1, colors[(level-1)%25])
    } else if input[len(input)-1:] == "t" {
	totlevel, _ = strconv.ParseUint(input[0:len(input)-1], 10, 64)
	var i uint64 = 1
	for ; i <= totlevel; i++ {
		rating = getRating(i)
		fmt.Println("\nLevel", i, " Tier", ((i-1)/25)+1, colors[(i-1)%25], " --> ", rating, "rating")
	}
     }
  }
}
