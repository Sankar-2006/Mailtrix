// Copyright 2023 Krisna Pranav, Sankar-2006. All rights reserved.
// Use of this source code is governed by a Apache-2.0 License
// license that can be found in the LICENSE file

package apiv1

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/jhillyerd/enmime"
	"github.com/krishpranav/Mailtrix/storage"
	"github.com/krishpranav/Mailtrix/utils/logger"
)

var (
	thumbWidth  = 180
	thumbHeight = 120
)

func Thumbnail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	partID := vars["partID"]

	a, err := storage.GetAttachmentPart(id, partID)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	fileName := a.FileName
	if fileName == "" {
		fileName = a.ContentID
	}

	if !strings.HasPrefix(a.ContentType, "image/") {
		blankImage(a, w)
		return
	}

	buf := bytes.NewBuffer(a.Content)

	img, err := imaging.Decode(buf)
	if err != nil {
		logger.Log().Warning(err)
		blankImage(a, w)
		return
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	var dstImageFill *image.NRGBA

	if img.Bounds().Dx() < thumbWidth || img.Bounds().Dy() < thumbHeight {
		dstImageFill = imaging.Fit(img, thumbWidth, thumbHeight, imaging.Lanczos)
	} else {
		dstImageFill = imaging.Fill(img, thumbWidth, thumbHeight, imaging.Center, imaging.Lanczos)
	}
	dst := imaging.New(thumbWidth, thumbHeight, color.White)
	dst = imaging.OverlayCenter(dst, dstImageFill, 1.0)

	if err := jpeg.Encode(foo, dst, &jpeg.Options{Quality: 70}); err != nil {
		logger.Log().Warning(err)
		blankImage(a, w)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "filename=\""+fileName+"\"")
	_, _ = w.Write(b.Bytes())
}

func blankImage(a *enmime.Part, w http.ResponseWriter) {
	rect := image.Rect(0, 0, thumbWidth, thumbHeight)
	img := image.NewRGBA(rect)
	background := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	dstImageFill := imaging.Fill(img, thumbWidth, thumbHeight, imaging.Center, imaging.Lanczos)

	if err := jpeg.Encode(foo, dstImageFill, &jpeg.Options{Quality: 70}); err != nil {
		logger.Log().Warning(err)
	}

	fileName := a.FileName
	if fileName == "" {
		fileName = a.ContentID
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "filename=\""+fileName+"\"")
	_, _ = w.Write(b.Bytes())
}
