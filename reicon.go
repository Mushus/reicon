package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/Mushus/twtr"
	"github.com/disintegration/imaging"
)

// ImageSize is size of upload image to twitter
// see: https://dev.twitter.com/rest/reference/post/account/update_profile_image
const ImageSize = 400

var (
	regexpRGBA = regexp.MustCompile(`^\s*#([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})\s*$`)
	regexpRGB  = regexp.MustCompile(`^\s*#([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})\s*$`)
)

var option = twtr.AuthOption{
	ConsumerCreds: twtr.ConsumerCreds{
		Token:  "VG5ZZjwkMOQMK6jqClQ27fcVf",
		Secret: "6u83aMUYdPYdEUEF9OUQjXLmT1RJHQdLFIamsl6yeRvGEYQkCr",
	},
}

var (
	imageProcess  = flag.String("image", "", `image file list (separator character is ;)`)
	colorCodeList = flag.String("color", "#ffffff00", `color list (separator character is ",")`)
	configfile    = flag.String("config", "", "config file")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 150; i++ {
		rand.Intn(255)
	}

	flag.Parse()
	if configfile == nil {
		flag.Usage()
		return
	}

	isEdit := false
	config, isExist, err := getConfig(*configfile)
	if err != nil {
		log.Fatal("failed to open configuration file: ", err)
	}

	isEdit = isEdit || !isExist
	if isExist {
		var tw twtr.Twtr
		oldConfig := config
		tw, err = createTwitterClient(&config)
		isEdit = isEdit || config != oldConfig
		if err != nil {
			log.Fatal("failed to authorize twitter: ", err)
		}
		if *imageProcess == "" {
			flag.Usage()
		} else {
			err = rollingIcon(tw, *colorCodeList, *imageProcess)
			if err != nil {
				log.Print("failed to change icon: ", err)
			}
		}
	}

	if isEdit {
		err = saveConfig(*configfile, config)
		if err != nil {
			log.Fatal("failed to save configuration file", err)
		}
		log.Println("save config")
	}
}

type config struct {
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func getConfig(file string) (config, bool, error) {
	cfg := config{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, false, nil
		}
		return cfg, false, err
	}
	json.Unmarshal(b, &cfg)
	return cfg, true, nil
}

func saveConfig(file string, cfg config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode configuration file: %v", err)
	}
	err = ioutil.WriteFile(file, b, 0700)
	if err != nil {
		return fmt.Errorf("failed to store file: %v", err)
	}
	return nil
}

func createTwitterClient(cfg *config) (twtr.Twtr, error) {
	var tw twtr.Twtr

	if cfg.ConsumerKey != "" || cfg.ConsumerSecret != "" {
		option.ConsumerCreds.Token = cfg.ConsumerKey
		option.ConsumerCreds.Secret = cfg.ConsumerSecret
	}

	if cfg.AccessToken != "" || cfg.AccessTokenSecret != "" {
		accessCreds := twtr.AccessCreds{
			Token:  cfg.AccessToken,
			Secret: cfg.AccessTokenSecret,
		}
		tw = twtr.NewTwtr(option, accessCreds)
	} else {
		auth := twtr.NewAuth(option)
		info, err := auth.GenerateAuthorizationInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to generate authorization infomation: %v", err)
		}

		fmt.Printf("open this url: %s\n", info.AuthorizationURL)
		fmt.Print("pin? ")
		stdin := bufio.NewScanner(os.Stdin)
		if !stdin.Scan() {
			return nil, fmt.Errorf("canceled")
		}

		pin := stdin.Text()
		tw, err = auth.CreateClient(pin)
		if err != nil {
			return nil, fmt.Errorf("failed to get Twitter API: %v", err)
		}
		creds := tw.GetAccessCreds()
		cfg.AccessToken = creds.Token
		cfg.AccessTokenSecret = creds.Secret
	}
	return tw, nil
}

func rollingIcon(tw twtr.Twtr, colorList string, imageProcess string) error {
	baseImg := image.NewRGBA(image.Rect(0, 0, ImageSize, ImageSize))

	colorImg, err := createColor(colorList)
	if err != nil {
		return fmt.Errorf("failed to create color: %v", err)
	}

	draw.Draw(baseImg, baseImg.Bounds(), colorImg, image.ZP, draw.Src)

	var img image.Image
	processes := strings.Split(imageProcess, ";")
	for _, process := range processes {
		img, err = createImage(process)
		if err != nil {
			return fmt.Errorf("failed to load image: %v", err)
		}
		draw.Draw(baseImg, baseImg.Bounds(), img, image.ZP, draw.Over)
	}

	writer := &bytes.Buffer{}

	err = png.Encode(writer, baseImg)
	if err != nil {
		return fmt.Errorf("failed to encode profile icon: %v", err)
	}

	upiOpt := twtr.UpdateProfileImageOption{
		Image: writer.Bytes(),
	}

	_, err = tw.Account().UpdateProfileImage(upiOpt)
	if err != nil {
		return fmt.Errorf("failed to update profile icon: %v", err)
	}
	return nil
}

func createColor(colorsCodeList string) (image.Image, error) {
	colorCodes := strings.Split(colorsCodeList, ",")

	rndIdx := rand.Intn(len(colorCodes))
	colorCode := colorCodes[rndIdx]

	var clr color.Color
	if hexCode := regexpRGBA.FindStringSubmatch(colorCode); hexCode != nil {
		r, _ := strconv.ParseInt(hexCode[1], 16, 0)
		g, _ := strconv.ParseInt(hexCode[2], 16, 0)
		b, _ := strconv.ParseInt(hexCode[3], 16, 0)
		a, _ := strconv.ParseInt(hexCode[4], 16, 0)
		clr = color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		}
	} else if hexCode := regexpRGB.FindStringSubmatch(colorCode); hexCode != nil {
		r, _ := strconv.ParseInt(hexCode[1], 16, 0)
		g, _ := strconv.ParseInt(hexCode[2], 16, 0)
		b, _ := strconv.ParseInt(hexCode[3], 16, 0)
		clr = color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: 255,
		}
	} else {
		return nil, fmt.Errorf("cannot convert \"%s\" to color", colorCode)
	}

	return image.NewUniform(clr), nil
}

func createImage(imagename string) (image.Image, error) {
	rndFiles := make([]string, 0)
	matchedFile, err := filepath.Glob(imagename)
	if err != nil {
		return nil, fmt.Errorf("failed to listing files: %v", err)
	}

	for _, filename := range matchedFile {
		ext := filepath.Ext(filename)
		switch ext {
		case ".gif", ".png", ".jpeg", ".jpg":
			rndFiles = append(rndFiles, filename)
		}
	}

	if len(rndFiles) == 0 {
		return nil, fmt.Errorf("file not found: %s", imagename)
	}

	rndIdx := rand.Intn(len(rndFiles))

	decideFilename := rndFiles[rndIdx]
	reader, err := os.Open(decideFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %v", err)
	}

	fmt.Printf("using image:%s\n", decideFilename)

	return imaging.Fill(img, ImageSize, ImageSize, imaging.Center, imaging.Lanczos), nil
}
