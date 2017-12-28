package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/guitmz/go-weatherbit"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func aboutDialog() *gtk.AboutDialog {
	aboutPB := gdkpixbuf.NewPixbufFromData(aboutLogo)
	ad := gtk.NewAboutDialog()
	ad.SetProgramName("Go Weather Indicator")
	ad.SetVersion("0.1")
	ad.SetCopyright("(c) Guilherme Thomazi Bonicontro")
	ad.SetComments("GTK Weather indicator application written in Go and using Weatherbit.io as provider.")
	ad.SetWebsite("https://www.guitmz.com")
	ad.SetLogo(aboutPB)
	return ad
}

func moreDialog(country, city, apiKey string) *gtk.Dialog {
	weatherResults, err := weatherbit.GetWeather(country, city, apiKey)
	check(err)
	weather := weatherResults.Data[0] // this API supports multiple datas but here I only want the current one.

	md := gtk.NewDialog()
	md.AddButton("Cool!", gtk.RESPONSE_OK)
	md.SetTitle("Detailed weather information")
	md.SetDefaultSize(300, 150)

	labelSunrise := gtk.NewLabel(fmt.Sprintf("Sunrise: %s", weather.Sunrise))
	labelSunset := gtk.NewLabel(fmt.Sprintf("Sunset: %s", weather.Sunset))
	labelDev := gtk.NewLabel("More to come soon :)")

	box := md.GetVBox()
	box.Add(labelSunrise)
	box.Add(labelSunset)
	box.Add(labelDev)

	md.ShowAll()
	return md
}

func buildIconURL(iconCode string) string {
	baseURL := "https://www.weatherbit.io/static/img/icons/%s.png"
	return fmt.Sprintf(baseURL, iconCode)
}

func updateStatusIcon(si *gtk.StatusIcon, country, city, apiKey string) bool {
	weatherResults, err := weatherbit.GetWeather(country, city, apiKey)
	check(err)
	weather := weatherResults.Data[0] // this API supports multiple datas but here I only want the current one.

	tmpfile, _ := ioutil.TempFile(os.TempDir(), "weather")
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	iconURL := buildIconURL(weather.Weather.Icon)
	resp, _ := http.Get(iconURL)
	_, _ = io.Copy(tmpfile, resp.Body)
	si.SetFromFile(tmpfile.Name())

	si.SetTooltipMarkup(fmt.Sprintf("%s\n%.1f°C\nFeels like %.1f°C\n%s", weather.CityName, weather.Temp, weather.AppTemp, weather.Weather.Description))
	return true
}

func main() {
	countryPtr := flag.String("country", "Germany", "Name of your country.")
	cityPtr := flag.String("city", "Berlin", "Name of your city.")
	apiPtr := flag.String("key", "", "Your wunderground api key")
	flag.Parse()

	gtk.Init(&os.Args)
	glib.SetApplicationName("go-gtk-weather-indicator")

	miQuit := gtk.NewMenuItemWithLabel("Quit")
	miQuit.Connect("activate", func() {
		gtk.MainQuit()
	})
	miAbout := gtk.NewMenuItemWithLabel("About")
	miAbout.Connect("activate", func() {
		ad := aboutDialog()
		ad.Run()
		ad.Destroy()
	})
	miMore := gtk.NewMenuItemWithLabel("Weather details")
	miMore.Connect("activate", func() {
		md := moreDialog(*countryPtr, *cityPtr, *apiPtr)
		md.Run()
		md.Destroy()
	})
	nm := gtk.NewMenu()
	nm.Append(miMore)
	nm.Append(miAbout)
	nm.Append(miQuit)
	nm.ShowAll()

	si := gtk.NewStatusIcon()
	si.SetTitle("Go Weather Indicator")
	si.Connect("popup-menu", func(cbx *glib.CallbackContext) {
		nm.Popup(nil, nil, gtk.StatusIconPositionMenu, si, uint(cbx.Args(0)), uint32(cbx.Args(1)))
	})

	glib.TimeoutAdd(3600000, updateStatusIcon(si, *countryPtr, *cityPtr, *apiPtr))

	gtk.Main()
	glib.NewMainLoop(nil, false).Run()
}
