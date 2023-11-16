package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// 输入文法 (可修改为文件读入)
var input = []string{
	//"E->aA|bB",
	//"A->cA|d",
	//"B->cB|d",

	"S->BB",
	"B->aB",
	"B->b",
}

var (
	startChar      string                              // 开始符号
	OVER           int                                 // 开始符号结束 S'->S.
	END            []int                               // 结束状态集合  A->a.
	convert        = make(map[int]map[string]int)      // convert 用于记录状态转换 eg: 0 -> 1 [0]状态经过'E' -> [2]状态
	handlerGrammar = make(map[string][]string)         // 处理后的input
	dfa            = make(map[int]map[string][]string) // 状态集合
	rightList      []string                            // 非终结符集合
	leftList       []string                            // 终结符结合
)

// Stack 自定义栈
type Stack []interface{}

func (s *Stack) Push(value interface{}) {
	*s = append(*s, value)
}

func (s *Stack) Pop() interface{} {
	if len(*s) == 0 {
		return nil
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value
}

func (s *Stack) Peek() interface{} {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func main() {

	// 1 对 input 进行处理 //
	// 1.1 将input处理成map类型方便后续处理  E->aA|bB 处理为 E->aA 和 E->bB 存入handlerGrammar
	handlerGrammar := make(map[string][]string)

	for _, grammar := range input {

		// 取出E->aA|bB 中的 aA|bB
		str := strings.Split(grammar, "->") // fmt.Println(rightStr)
		leftStr := str[0]
		leftList = append(leftList, leftStr) // 将非终结符存入leftList
		rightStr := str[1]
		// 以'|' 分割aA|bB
		rightGrammars := strings.Split(rightStr, "|")
		for _, rg := range rightGrammars {
			handlerGrammar[leftStr] = append(handlerGrammar[leftStr], rg)
		}
	}

	// 1.2 找出开始符号（未在右端出现的非终结符)
	for _, r := range leftList {
		isStart := true
		for _, grammar := range handlerGrammar {
			for _, s := range grammar {
				if strings.Contains(s, r) {
					isStart = false
					break
				}
			}
		}
		if isStart {
			startChar = r
		}
	}

	// 1.3 将所有字符存入set
	var pSet []string
	for k, v := range handlerGrammar {
		pSet = append(pSet, k)
		for _, str := range v {
			for _, c := range str {
				s := string(c)
				if !slices.Contains(leftList, s) {
					rightList = append(rightList, s)
				}
				pSet = append(pSet, s)
			}
		}
	}
	slices.Sort(pSet)

	// 添加开始状态 eg: S'-> S
	handlerGrammar[startChar+"'"] = append(handlerGrammar[startChar+"'"], startChar)

	// 2.生成拓广文法 //
	expandGrammar := generateExpandGrammar(handlerGrammar)
	firstExpandGrammar := make(map[string][]string) // 得到所有 以'.'开始的语句
	for s, rList := range expandGrammar {
		for _, str := range rList {
			if str[0] == '.' {
				firstExpandGrammar[s] = append(firstExpandGrammar[s], str)
			}

		}
	}
	fmt.Println("// ================== 拓广文法集合 =============== //")
	for k, eg := range expandGrammar {
		for _, s := range eg {
			fmt.Println(k, "->", s)
		}
	}
	fmt.Println()

	// 3.DFA 生成 //

	// 3.1 构造I0, 找到开始符号 将其加入dfa[0]
	keyMap := make(map[string][]string)
	keyMap[startChar+"'"] = append(keyMap[startChar+"'"], "."+startChar)

	idx := 0
	dfa[idx] = keyMap

	var dfs func() //定义递归函数

	// 3.2 依次遍历每个字符 进行状态转换(深度搜索）
	dfs = func() {
		// 构造I[idx]集合
		for key, rList := range dfa[idx] {
			for _, str := range rList {
				if str[len(str)-1] == '.' { // 如果该语句为"aB." 以'.'结尾 则无须继续判断
					if str == startChar+"." { // 是否为 开始符号的归约 eg: S'->S.
						OVER = idx // 记录在第几状态
					}
					continue
				}
				for i := 1; i < len(str); i++ {
					// expandGrammar[key] 对应的语句全加入当前状态集合
					cs := string(str[i])
					if slices.Contains(leftList, cs) && str[i-1] == '.' {
						for _, s := range firstExpandGrammar[cs] {
							dfa[idx][cs] = append(dfa[idx][cs], s)
						}
					}
				}
			}
			slices.Sort(dfa[idx][key]) // 将状态内的文法进行排序 方便后续判断是否到达同一状态
		}

		// 检查这个状态是否已经存在
		for i := 0; i < idx; i++ {
			isSame := true
			if len(dfa[i]) != len(dfa[idx]) { // 两状态所含产生式数量不一致
				continue
			}
			for sKey, list := range dfa[idx] { // 遍历要判断的状态中的每个产生式
				if _, ok := dfa[i][sKey]; ok {
					if len(list) != len(dfa[idx][sKey]) {
						isSame = false
						break
					}
					for j := range dfa[idx][sKey] {
						if dfa[idx][sKey][j] != dfa[i][sKey][j] {
							isSame = false
							break
						}
					}
				} else {
					isSame = false
					break
				}
				if !isSame {
					break
				}
			}
			if isSame {
				delete(dfa, idx)
				return
			}
		}

		preIdx := idx // 记录当前状态的下标
		idx++
		//状态转移
		for _, ch := range pSet {
			for key, rList := range dfa[preIdx] { // key 为产生式左边部分
				for _, str := range rList {
					if str[len(str)-1] == '.' { // 无法发生转移的产生式 E -> a.
						continue
					}
					for i := 1; i < len(str); i++ {
						s := string(str[i])
						if s == ch && str[i-1] == '.' {
							newStr := str[:i-1] + s + string(str[i-1]) + str[i+1:] // 变换为转换后的语句 eg: a.B => aB.
							keyMap = make(map[string][]string)
							keyMap[key] = append(keyMap[key], newStr)
							dfa[idx] = keyMap
							dfs()
						}
					}
				}
			}
		}

	}
	dfs()

	/* 3.3 构造状态转移信息
	i为要转移的状态I下标, i2 为可能转移的状态I下标
	在dfa[i]中找出会发生转移的字符str, 循环所有的dfa, 如果在dfa中找出一个可以由字符str 转移得到
	即 dfa[i] -> dfa[i2], 加入convert中
	*/
	for i, m := range dfa {
		for _, rList := range m {
			for _, cs := range rList {
				if cs[len(cs)-1] == '.' { // 无法发生转移
					continue
				}
				index := strings.Index(cs, ".")
				str := string(cs[index+1])
				cs = cs[:index] + string(cs[index+1]) + "." + cs[index+2:] // 将B.B -> BB.
				isTranslate := false                                       // 是否转移成功
				for i2, m2 := range dfa {                                  // 可以自己转移到自己
					for _, rList2 := range m2 {
						if slices.Contains(rList2, cs) {
							if convert[i] == nil {
								nxtFlow := make(map[string]int)
								nxtFlow[str] = i2
								convert[i] = nxtFlow
							} else {
								convert[i][str] = i2
							}
							isTranslate = true
							break
						}
					}
					if isTranslate {
						break
					}
				}
			}
		}
	}
	fmt.Println("// ================== DFA =============== //")

	for i := 0; i < len(dfa); i++ {
		fmt.Println(fmt.Sprintf("I%d", i))
		for key, list := range dfa[i] {
			for _, s := range list {
				fmt.Println(key, "->", s)
			}
		}
		fmt.Println()
	}
	fmt.Println()
	// 将未发生转移的状态存入END中 例: B->b.
	for i := 0; i <= idx; i++ {
		if _, ok := convert[i]; !ok {
			END = append(END, i)
		}
	}

	var (
		cToI = make(map[string]int) // 将每个字符映射到数组的第几个位置
	)
	pIdx := 0
	for _, s := range rightList {
		cToI[s] = pIdx
		pIdx++
	}
	rightList = append(rightList, "#")
	cToI["#"] = pIdx
	pIdx++
	for _, s := range leftList {
		cToI[s] = pIdx
		pIdx++
	}
	analyzeTable := make(map[int][]string)
	for i := 0; i < idx; i++ {

		if slices.Contains(END, i) && i != OVER {
			var run int
			for _, rList := range dfa[i] {
				for _, s := range rList {
					s = s[:len(s)-1]
					for i, s2 := range input {
						if strings.Contains(s2, s) {
							run = i + 1
							break
						}
					}
				}
			}
			for range rightList {
				str := "r" + strconv.Itoa(run)
				analyzeTable[i] = append(analyzeTable[i], str)
			}
		} else {
			for _, s := range rightList {
				if nxtId, ok := convert[i][s]; ok {
					flag := fmt.Sprintf("S%d", nxtId)
					analyzeTable[i] = append(analyzeTable[i], flag)
				} else if i == OVER && s == "#" {
					analyzeTable[i] = append(analyzeTable[i], "acc")
				} else {
					analyzeTable[i] = append(analyzeTable[i], " ")
				}
			}
		}
		for _, s := range leftList {
			if nxtId, ok := convert[i][s]; ok {
				analyzeTable[i] = append(analyzeTable[i], strconv.Itoa(nxtId))
			} else {
				analyzeTable[i] = append(analyzeTable[i], " ")
			}
		}
	}
	fmt.Println("// ================== LR(0)分析表 =============== //")
	fmt.Println("\tAction \t\t Goto")
	fmt.Println("--------------------------------------------------------")
	for _, s := range rightList {
		fmt.Printf("\t%v", s)
	}
	for _, s := range leftList {
		fmt.Printf("\t%v", s)
	}
	fmt.Println("\n-------------------------------------------------------")
	for i := 0; i < len(analyzeTable); i++ {
		fmt.Printf("%d\t", i)
		for _, s := range analyzeTable[i] {
			fmt.Printf("%v\t", s)
		}
		fmt.Println()
	}
	fmt.Println()

	// 4.利用LR(0)分析表分析语句
	// 4.1 初始化状态栈 与 归约栈
	convertStk := make(Stack, 0)
	convertStk.Push(0)

	reductionStk := make(Stack, 0)
	reductionStk.Push("#")

	// 4.2 读入要分许的语句
	var inputStr string
	fmt.Print("输入以#结束的字符串: ")
	_, _ = fmt.Scan(&inputStr)
	isSuccess := true // 判断是否归约成功

	// 4.3 进行归约
	k := 0
	for k < len(inputStr) {
		cs := string(inputStr[k])
		i := convertStk.Peek().(int)   // 取出状态栈的栈顶状态
		s := analyzeTable[i][cToI[cs]] // 根据分析表得到可能发生的转移s ["r%d", "S%d", "acc", "%d"]
		if s != " " {
			if s[0] == 'r' { // 发生归约
				// action
				zIdx, _ := strconv.Atoi(s[1:])               // 使用第几个规则进行归约
				cvt := strings.Split(input[zIdx-1], "->")[0] //得到推导式的左部

				// 对产生式的右半边进行遍历判断进行归约的是哪个语句 得到语句的长度
				parts := strings.Split(input[zIdx-1], "->")[1]
				partList := strings.Split(parts, "|")
				var popSize int
				for _, s2 := range partList {
					if string(s2[len(s2)-1]) == reductionStk.Peek().(string) {
						popSize = len(s2)
						break
					}
				}

				// 状态回退
				for i := 0; i < popSize; i++ {
					convertStk.Pop()
					reductionStk.Pop()
				}

				reductionStk.Push(cvt) // 将产生式左边加入归约栈

				// goto
				j := convertStk.Peek().(int)
				p := analyzeTable[j][cToI[cvt]]
				pId, _ := strconv.Atoi(p)
				convertStk.Push(pId)
			} else if s[0] == 'S' {
				zIdx, _ := strconv.Atoi(s[1:])
				convertStk.Push(zIdx)
				reductionStk.Push(cs)
				k++
			} else if s == "acc" { // "acc" 归约成功 跳出循环
				break
			}
		} else {
			isSuccess = false
			break
		}

	}
	if isSuccess {
		fmt.Println("归约成功")
	} else {
		fmt.Println("归约失败")
	}

}

func generateExpandGrammar(grammar map[string][]string) map[string][]string {
	expandGrammar := make(map[string][]string)
	for k, vStrList := range grammar {
		for _, vStr := range vStrList {
			var g string
			for i := range vStr {
				if i == 0 {
					g = "." + vStr
				} else {
					g = vStr[:i] + "." + vStr[i:]
				}
				expandGrammar[k] = append(expandGrammar[k], g)
			}
			expandGrammar[k] = append(expandGrammar[k], vStr+".")
		}

	}
	return expandGrammar
}
