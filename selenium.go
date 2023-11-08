package go_selenium

import (
	"bytes"
	"fmt"
	go_file "github.com/pefish/go-file"
	go_time "github.com/pefish/go-time"
	"github.com/tebeka/selenium"
	"image"
	"image/png"
	"os"
)

func GetChild(ele selenium.WebElement, index int) (selenium.WebElement, error) {
	children, err := ele.FindElements(selenium.ByXPATH, "*")
	if err != nil {
		return nil, err
	}
	return children[index], nil
}

func ScrollToBottom(wd selenium.WebDriver) error {
	_, err := wd.ExecuteScript(`window.scrollTo(0, document.body.scrollHeight)`, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetChildByPath(ele selenium.WebElement, indexs []int) (selenium.WebElement, error) {
	result := ele
	for _, index := range indexs {
		children, err := result.FindElements(selenium.ByXPATH, "*")
		if err != nil {
			return nil, err
		}
		if index < 0 {
			index += len(children)
		}
		result = children[index]
	}

	return result, nil
}

func SaveSnapshot(wd selenium.WebDriver, name string) error {
	pngBytes, err := wd.Screenshot()
	if err != nil {
		return err
	}
	img, _, err := image.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		return err
	}
	dir := fmt.Sprintf("./snapshot_%s", name)
	err = go_file.FileInstance.AssertPathExist(dir)
	if err != nil {
		return err
	}
	out, err := os.Create(fmt.Sprintf(
		"%s/%d_%s.png",
		dir,
		go_time.TimeInstance.CurrentTimestamp(go_time.TimeUnit_SECOND),
		name,
	))
	if err != nil {
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
