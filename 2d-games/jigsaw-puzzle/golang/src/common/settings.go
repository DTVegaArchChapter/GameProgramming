package common

import "image/color"

var (
	ScreenWidth  = 1280
	ScreenHeight = 720

	// Dark neutral gray for game background
	BackgroundColor = color.RGBA{R: 36, G: 36, B: 36, A: 255} // #242424

	// Header color - dark gray
	HeaderColor = color.RGBA{R: 20, G: 20, B: 20, A: 255} // #141414

	// Title - warm golden yellow (stands out without being too aggressive)
	TitleColor = color.RGBA{R: 255, G: 215, B: 0, A: 255} // #FFD700

	// General text - near-white for better contrast on dark bg
	BodyTextColor = color.RGBA{R: 245, G: 245, B: 245, A: 255} // #F5F5F5

	// Footer color - dark gray
	FooterColor = color.RGBA{R: 20, G: 20, B: 20, A: 255} // #141414

	// Header button color - slightly lighter gray for visibility
	HeaderButtonColor = color.RGBA{R: 40, G: 40, B: 40, A: 255} // #282828

	// Header button hover color - lighter gray for interaction feedback
	HeaderButtonHoverColor = color.RGBA{R: 56, G: 56, B: 56, A: 255} // #383838
)
