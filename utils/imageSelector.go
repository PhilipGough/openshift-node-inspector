package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const baseImagePath string = "docker.io/rhmap/"

func ConfirmImage(image string, isDefault bool) string {
	image = determineImage(image, isDefault)
	fmt.Printf("Is this the image you want to use --- > %s? \nPress enter to accept or input different value \n", image)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.TrimSpace(text) != "" {
		image = text

	}
	fmt.Printf("Using image %s \n", image)
	return image
}

func determineImage(image string, useDefault bool) string {
	if useDefault {
		return image
	}
	return baseImagePath + image
}
