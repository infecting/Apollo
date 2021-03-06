package cli

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

/*
 * Writes a string to the Cli ReadWrite interface
 */
func (c Cli) writeString(text string) error {
	return c.write([]byte(text))
}

/*
 * Writes a byte array to the Cli ReadWrite interface
 */
func (c Cli) write(data []byte) error {
	_, err := c.Writer.Write(data)
	return err
}

/*
 * Reads from the Cli ReadWrite interface
 * This function will return the data returned
 * from the Cli and any errors returned
 *
 * TODO: Add a timeout to this function
 */
func (c Cli) read() (string, error) {
	str, err := bufio.NewReader(c.Writer).ReadString('\n')
	if err != nil {
		return "", err
	}
	str = strings.ReplaceAll(str, "\r", "") // Remove \r in the string
	str = strings.ReplaceAll(str, "\n", "") // Remove \r in the string
	return str, nil
}

/*
 * Closes the ReadWriteCloser and removes the cli from the list
 */
func (c Cli) remove() {
	_ = c.Writer.Close()
	CliList[c] = false
}

/*
 * Prints formatted text to the client.
 * If there is an error writing, the client
 * will be removed
 * Automatically prints a new line to the end of the string if it's not there
 */
func (c Cli) Printf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	if len(text) == 0 { // If there is no text
		return
	}
	if text[len(text)-1] != '\n' { // If the end of the text is not a newline
		text += "\n" // Add a newline
	}
	text = strings.ReplaceAll(text, "\n", "\r\n") // Replaces newlines with \r\n
	err := c.writeString(text)
	if err != nil {
		log.Debugf("error writing to cli: %s", err)
		c.remove()
	}
}

/*
 * Prints the text to the client.
 * If there is an error writing, the client
 * will be removed
 */
func (c Cli) Print(data ...interface{}) {
	text := fmt.Sprint(data...)
	err := c.writeString(text)
	if err != nil {
		log.Debugf("error writing to cli: %s", err)
		c.remove()
	}
}

/*
 * Clears the cli by sending the clear control sequence character
 */
func (c Cli) Clear() {
	c.Print("\033[2J")
}

/*
 * Sets the title of the cli only if color is set to true
 */
func (c Cli) SetTitle(title string) {
	c.Print(fmt.Sprintf("\033]0;%s\007", title))
}
