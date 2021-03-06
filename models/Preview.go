package models

import "strings"

//PreviewMimes mimes assigned to preview
var PreviewMimes map[PreviewType][]string = map[PreviewType][]string{
	ImagePreviewType: []string{
		"image/*",
	},
	TextPreviewType: []string{
		"text/*",
	},
}

//PreviewType type of preview
type PreviewType uint8

//Preview types
const (
	DefaultPreviewType PreviewType = iota
	ImagePreviewType
	TextPreviewType
)

//PreviewTemplate template struct for preview
type PreviewTemplate struct {
	Filename       string
	PublicFilename string
	PreviewType    PreviewType
	Host           string
	FileSizeStr    string
	Encrypted      bool
}

//PreviewTypeFromMime get Type to preview from mime
func PreviewTypeFromMime(sMime string) PreviewType {
	if len(strings.TrimSpace(sMime)) == 0 {
		return DefaultPreviewType
	}

	for ptype, mimes := range PreviewMimes {
		for _, mime := range mimes {
			if strings.HasSuffix(mime, "*") {
				if strings.HasPrefix(sMime, mime[:len(mime)-1]) {
					return ptype
				}
			} else {
				if mime == sMime {
					return ptype
				}
			}
		}
	}

	return DefaultPreviewType
}

//IsImagePreview return true if pt is image previewtype
func IsImagePreview(pt PreviewType) bool {
	return pt == ImagePreviewType
}

//IsTextPreview return true if pt is text previewtype
func IsTextPreview(pt PreviewType) bool {
	return pt == TextPreviewType
}

//IsDefaultPreview return true if pt is default previewtype
func IsDefaultPreview(pt PreviewType) bool {
	return pt == DefaultPreviewType
}
