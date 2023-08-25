package form

type Button struct {
	// Text holds the text displayed on the button. It may use Minecraft formatting codes and may have
	// newlines.
	Text string `json:"text"`
	// Image holds a path to an image for the button. The Image may either be a URL pointing to an image,
	// such as 'https://someimagewebsite.com/someimage.png', or a path pointing to a local asset, such as
	// 'textures/block/grass_carried'.
	Image any `json:"image"`
}

type ImageType struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
