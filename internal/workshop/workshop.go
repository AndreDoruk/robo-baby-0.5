package workshop

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const workshop_item_url string = "https://steamcommunity.com/sharedfiles/filedetails/?id="
const weekly_most_popular_url string = "https://steamcommunity.com/workshop/browse/?appid=250900&browsesort=trend&section=readytouseitems&days=7&actualsort=trend&p=1"
const random_item_url_prefix string = "https://steamcommunity.com/workshop/browse/?appid=250900&browsesort=trend&section=readytouseitems&actualsort=trend&days=7&p=3"

const weekly_item_num int = 9
const max_page_num int = 50

const selection_item_num = 31

type WorkshopItem struct {
	URL         string
	Name        string
	Icon        string
	Visitors    int
	Subscribers int
	Favorites   int
	//description string
}

type WorkshopComment struct {
	Creator string
	Comment string
	IconURL string
}

const fallback_image = "https://cdn.discordapp.com/attachments/1146773305693585420/1151462507706331158/image.png"

func getItem(url string) WorkshopItem {
	document := getPageDocument(url)

	itemName := document.Find(".workshopItemTitle").Text()
	previewImage := document.Find(".workshopItemPreviewImageMain").Last().AttrOr("src", fallback_image)

	if previewImage == fallback_image {
		previewImage = document.Find(".workshopItemPreviewImageEnlargeable").Last().AttrOr("src", fallback_image)
	}

	tbody := document.Find("tbody") // Only tbody in the html lol!
	itemStats := getStatsFromChildren(tbody)

	visitorNum := itemStats[0]
	subscriberNum := itemStats[1]
	favoriteNum := itemStats[2]

	return WorkshopItem{url, itemName, previewImage, visitorNum, subscriberNum, favoriteNum}
}

func getPageDocument(url string) *goquery.Document {
	response, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return document
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
	selection := getPageDocument(url).Find(".workshopItem")

	items := make([]WorkshopItem, numItems)
	var waitGroup sync.WaitGroup

	selection.Each(func(index int, selection *goquery.Selection) {
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

func GetRandomItem() WorkshopItem {
	for {
		url := random_item_url_prefix + strconv.Itoa(rand.Intn(max_page_num))
		selection := getPageDocument(url).Find(".workshopItem")

		var workshopItem WorkshopItem

		selectionItems := getSelectionItems(selection)
		if len(selectionItems) == 0 {
			fmt.Println("god is dead")
			continue
		}

		selectedItem := selectionItems[rand.Intn(len(selectionItems))]

		hrefLink := selectedItem.Find("a").AttrOr("href", "https://dontasktoask.com")
		workshopItem = getItem(hrefLink)

		return workshopItem
	}
}

func getSelectionItems(selection *goquery.Selection) []*goquery.Selection {
	items := make([]*goquery.Selection, 0, selection_item_num)

	selection.Each(func(_ int, selection *goquery.Selection) {
		items = append(items, selection)
	})

	return items
}

func GetRandomCommentAndItem() (WorkshopComment, WorkshopItem) {
	for {
		item := GetRandomItem()
		comment := getRandomCommentFromItem(item)

		if comment != nil {
			return *comment, item
		}
	}
}

func getRandomCommentFromItem(item WorkshopItem) *WorkshopComment {
	selection := getPageDocument(item.URL).Find(".commentthread_comment")
	selectionItems := getSelectionItems(selection)

	if len(selectionItems) == 0 {
		return nil
	}

	selectedComment := selectionItems[rand.Intn(len(selectionItems))]

	return &WorkshopComment{
		Creator: selectedComment.Find(".commentthread_author_link").Find("bdi").Text(),
		Comment: selectedComment.Find(".commentthread_comment_text").Text(),
		IconURL: selectedComment.Find(".commentthread_comment_avatar").Find("img").AttrOr("src", "https://tenor.com/view/epico-mandela-catalog-mandela-catalogue-intruder-creepy-gif-23565190"),
	}
}
