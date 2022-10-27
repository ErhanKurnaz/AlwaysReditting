package main

import (
	"AlwaysRedditing/api"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/reujab/wallpaper"
)

func main() {
	a := app.New()
	window := a.NewWindow("AlwaysRedditing")
	window.Resize(fyne.NewSize(500, 500))

	window.SetContent(App())

	window.ShowAndRun()
}

func App() *fyne.Container {
	return container.NewVBox(
		SearchBox(),
	)
}

func SearchBox() *fyne.Container {
	searchValue := binding.NewString()
	searchInput := widget.NewSelectEntry([]string{})
	searchInput.Bind(searchValue)
	searchInput.SetPlaceHolder("Subreddit name")

	show18 := binding.NewBool()
	var debounceTimer *time.Timer = nil

	searchInput.OnChanged = func(newInput string) {
		if debounceTimer != nil {
			debounceTimer.Stop()
		}

		fetchAutoComplete := func() {
			over18, _ := show18.Get()
			redditOptions, err := api.RedditAutoComplete(newInput, over18)
			if err != nil {
				fmt.Printf("Error fetching autocomplete: %s", err)
				return
			}

			newOptions := make([]string, len(redditOptions.Children))

			for i, option := range redditOptions.Children {
				newOptions[i] = option.Data.DisplayName
			}

			searchInput.SetOptions(newOptions)
			debounceTimer = nil
		}

		debounceTimer = time.AfterFunc(time.Duration(2)*time.Second, fetchAutoComplete)
	}

	selectButton := widget.NewButton("Select", func() {
		fmt.Println("selected")
		background, _ := wallpaper.Get()
		fmt.Println(background)

		over18, _ := show18.Get()
		posts, err := api.RedditGetPost(searchInput.Text, over18)
		if err != nil {
			fmt.Printf("Error fetching posts: %s", err)
			return
		}

		var selectedImage *string = nil
		for _, post := range posts.Children {
			if !post.Data.IsVideo && post.Data.PostHint == "image" && (!post.Data.Over18 || over18) {
				selectedImage = &post.Data.Url

				break
			}
		}

		if selectedImage == nil {
			fmt.Println("No suitable image found")
			return
		}

		wallpaper.SetFromURL(*selectedImage)
		// fileName, err := api.RedditSaveImage(*selectedImage)
		// if err != nil {
		// fmt.Printf("Could not save image: %s", err)
		// return
		// }
		//
		// fmt.Println(fileName)
	})

	optionContainer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewCheckWithData("show 18+", show18),
		layout.NewSpacer(),
		selectButton,
	)

	return container.New(
		layout.NewGridLayout(2),
		searchInput,
		optionContainer,
	)
}
