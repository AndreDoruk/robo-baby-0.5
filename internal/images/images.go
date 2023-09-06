package images

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/trustig/robobaby0.5/internal/workshop"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

const font_face_path string = "./src/FontSouls.ttf"

// Mod
const mod_background_image_path string = "./src/modBackground.png"
const tomato_image_path string = "./src/tomato.png"

const thumbnail_position_x int = 51
const thumbnail_position_y int = 51

const scale_mul float64 = 0.92

const letters_before_scale int = 12
const mod_name_size float64 = 90

const mod_name_position_x float64 = 401
const mod_name_position_y float64 = 157

const stat_size float64 = 59

const visitor_prefix string = "Visitors: " // Indeed, I have slept long enough
const visitor_position_x float64 = 409
const visitor_position_y float64 = 140 + 81

const subscriber_prefix string = "Subscribers: "
const subscriber_position_x float64 = 409
const subscriber_position_y float64 = 180 + 81

const favorite_prefix string = "Favorites: "
const favorite_position_x float64 = 409
const favorite_position_y float64 = 220 + 81

const adoption_suffix string = "% Adoption rate"
const adoption_size float64 = 40
const adoption_position_x float64 = 409
const adoption_position_y float64 = 260 + 70

// Votes
const vote_background_image_path string = "./src/voteBackground.png"
const vote_time_overlay_path string = "./src/voteTimeOverlay.png"

const user_avatar_x int = 34
const user_avatar_y int = 16

const join_vote_text string = "Join Vote"
const join_vote_size float64 = 30
const join_vote_x float64 = 184
const join_vote_y float64 = 74

const username_size float64 = 50
const username_x float64 = 184
const username_y float64 = 104

const time_prefix string = "Ends in "
const time_size float64 = 30
const time_x float64 = 184
const time_y float64 = 124

const hour_suffix string = " hour(s)"
const minute_suffix string = " minute(s)"

const victory_text string = "Vote succeeded!"
const loss_text string = "Vote failed!"

type color struct {
	r float64
	g float64
	b float64
}

var WHITE_COLOR color = color{1, 1, 1}
var YELLOW_COLOR color = color{1, 0.89, 0.12}
var GREEN_COLOR color = color{0.25, 0.75, 0.4}
var RED_COLOR color = color{0.75, 0.25, 0.25}

func CreateWorkshopImage(item workshop.WorkshopItem) image.Image {
	context := getContextFromPath(mod_background_image_path)

	context.DrawImage(getItemThumbnail(item), thumbnail_position_x, thumbnail_position_y)

	drawText(context, item.Name, mod_name_position_x, mod_name_position_y, getScale(item.Name, mod_name_size), WHITE_COLOR)
	drawText(context, visitor_prefix+strconv.Itoa(item.Visitors), visitor_position_x, visitor_position_y, stat_size, WHITE_COLOR)
	drawText(context, subscriber_prefix+strconv.Itoa(item.Subscribers), subscriber_position_x, subscriber_position_y, stat_size, WHITE_COLOR)
	drawText(context, favorite_prefix+strconv.Itoa(item.Favorites), favorite_position_x, favorite_position_y, stat_size, WHITE_COLOR)

	adoptionRate := item.Subscribers * 100 / item.Visitors
	drawText(context, strconv.Itoa(adoptionRate)+adoption_suffix, adoption_position_x, adoption_position_y, adoption_size, YELLOW_COLOR)

	return context.Image()
}

func getScale(text string, startingScale float64) float64 {
	offsetLetters := len(text) - letters_before_scale

	if offsetLetters < 0 {
		return startingScale
	}

	return startingScale * math.Pow(scale_mul, float64(offsetLetters))
}

func CreateVoteImage(session *discordgo.Session, user *discordgo.User) image.Image {
	userAvatar, err := session.UserAvatarDecode(user)

	if err != nil {
		log.Fatalln(err)
	}

	context := getContextFromPath(vote_background_image_path)
	context.DrawImage(userAvatar, user_avatar_x, user_avatar_y)

	drawText(context, join_vote_text, join_vote_x, join_vote_y, join_vote_size, WHITE_COLOR)
	drawText(context, user.Username, username_x, username_y, getScale(user.Username, username_size), WHITE_COLOR)
	drawText(context, time_prefix+"12"+hour_suffix, time_x, time_y, time_size, WHITE_COLOR)

	return context.Image()
}

func UpdateVoteTimestamp(image image.Image, minutesRemaining float64) image.Image {
	context := GetContextFromImage(image)
	timeOverlay, err := gg.LoadImage(vote_time_overlay_path)

	if err != nil {
		log.Fatalln(err)
	}

	context.DrawImage(timeOverlay, 0, 0)

	if minutesRemaining > 60 {
		hours := strconv.FormatFloat(math.Floor(minutesRemaining/60), 'f', -1, 64)
		drawText(context, time_prefix+hours+hour_suffix, time_x, time_y, time_size, WHITE_COLOR)
	} else if minutesRemaining > 0 {
		minutes := strconv.FormatFloat(minutesRemaining, 'f', -1, 64)
		drawText(context, time_prefix+minutes+minute_suffix, time_x, time_y, time_size, WHITE_COLOR)
	} else {
		drawText(context, "Vote ended!", time_x, time_y, time_size, WHITE_COLOR)
	}

	return context.Image()
}

func UpdateVoteVictoryText(image image.Image, win bool) image.Image {
	context := GetContextFromImage(image)
	timeOverlay, err := gg.LoadImage(vote_time_overlay_path)

	if err != nil {
		log.Fatalln(err)
	}

	context.DrawImage(timeOverlay, 0, 0)

	if win {
		drawText(context, victory_text, time_x, time_y, time_size, GREEN_COLOR)
	} else {
		drawText(context, loss_text, time_x, time_y, time_size, RED_COLOR)
	}

	return context.Image()
}

func GetContextFromImage(backgroundImage image.Image) *gg.Context {
	backgroundBounds := backgroundImage.Bounds()
	width := backgroundBounds.Dx()
	height := backgroundBounds.Dy()

	context := gg.NewContext(width, height)
	context.DrawImage(backgroundImage, 0, 0)

	return context
}

func getContextFromPath(path string) *gg.Context {
	backgroundImage, err := gg.LoadImage(path)

	if err != nil {
		log.Fatalln(err)
	}

	return GetContextFromImage(backgroundImage)
}

func drawText(context *gg.Context, text string, textX float64, textY float64, size float64, color color) {
	fontface, err := gg.LoadFontFace(font_face_path, size)

	if err != nil {
		log.Fatalln(err)
	}

	context.SetFontFace(fontface)
	context.SetRGB(color.r, color.g, color.b)

	context.DrawStringAnchored(text, textX, textY, 0, 0)
}

func getItemThumbnail(item workshop.WorkshopItem) image.Image {
	return GetImageFromUrl(item.Icon)
}

func GetImageFromUrl(url string) image.Image {
	response, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	thumbnailBytes, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	byteReader := bytes.NewReader(thumbnailBytes)
	image, _, err := image.Decode(byteReader)

	if err != nil {
		image, err = gif.Decode(byteReader)

		if err != nil {
			log.Fatalln(err)
		}
	}

	return image
}

func Tomato(originalImage image.Image) image.Image {
	bounds := originalImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	context := gg.NewContext(width, height)
	context.DrawImage(originalImage, 0, 0)

	tomatoImage, err := gg.LoadImage(tomato_image_path)

	if err != nil {
		log.Fatalln(err)
	}

	context.DrawImage(tomatoImage, rand.Intn(width), rand.Intn(height))

	return context.Image()
}

func ImageToBytesReader(image image.Image) *bytes.Reader {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, image)

	if err != nil {
		log.Fatalln(err)
	}

	reader := bytes.NewReader(buffer.Bytes())

	if err != nil {
		log.Fatalln(err)
	}

	return reader
}

func ImageToRaw(image image.Image) string {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, image)

	if err != nil {
		log.Fatalln(err)
	}

	return buffer.String()
}
