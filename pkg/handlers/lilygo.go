package handlers

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/chux0519/yeti/pkg/service/rank"
	"github.com/chux0519/yeti/pkg/utils"
)

func LilygoRankHandler(R *rank.YetiRankService, ign string) (*rank.RankData, error) {
	data, err := R.FetchUserRank(ign)
	if err != nil {
		return nil, err
	}

	// avatarBytes, _ := rank.HttpGet(data.CharacterImageURL)

	// avartarImg, err := png.Decode(bytes.NewBuffer(avatarBytes))
	// if err != nil {
	// 	return nil, err
	// }

	profileImg, err := data.GetProfileImageGo()
	if err != nil {
		return nil, err
	}

	raw := utils.ImgToRaw(profileImg, 0xD69A)
	pixels := []string{}
	for _, pix := range raw {
		pixHex := fmt.Sprintf("0x%X", pix)
		pixels = append(pixels, pixHex)
	}
	hexData := strings.Join(pixels, ",")

	header := fmt.Sprintf("int avatar_w = %d;int avatar_h = %d;\n", profileImg.Bounds().Dx(), profileImg.Bounds().Dy())
	header += fmt.Sprintf("const unsigned short avatar[%d] PROGMEM={%s};", len(raw), hexData)

	// println(header)

	if err := os.WriteFile(fmt.Sprintf("%s.h", ign), []byte(header), 0666); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer([]byte{})

	err = png.Encode(buffer, profileImg)
	if err != nil {
		return nil, err
	}
	profileBytes, _ := ioutil.ReadAll(buffer)
	if err := R.SaveProfileImageToCache(R.GetProfileImageName(data), profileBytes); err != nil {
		return nil, err
	}

	return data, nil
}
