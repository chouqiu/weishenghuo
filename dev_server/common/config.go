package common 

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//解析类
type Config struct {
	filepath string
	conflist map[string]map[string]string
}

//创建一个空的解析文件类
func LoadConfig(filepath string) *Config {
	c := new(Config)
	c.filepath = filepath

	c.ReadList()

	return c
}

//根据section和name获取value
func (c *Config) GetValue(section, name string) string {

	for k, v := range c.conflist {

		if k == section {
			return v[name]
		}

	}
	return ""
}

//读取文件
func (c *Config) ReadList() {

	file, err := os.Open(c.filepath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var section string = ""

	c.conflist = make(map[string]map[string]string)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				continue
			}
		}

		line = strings.TrimSpace(line)

		switch {
		case len(line) == 0:
			//去除空行
		case line[0] == '#':
			//去除注释行
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			c.conflist[section] = make(map[string]string)
		default:
			i := strings.Index(line, "#")
			//去除行尾注释
			if i >= 0 {
				line = line[:i]
			}

			i = strings.Index(line, "=")

			if i >= 0 {
				c.conflist[section][strings.TrimSpace(line[0:i])] = strings.TrimSpace(line[i+1:])
			}
		}
	}
}
