package main

import "fmt"

type Cell func() string

type CellBorder rune

type Color string

func (cb CellBorder) String() string {
	return fmt.Sprintf("%c", cb)
}

func createTemplate(emoji, name, value string) Cell {
	return func() string {
		return emoji + " " + name + ": " + value
	}
}

func colorPrint(str string, c Color) {
	colorsCode := map[Color]string{
		"reset":   "\033[0m",
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
	}

	fmt.Print(colorsCode[c] + str + colorsCode["reset"])
}

const maxCharCount = 39

func drawTemplateWithBorders(str string, c Color) {
	colorPrint("|", c)

	charCount := 1
	for _, ch := range str {
		charCount++

		// move to the next line if str is so long
		if charCount == maxCharCount {
			charCount = 1

			colorPrint("|\n", c)
			colorPrint("|", c)
		}

		fmt.Printf("%c", ch)
	}

	// correct draw right border
	for charCount < maxCharCount-1 {
		fmt.Print(" ")
		charCount++
	}
	colorPrint("|\n", c)
}

func drawBetweenBorder(between CellBorder, c Color) {
	colorPrint("|", c)

	for i := 0; i < 38; i++ {
		colorPrint(string(between), c)
	}
	colorPrint("|\n", c)
}

func drawCell(between CellBorder, c Color, name Cell, templates ...Cell) {
	colorPrint("----------------------------------------\n", c)

	drawTemplateWithBorders(name(), c)

	for _, template := range templates {
		if between != 0 {
			drawBetweenBorder(between, c)
		}

		drawTemplateWithBorders(template(), c)
	}

	colorPrint("----------------------------------------\n", c)
}

func main() {
	var between CellBorder = '$'

	var color Color = "yellow"

	name := createTemplate("💬", "Название", "станок сторона китай качественный гайка")
	drawCell(between, color, name)

	color = "green"
	description := createTemplate("📖", "Описание", "станок для дерева")
	drawCell(between, color, name, description)

	between = '#'
	price := createTemplate("💵", "Цена", "100")
	drawCell(between, color, name, description, price)

	between = 0
	color = "red"
	location := createTemplate("📍", "Локация, где можно будет забрать товар, который будет получен", "Москва")
	drawCell(between, color, name, description, price, location)
}
