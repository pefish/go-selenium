package go_selenium

import (
	"bytes"
	"fmt"
	go_file "github.com/pefish/go-file"
	go_time "github.com/pefish/go-time"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"image"
	"image/png"
	"os"
	"time"
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

func WaitDocumentReady(wd selenium.WebDriver, moreCondition func(wd selenium.WebDriver) (bool, error)) error {
	err := wd.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		result, err := wd.ExecuteScript("return document.readyState", nil)
		if err != nil {
			return false, err
		}
		if result.(string) != "complete" {
			return false, nil
		}
		if moreCondition != nil {
			return moreCondition(wd)
		}
		return true, nil
	}, 10*time.Second, time.Second)
	if err != nil {
		return err
	}
	return nil
}

func GetChildByPath(ele selenium.WebElement, indexes []int) (selenium.WebElement, error) {
	result := ele
	for i, index := range indexes {
		children, err := result.FindElements(selenium.ByXPATH, "*")
		if err != nil {
			return nil, err
		}
		if index < 0 {
			index += len(children)
		}
		if index >= len(children) {
			return nil, errors.Errorf("Not found on index %d", i)
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
