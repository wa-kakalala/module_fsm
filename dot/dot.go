package dot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DotJson_s struct {
	Dot string `json:"dot"`
}

// 定义一个结构体来存储解析后的数据
type DotData struct {
	Source      string
	Trigger     string
	Destination string
	Color       string
	Type        int
}

// 解析函数
// dotString := `source:"xxx"; trigger:"yyy"; destination:"zzz"; color:"red";`
func parseDotString(input string) (DotData, error) {
	var result DotData
	result.Type = 0
	input = strings.TrimSpace(input)
	pairs := strings.Split(input, ";")
	for _, pair := range pairs {
		fmt.Println(pair)
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			fmt.Println("invalid key-value pair")
			return result, fmt.Errorf("invalid key-value pair: %s", pair)
		}
		key, value := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		value = strings.Trim(value, "\"") // 去除值两端的双引号

		switch key {
		case "source":
			result.Source = value
		case "trigger":
			result.Trigger = value
			if value != "" {
				result.Type += 1
			}
		case "destination":
			result.Destination = value
		case "color":
			result.Type += 2
			result.Color = value
		}
	}
	return result, nil
}

func (dj *DotJson_s) GenDotFile() (string, error) {
//	fmt.Println("enter GetDotFile")
	dotStr := strings.Split(dj.Dot, "\n")
	var edgeinfo []DotData
	var nodeinfo []string
	for _, s := range dotStr {
		dotdata, err := parseDotString(s)
		if err != nil {
			return "", err
		} else {
			edgeinfo = append(edgeinfo, dotdata)
			exist := false
			for _, element := range nodeinfo {
				if element == dotdata.Source {
					exist = true
					break
				}
			}
			if !exist {
				nodeinfo = append(nodeinfo, dotdata.Source)
			}

			exist = false
			for _, element := range nodeinfo {
				if element == dotdata.Destination {
					exist = true
					break
				}
			}
			if !exist {
				nodeinfo = append(nodeinfo, dotdata.Destination)
			}
		}
	}
	// 产生dot文件和fsm
	currentTime := time.Now()
	const format = "20060102150405"
	formattedTime := currentTime.Format(format)
	filename := "./output/" + formattedTime + ".dot"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	content := ""
	content += "digraph graphviz { \n"
	content += "\t graph [\n \t fontname=\"Arial\" \t];\n"
	content += "node [\n\t\tfontsize = \"16\"\n \t\tshape = \"ellipse\" \n \t];\n"
	for _, node := range nodeinfo {
		content += node + ";" + "\n"
	}
	// 添加边
	for _, edge := range edgeinfo {
		content += fmt.Sprintf("\"%v\"->\"%v\"", edge.Source, edge.Destination)
		if edge.Type == 3 {
			content += fmt.Sprintf("[label =\"%v\", color = \"%v\"]", edge.Trigger, edge.Color)
		} else if edge.Type == 1 {
			content += fmt.Sprintf("[label =\"%v\"]", edge.Trigger)
		} else if edge.Type == 2 {
			content += fmt.Sprintf("[ color = \"%v\"]", edge.Color)
		}
	}

	content += "}\n"
//	fmt.Println(content)
	file.Write([]byte(content))
	file.Close()
	pngfile := "./output/" + formattedTime + ".png"                                                        // 确保文件在操作完成后被关闭
	cmd := exec.Command("circo", "-Tpng", "-o", pngfile, filename) // 在 Windows 上可以是 "cmd", "/c", "dir"
//	fmt.Println(cmd)
	// 捕获标准输出和标准错误
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("命令执行失败: %v\n", err)
		fmt.Printf("标准错误输出:\n%s\n", stderr.String())
	} else {
		os.Remove(filename)
	}
	return pngfile, nil
}
