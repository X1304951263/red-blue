package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

type Card struct {
	Suit  string
	Value string
}

var deck []Card //牌盒
var m map[string]int

func main() {

	sampleNum := 1 << 3
	bannerWinNum := make([]int, 72)
	playerWinNum := make([]int, 72)
	drawNum := make([]int, 72)
	for i := 0; i < sampleNum; i++ {

		//初始化路牌
		roadSign := make([]string, 0)

		// 初始化一副扑克牌
		deck = initializeDeck()

		// 初始化分数
		m1 := make(map[string]int)
		initializeScoreMap(m1)
		m = m1
		// 洗牌
		shuffledDeck := shuffleDeck(deck)
		shuffledDeck = removeDeckHead(shuffledDeck)

		totalNums, bankerWin, playerWin, draw := 0, 0, 0, 0
		luckySix := 0
		for len(shuffledDeck) > 6 {
			playerCard := make([]Card, 0)
			bankerCard := make([]Card, 0)
			// 发牌
			bankerCard = dealCard(bankerCard, &shuffledDeck, 1)
			playerCard = dealCard(playerCard, &shuffledDeck, 1)
			bankerCard = dealCard(bankerCard, &shuffledDeck, 1)
			playerCard = dealCard(playerCard, &shuffledDeck, 1)

			// 根据规则决定是否发第三张牌
			if shouldPlayerDrawThirdCard(playerCard, bankerCard) { //闲家补牌
				playerCard = dealCard(playerCard, &shuffledDeck, 1)
			}

			if shouldBankerDrawThirdCard(playerCard, bankerCard) { //庄家补牌
				bankerCard = dealCard(bankerCard, &shuffledDeck, 1)
			}

			playerPoints := calculatePoints(playerCard)
			bankerPoints := calculatePoints(bankerCard)

			if playerPoints > bankerPoints {
				playerWin++
				roadSign = append(roadSign, "闲")
				if totalNums < 72 {
					playerWinNum[totalNums]++
				}
			} else if bankerPoints > playerPoints {
				bankerWin++
				roadSign = append(roadSign, "庄")
				if bankerPoints == 6 {
					luckySix++
				}
				if totalNums < 72 {
					bannerWinNum[totalNums]++
				}
			} else {
				if totalNums < 72 {
					drawNum[totalNums]++
				}
				draw++
				roadSign = append(roadSign, "和")
			}
			totalNums++

		}
		printRoadSign(&roadSign)
		fmt.Println("闲赢:", playerWin, "庄赢:", bankerWin, "和局:", draw, "总局数:", totalNums,
			"闲赢率:", fmt.Sprintf("%.2f", float64(playerWin)/float64(totalNums)*100),
			"庄赢率:", fmt.Sprintf("%.2f", float64(bankerWin)/float64(totalNums)*100),
			"和局率:", fmt.Sprintf("%.2f", float64(draw)/float64(totalNums)*100), "幸运6:", luckySix)
	}
	fmt.Println(bannerWinNum)
	fmt.Println(playerWinNum)
	fmt.Println(drawNum)

}

func printRoadSign(roadSign *[]string) {
	n := len(*roadSign)
	if n%6 > 0 {
		for i := 0; i < 6-(n%6); i++ {
			*roadSign = append(*roadSign, "")
		}
	}

	res := ""
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	circle := "\u25CF" // 使用黑色圆圈字符

	for i := 0; i < 6; i++ {
		for j := 0; j < len(*roadSign)/6; j++ {
			if (*roadSign)[i+j*6] == "庄" {
				res += red("\x1b[31m" + circle + "\x1b[0m" + " ")
			} else if (*roadSign)[i+j*6] == "闲" {
				res += blue("\x1b[94m" + circle + "\x1b[0m" + " ")
			} else if (*roadSign)[i+j*6] == "和" {
				res += green("\x1b[92m" + circle + "\x1b[0m" + " ")
			}
		}
		res += "\n"
	}
	fmt.Println(res)
}

func initializeDeck() []Card {
	suits := []string{"红心", "方块", "梅花", "黑桃"}
	values := []string{
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	} //8副牌

	var deck []Card

	for _, suit := range suits {
		for _, value := range values {
			card := Card{Suit: suit, Value: value}
			deck = append(deck, card)
		}
	}

	return deck
}

func initializeScoreMap(m map[string]int) {
	m["A"] = 1
	m["2"] = 2
	m["3"] = 3
	m["4"] = 4
	m["5"] = 5
	m["6"] = 6
	m["7"] = 7
	m["8"] = 8
	m["9"] = 9
	m["10"] = 10
	m["J"] = 10
	m["Q"] = 10
	m["K"] = 10
}

func shuffleDeck(deck []Card) []Card {
	rand.Seed(time.Now().UnixNano())

	// Fisher-Yates 洗牌算法
	for i := len(deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}
func dealCard(hand []Card, deck *[]Card, num int) []Card {
	card := (*deck)[:num]
	*deck = (*deck)[num:]
	hand = append(hand, card...)
	return hand
}

func removeDeckHead(deck []Card) []Card {
	return deck[11:]
}

func shouldPlayerDrawThirdCard(playerHand, bankerHand []Card) bool {
	playerPoints := calculatePoints(playerHand) //闲家点数
	bankerPoints := calculatePoints(bankerHand) //庄家点数

	if playerPoints >= 8 || bankerPoints >= 8 { //闲家或者庄家点数大于8
		return false //不能发牌
	}

	if playerPoints <= 5 { //闲家点数小于5，可以发牌
		return true
	}

	return false
}

func shouldBankerDrawThirdCard(playerHand, bankerHand []Card) bool {
	playerPoints := calculatePoints(playerHand) //闲家点数
	bankerPoints := calculatePoints(bankerHand) //庄家点数
	if bankerPoints >= 8 || playerPoints >= 8 { //庄家点数大于8,不能发牌
		return false
	}

	switch bankerPoints {
	case 0, 1, 2:
		return true
	case 3:
		if len(playerHand) == 3 && playerHand[2].Value != "8" {
			return true
		}
	case 4:
		if len(playerHand) == 3 && (playerHand[2].Value >= "2" && playerHand[2].Value <= "7") {
			return true
		}
	case 5:
		if len(playerHand) == 3 && (playerHand[2].Value >= "4" && playerHand[2].Value <= "7") {
			return true
		}
	case 6:
		if len(playerHand) == 3 && (playerHand[2].Value == "6" || playerHand[2].Value == "7") {
			return true
		}
	}

	return false
}

func displayHand(hand []Card) {
	for _, card := range hand {
		fmt.Printf("%s%s ", card.Suit, card.Value)
	}
	fmt.Println()
}

func calculatePoints(hand []Card) int {
	points := 0
	for _, card := range hand {
		points += m[card.Value]
	}
	return points % 10
}
