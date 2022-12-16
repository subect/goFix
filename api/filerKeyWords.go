package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"unicode/utf8"
)

type KeyWordFilter struct {
	KeyWord []string `json:"keyWords"` //需要过滤的关键词
	Text    string   `json:"text"`     //需要过滤的文本
}

// Trie 敏感词过滤
type Trie struct {
	child map[rune]*Trie
	word  string
}

// 插入
func (trie *Trie) insert(word string) *Trie {
	cur := trie
	for _, v := range []rune(word) {
		if _, ok := cur.child[v]; !ok {
			newTrie := NewTrie()
			cur.child[v] = newTrie
		}
		cur = cur.child[v]
	}
	cur.word = word
	return trie
}

// 过滤
func (trie *Trie) filerKeyWords(word string) string {
	cur := trie
	for i, v := range []rune(word) {
		if _, ok := cur.child[v]; ok {
			cur = cur.child[v]
			if cur.word != "" {
				word = replaceStr(word, "*", i-utf8.RuneCountInString(cur.word)+1, i)
				cur = trie
			}
		} else {
			cur = trie
		}
	}
	return word
}
func replaceStr(word string, replace string, left, right int) string {
	str := ""
	for i, v := range []rune(word) {
		if i >= left && i <= right {
			str += replace
		} else {
			str += string(v)
		}
	}
	return str
}
func NewTrie() *Trie {
	return &Trie{
		word:  "",
		child: make(map[rune]*Trie, 0),
	}
}

func FilerKeyWords(c *gin.Context) {
	trie := NewTrie()

	//str := c.Query("text")
	//trie.insert(str)
	//trie.filerKeyWords()

	r := c.Request
	req := KeyWordFilter{}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		basicLog.Errorf("FlowInteraction getReq err %v bs [%s]\n", err, bs)
		return
	}
	err = json.Unmarshal(bs, &req)
	if err != nil {
		basicLog.Errorf("FlowInteraction Unmarshal err %v bs [%s]\n", err, bs)
		return
	}

	for _, keyWord := range req.KeyWord {
		trie.insert(keyWord)
	}
	result := trie.filerKeyWords(req.Text)
	basicLog.Debugf(result)
	c.JSON(200, result)
}
