package main

import "fmt"

type drawSubCell func() string

type ColorCode int

type Color string

type CellBorder rune

type CellShape func(CellBorder, Color, drawSubCell, []drawSubCell)

const (
	codeReset = ColorCode(iota)
	codeBlack = ColorCode(iota + 29)
	codeRed
	codeGreen
	codeYellow
	codeBlue
	codeMagenta
	codeCyan
	codeWhite
)

func convertCodeToColor(color ColorCode) Color {
	return Color(fmt.Sprintf("\033[%dm", color))
}

func createTemplate(emoji, name, value string) drawSubCell {
	return func() string {
		return fmt.Sprintf("%s %s: %s", emoji, name, value)
	}
}

func printColoredString(str string, c Color) {
	fmt.Printf("%s%s%s", c, str, convertCodeToColor(codeReset))
}

const maxLengthRowCell = 40 // cell width

func drawTemplateWithBorders(subCell string, c Color) {
	printColoredString("|", c)

	charCount := 1
	for _, ch := range subCell {
		charCount++

		// move to the next line if str is so long
		if charCount == maxLengthRowCell-1 {
			charCount = 1

			printColoredString("|\n", c)
			printColoredString("|", c)
		}

		fmt.Printf("%c", ch)
	}

	// correct draw right border
	for charCount < maxLengthRowCell-2 {
		fmt.Print(" ")
		charCount++
	}
	printColoredString("|\n", c)
}

func drawBetweenBorder(borderBetween CellBorder, c Color) {
	printColoredString("|", c)

	for i := 0; i < maxLengthRowCell-2; i++ {
		printColoredString(string(borderBetween), c)
	}
	printColoredString("|\n", c)
}

func drawCell(drawForm CellShape, borderBetween CellBorder, c Color, titleCell drawSubCell, templates ...drawSubCell) {
	drawForm(borderBetween, c, titleCell, templates)
}

func drawSquare(borderBetween CellBorder, c Color, name drawSubCell, templates []drawSubCell) {
	printColoredString("----------------------------------------\n", c)

	drawTemplateWithBorders(name(), c)

	for _, template := range templates {
		if borderBetween != 0 {
			drawBetweenBorder(borderBetween, c)
		}

		drawTemplateWithBorders(template(), c)
	}

	printColoredString("----------------------------------------\n", c)
}

func main() {
	var borderBetween CellBorder = '$'

	color := convertCodeToColor(codeYellow)

	name := createTemplate("💬", "Название", "станок сторона китай качественный гайка")

	var square CellShape = func(borderBetween CellBorder, c Color, name drawSubCell, templates []drawSubCell) {
		drawSquare(borderBetween, c, name, templates)
	}

	drawCell(square, borderBetween, color, name)

	color = convertCodeToColor(codeGreen)
	description := createTemplate("📖", "Описание", "станок для дерева")
	drawCell(square, borderBetween, color, name, description)

	borderBetween = '#'
	price := createTemplate("💵", "Цена", "100")
	drawCell(square, borderBetween, color, name, description, price)

	borderBetween = 0
	color = convertCodeToColor(codeRed)
	location := createTemplate("📍", "Локация, где можно будет забрать товар, который будет получен", "Москва")
	drawCell(square, borderBetween, color, name, description, price, location)
}
