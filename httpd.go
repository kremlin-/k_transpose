package main

import (
	"fmt"
	"os"

	"bufio"
	"io"
	"io/ioutil"

	"strconv"
	"strings"

	"net/http"
	"net/url"

	"github.com/wsxiaoys/terminal/color"
)

/* holds our palette, or 16 ANSI colors (8 normal colors + bright complements)
   and two foreground/background colors. colors are 3 byte arrays (RGB) */
type ktPalette struct {
	black  [3]byte
	bblack [3]byte

	red  [3]byte
	bred [3]byte

	green  [3]byte
	bgreen [3]byte

	yellow  [3]byte
	byellow [3]byte

	blue  [3]byte
	bblue [3]byte

	purple  [3]byte
	bpurple [3]byte

	cyan  [3]byte
	bcyan [3]byte

	white  [3]byte
	bwhite [3]byte

	fg [3]byte
	bg [3]byte
}

/* the default "control" ANSI color set (boring) */
var ansiColors = ktPalette{

	black:  [3]byte{0, 0, 0},
	bblack: [3]byte{128, 128, 128},

	red:  [3]byte{128, 0, 0},
	bred: [3]byte{255, 0, 0},

	green:  [3]byte{0, 128, 0},
	bgreen: [3]byte{0, 255, 0},

	yellow:  [3]byte{128, 128, 0},
	byellow: [3]byte{255, 255, 0},

	blue:  [3]byte{0, 0, 128},
	bblue: [3]byte{0, 0, 255},

	purple:  [3]byte{128, 0, 128},
	bpurple: [3]byte{255, 0, 255},

	cyan:  [3]byte{0, 128, 128},
	bcyan: [3]byte{0, 255, 255},

	white:  [3]byte{128, 128, 128},
	bwhite: [3]byte{255, 255, 255},

	fg: [3]byte{0, 0, 0},
	bg: [3]byte{255, 255, 255},
}

/* lets get our money's worth out of this utf-8 crap */
var checkM = "✓"
var xM = "✗"

/* parses a colorfile, returns palette struct. given a nil file pointer,
   returns standard ANSI color set (our "control") */
func parseColors(colorfile *os.File) (pal ktPalette, err error) {

	if colorfile == nil {
		return ansiColors, nil
	}

	ret := ktPalette{}
	var e errConsumer

	scanner := bufio.NewScanner(colorfile)
	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "color") &&
			!strings.Contains(scanner.Text(), "!") {

			/* i am so sorry */
			cur := strings.Replace(scanner.Text(), " ", "", -1)
			split := strings.Split(cur, ":")
			hexColor := strings.Replace(split[1], "#", "", -1)
			coNumStr := strings.Replace(split[0], "URxvt*color", "", -1)

			r, err := strconv.ParseUint(hexColor[0:2], 16, 8)
			e.Consume(err)
			g, err := strconv.ParseUint(hexColor[2:4], 16, 8)
			e.Consume(err)
			b, err := strconv.ParseUint(hexColor[4:], 16, 8)
			e.Consume(err)

			colorNo, err := strconv.Atoi(coNumStr)
			e.Consume(err)

			if e.err != nil {
				return ktPalette{}, e.err
			}

			switch colorNo {
			case 0:
				ret.black[0] = byte(r)
				ret.black[1] = byte(g)
				ret.black[2] = byte(b)
			case 1:
				ret.red[0] = byte(r)
				ret.red[1] = byte(g)
				ret.red[2] = byte(b)
			case 2:
				ret.green[0] = byte(r)
				ret.green[1] = byte(g)
				ret.green[2] = byte(b)
			case 3:
				ret.yellow[0] = byte(r)
				ret.yellow[1] = byte(g)
				ret.yellow[2] = byte(b)
			case 4:
				ret.blue[0] = byte(r)
				ret.blue[1] = byte(g)
				ret.blue[2] = byte(b)
			case 5:
				ret.purple[0] = byte(r)
				ret.purple[1] = byte(g)
				ret.purple[2] = byte(b)
			case 6:
				ret.cyan[0] = byte(r)
				ret.cyan[1] = byte(g)
				ret.cyan[2] = byte(b)
			case 7:
				ret.white[0] = byte(r)
				ret.white[1] = byte(g)
				ret.white[2] = byte(b)
			case 8:
				ret.bblack[0] = byte(r)
				ret.bblack[1] = byte(g)
				ret.bblack[2] = byte(b)
			case 9:
				ret.bred[0] = byte(r)
				ret.bred[1] = byte(g)
				ret.bred[2] = byte(b)
			case 10:
				ret.bgreen[0] = byte(r)
				ret.bgreen[1] = byte(g)
				ret.bgreen[2] = byte(b)
			case 11:
				ret.byellow[0] = byte(r)
				ret.byellow[1] = byte(g)
				ret.byellow[2] = byte(b)
			case 12:
				ret.bblue[0] = byte(r)
				ret.bblue[1] = byte(g)
				ret.bblue[2] = byte(b)
			case 13:
				ret.bpurple[0] = byte(r)
				ret.bpurple[1] = byte(g)
				ret.bpurple[2] = byte(b)
			case 14:
				ret.bcyan[0] = byte(r)
				ret.bcyan[1] = byte(g)
				ret.bcyan[2] = byte(b)
			case 15:
				ret.bwhite[0] = byte(r)
				ret.bwhite[1] = byte(g)
				ret.bwhite[2] = byte(b)
			}

		} else if strings.Contains(scanner.Text(), "background") {

			hex := strings.Split(scanner.Text(), "#")
			het := hex[1]

			r, err := strconv.ParseUint(het[0:2], 16, 8)
			e.Consume(err)
			g, err := strconv.ParseUint(het[2:4], 16, 8)
			e.Consume(err)
			b, err := strconv.ParseUint(het[4:], 16, 8)
			e.Consume(err)

			if e.err != nil {
				return ktPalette{}, e.err
			}

			ret.bg[0] = byte(r)
			ret.bg[1] = byte(g)
			ret.bg[2] = byte(b)

		} else if strings.Contains(scanner.Text(), "foreground") {

			hex := strings.Split(scanner.Text(), "#")
			het := hex[1]

			r, err := strconv.ParseUint(het[0:2], 16, 8)
			e.Consume(err)
			g, err := strconv.ParseUint(het[2:4], 16, 8)
			e.Consume(err)
			b, err := strconv.ParseUint(het[4:], 16, 8)
			e.Consume(err)

			if e.err != nil {
				return ktPalette{}, e.err
			}

			ret.fg[0] = byte(r)
			ret.fg[1] = byte(g)
			ret.fg[2] = byte(b)
		}
	}

	return ret, e.err
}

var httpdStatus bool

func ktInit(dirPrepend string, port int, colorfilePath string) error {

	httpdStatus = false

	color.Print("@yparsing colorfile :: @{|}")
	file, err := os.Open(colorfilePath)
	if err != nil {
		color.Printf("@r[%s]@{|} - bad colorfile path\n", xM)
		return fmt.Errorf("%s\n", "bad colorfile path")
	}

	pal, err := parseColors(file)

	if err != nil {
		color.Printf("@r[%s]@{|} - malformed colorfile [%s]\n", xM, err)
		return fmt.Errorf("%s\n", "malformed colorfile")
	}

	color.Printf("@g[%s]@{|}\n", checkM)

	if pal == ansiColors {
	}

	color.Printf("@yverifying & preprocessing colorsets@{|} :: @y[SKIP]\n")
	color.Printf("@ygenerating transpositional colorspace@{|} :: @y[SKIP]\n")

	color.Printf("@ystarting httpd on port @b%d@{|} :: ", port)
	http.HandleFunc("/kt/", transposePage)

	var portString string
	fmt.Sprintf(portString, ":%d", port)
	err = http.ListenAndServe(portString, nil)
	if err != nil {

		color.Printf("@r[%s]@{-} - run me as root!\n", xM)
		return fmt.Errorf("%s\n", "failed to start httpd")
	}

	color.Printf("@g[%s]@{-}\n", checkM)

	return nil
}

func transposePage(writer http.ResponseWriter, req *http.Request) {

	if !httpdStatus {

		httpdStatus = true
		color.Printf("@g[%s]@{|}\n", checkM)
	}

	if req.URL.Path == "/kt/" {
		writer.Write([]byte("wtf"))
		return
	}

	fqdn := req.URL.Path[4:]
	targetURL := fmt.Sprintf("http://%s", fqdn)

	resp, err := http.Get(targetURL)

	if err != nil || resp.StatusCode != 200 {

		io.WriteString(writer, "failed to get that page! -kt\n")
		io.WriteString(writer, targetURL+"\n")

		io.WriteString(writer, resp.Status)
		return
	}

	conType := resp.Header.Get("Content-Type")

	switch conType[0:strings.Index(conType, ";")] {

	case "text/html":
		writer.Write(transposeHTML(bufio.NewScanner(resp.Body), fqdn))

	case "text/css":
		writer.Write(transposeCSS(bufio.NewScanner(resp.Body), fqdn))

	default:
		page, _ := ioutil.ReadAll(resp.Body)
		writer.Write(page)
	}

	resp.Body.Close()
}

/* swap href="" & src="" */
func transposeHTML(scan *bufio.Scanner, fqdn string) []byte {

	var ret []byte
	var i int

	scan.Split(bufio.ScanWords)
	for scan.Scan() {

		i++
		cur := scan.Text()

		//fmt.Printf("%s\n", cur)

		if len(cur) < 7 {

		} else if cur[0:6] == "href=\\" {

			urlStr := cur[7 : strings.Index(cur[7:], "\\")+7]

			u, err := url.Parse(urlStr)
			if err != nil {
				fmt.Printf("malformed URL: %s\n", urlStr)
			}

			if u.Host == "" {

				u.Host = fmt.Sprintf("localhost/kt/%s", fqdn)
				//				cur = append(cur[0:6],
			}

			fmt.Printf("[F] URL: %s // PATH: %s\n", u.Host, u.Path)
			if u == u {
			}

		} else if cur[0:5] == "href=" {

			urlStr := cur[6 : strings.Index(cur[6:], "\"")+6]

			u, err := url.Parse(urlStr)
			if err != nil {
				fmt.Printf("malformed URL: %s\n", urlStr)
			}

			if u == u {
			}
			fmt.Printf("URL: %s // PATH: %s\n", u.Host, u.Path)

		} else if cur[0:5] == "src=\"" {

			//fmt.Printf("%s\n", cur)
			urlStr := cur[5 : strings.Index(cur[5:], "\"")+5]

			u, err := url.Parse(urlStr)
			if err != nil {
				fmt.Printf("malformed URL: %s\n", urlStr)
			}

			if u.Host == "" {

				u.Host = fmt.Sprintf("localhost/kt/%s", fqdn)
				//				cur = append(cur[0:6],

				fmt.Printf("[S] URL: %s // PATH: %s\n", u.Host, u.Path)
				if u == u {
				}
			}

		}

		ret = append(ret, byte(' '))
		ret = append(ret, cur...)
	}

	fmt.Printf("%d\n", i)

	return ret
}

func transposeCSS(scan *bufio.Scanner, fqdn string) []byte {

	var ret []byte

	return ret
}

func main() {

	err := ktInit("kolors", 999, "/home/kremlin/.Xresources")

	if err != nil {
	}
}

type errConsumer struct {
	err error
}

func (e *errConsumer) Consume(err error) {
	if e.err == nil && err != nil {
		e.err = err
	}
}
