package workshop

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	//"github.com/boltdb/bolt"
)

const workshop_item_url string = "https://steamcommunity.com/sharedfiles/filedetails/?id="
const weekly_most_popular_url string = "https://steamcommunity.com/workshop/browse/?appid=250900&browsesort=trend&section=readytouseitems&days=7&actualsort=trend&p=1"

const weekly_item_num int = 9

type WorkshopItem struct {
	Name        string
	Icon        string
	Visitors    int
	Subscribers int
	Favorites   int
	//description string
}

func getItem(url string) WorkshopItem {
	response, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	itemName := document.Find(".workshopItemTitle").Text()
	previewImage := document.Find(".workshopItemPreviewImageMain").Last().AttrOr("src", "this world is cruel and painful")

	tbody := document.Find("tbody") // Only tbody in the html lol!
	itemStats := getStatsFromChildren(tbody)

	visitorNum := itemStats[0]
	subscriberNum := itemStats[1]
	favoriteNum := itemStats[2]

	return WorkshopItem{itemName, previewImage, visitorNum, subscriberNum, favoriteNum}
}

func GetItemFromId(id int) WorkshopItem {
	url := workshop_item_url + strconv.Itoa(id)
	return getItem(url)
}

func getStatsFromChildren(selection *goquery.Selection) []int {
	numbers := make([]int, 3)

	selection.Children().Each(func(index int, selection *goquery.Selection) {
		statSelection := selection.Children().First()

		statText := strings.Replace(statSelection.Text(), ",", "", -1)
		stat, err := strconv.Atoi(statText)

		if err != nil {
			log.Fatal(err)
		}

		numbers[index] = stat
	})

	return numbers
}

func GetMostPopularItems() []WorkshopItem {
	return getItemsFromSearch(weekly_most_popular_url, weekly_item_num)
}

func getItemsFromSearch(url string, numItems int) []WorkshopItem {
	response, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	items := make([]WorkshopItem, numItems)
	var waitGroup sync.WaitGroup

	document.Find(".workshopItem").Each(func(index int, selection *goquery.Selection) {
		hrefLink := selection.Find("a").AttrOr("href", "https://dontasktoask.com") //lol!

		if index < weekly_item_num {
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()
				items[index] = getItem(hrefLink)
			}()
		}
	})

	waitGroup.Wait()
	return items
}
