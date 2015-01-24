package main

import(
    "fmt"
    "net/http"
    "os"

    "github.com/wsxiaoys/terminal/color"
)

/* holds our palette, or 16 ANSI colors (8 normal colors + bright complements)
   and two foreground/background colors. colors are 3 byte arrays (RGB) */
type ktPalette struct {

    black    [3]byte
    bblack   [3]byte

    red      [3]byte
    bred     [3]byte

    green    [3]byte
    bgreen   [3]byte

    yellow   [3]byte
    byellow  [3]byte

    blue     [3]byte
    bblue    [3]byte

    purple   [3]byte
    bpurple  [3]byte

    cyan     [3]byte
    bcyan    [3]byte

    white    [3]byte
    bwhite   [3]byte

    fg       [3]byte
    bg       [3]byte
}

/* the default "control" ANSI color set (boring) */
var ansiColors = ktPalette {

    black   : [3]byte {0,0,0},
    bblack  : [3]byte {128,128,128},

    red     : [3]byte {128,0,0},
    bred    : [3]byte {255,0,0},

    green   : [3]byte {0,128,0},
    bgreen  : [3]byte {0,255,0},

    yellow  : [3]byte {128,128,0},
    byellow : [3]byte {255,255,0},

    blue    : [3]byte {0,0,128},
    bblue   : [3]byte {0,0,255},

    purple  : [3]byte {128,0,128},
    bpurple : [3]byte {255,0,255},

    cyan    : [3]byte {0,128,128},
    bcyan   : [3]byte {0,255,255},

    white   : [3]byte {128,128,128},
    bwhite  : [3]byte {255,255,255},

    fg      : [3]byte {0,0,0},
    bg      : [3]byte {255,255,255},
}

/* parses a colorfile, returns palette struct. given a nil file pointer,
   returns standard ANSI color set (our "control") */
func parseColors(colorfile *os.File) (pal ktPalette, err error) {

    return ansiColors, nil
}

func ktInit(dirPrepend string, port int, colorfilePath string) error {

    color.Print("@yparsing colorfile :: @{|}")
    file, err := os.Open(colorfilePath)
    if err != nil {
        color.Printf("@r[%s]@{|} - bad colorfile path\n", xM)
        return fmt.Errorf("%s\n", "bad colorfile path")
    }

    pal, err := parseColors(file)
    fmt.Print(pal)

    if err != nil {
        color.Printf("@r[%s]@{|} - malformed colorfile\n", xM)
        return fmt.Errorf("%s\n", "malformed colorfile")
    }

    color.Printf("@g[%s]@{|}\n", checkM)

    color.Printf("@ystarting httpd on port @b%d@{|} :: ", port)

    return nil
}

func main() {

    err := ktInit("kolors", 999, "/home/kremlin/go/src/k_transpose/kremlin_colors");

    /* make sure to close() anything you need to (you need to) */
    resp, err := http.Get("http://kremlin.cc")

    if err != nil {}
    fmt.Println(resp)
}

var checkM = "✓"
var xM     = "✗"

