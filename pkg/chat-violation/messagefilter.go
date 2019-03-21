package chat_violation

import (
    "unicode"
    "strings"
    "os"
    "bufio"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/chat-violation/trie"
)

type ChatViolation struct {
    Tree *trie.Trie
}

var MessageFilter ChatViolation

// Tokenize by first seperating chinese characters one by one and english by spaces, and then lowercased and all chinese is converted into simplified chinese
// all tabs and spaces are ignored between chinese characters
//
func (c *ChatViolation) tokenize(s string) (g []string) {
    g = make([]string, 0, len(s))
    temp := make([]rune, 0)

    for _, c := range s {
        if c <= unicode.MaxASCII {
            if c == ' ' || c == '\t' {
                if len(temp) > 0 {
                    g = append(g, strings.ToLower(string(temp)))
                    temp = make([]rune, 0)
                }
            } else {
                temp = append(temp, c)

            }
        } else {
            if len(temp) > 0 {
                g = append(g, strings.ToLower(string(temp)))
                temp = make([]rune, 0)
            }

            g = append(g, string(c))

        }
    }
    if len(temp) > 0 {
        g = append(g, strings.ToLower(string(temp)))
    }
    return
}

func (c *ChatViolation) InitChatViolation(filePath string) {
    c.Tree = trie.NewTrie()
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    scan := bufio.NewScanner(file)
    for scan.Scan() {
        s := scan.Text()

        g := c.tokenize(s)
        c.Tree.Insert(g)
    }
}

func (c *ChatViolation) WordsFilter(msg string) string{
    found := 0
    finstr := ""

    if len(msg) > 0 {
        g := c.tokenize(msg)
        if len(g) == 0 {
            found = 0
            return string([]byte(msg))
        }
        b := g
        for {
            found = <-c.Tree.FindMatchLen(b)
            if len(b) <= 1 {
                break
            }
            if found > 0 {
                if len(b) > found {
                    b = b[found:]
                    continue
                }
                break
            } else {
                finstr = finstr + b[0]
                b = b[1:]
            }
        }
        found = <-c.Tree.FindMatchLen(b)
        if found == 0 {
            finstr = finstr + b[0]
        }
    }
    return string([]byte(finstr))
}