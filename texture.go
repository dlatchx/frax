package main

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"errors"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var lastActiveTexture uint32 = 0

type TextureManager struct {
	textures map[string]*Texture
}

func NewTextureManager() *TextureManager {
	return &TextureManager{make(map[string]*Texture)}
}

func (tm *TextureManager) Get(name string) (tex *Texture, wasLoaded bool) {
	tex, wasLoaded = tm.textures[name]

	if !wasLoaded {
		tex = NewTexture()
		tm.textures[name] = tex
	}

	return
}

func (tm *TextureManager) MustGet(name string) *Texture {
	tex, found := tm.textures[name]

	if !found {
		log.Panicf(`Texture "%s" was not found in memory`, name)
	}

	return tex
}

func (tm *TextureManager) Has(name string) bool {
	_, found := tm.textures[name]
	return found
}

func (tm *TextureManager) Load(name, path string) (*Texture, error) {
	tex, wasLoaded := tm.Get(name)

	if !wasLoaded {
		log.Debugf(`Loading texture "%s" <- "%s"`, name, path)
		err := tex.SetImgFile(path)
		if err != nil {
			tm.Unload(name)
			return nil, err
		}
	}

	return tex, nil
}

func (tm *TextureManager) MustLoad(name, path string) *Texture {
	tex, err := tm.Load(name, path)

	if err != nil {
		log.Panicf(`Could not load texture "%s" ("%s")`, path, name)
	}

	return tex
}

func (tm *TextureManager) Reload(name, path string) (*Texture, error) {
	log.Debugf(`Reloading texture "%s" <- "%s"`, name, path)

	tex, _ := tm.Get(name)
	err := tex.SetImgFile(path)

	return tex, err
}

func (tm *TextureManager) Take(name string) *Texture {
	tex, found := tm.textures[name]

	if !found {
		return nil
	}

	delete(tm.textures, name)

	return tex
}

func (tm *TextureManager) MustTake(name string) *Texture {
	tex := tm.Take(name)

	if tex == nil {
		log.Panicf(`Texture "%s" was not found in memory`, name)
	}

	return tex
}

func (tm *TextureManager) Unload(name string) {
	log.Debugf(`Unloading texture "%s"`, name)

	tex, found := tm.textures[name]

	if found {
		tex.Delete()
		delete(tm.textures, name)
	}
}

type Texture struct {
	id       uint32
	lastSize image.Point
}

func NewTexture() *Texture {
	gl.Enable(gl.TEXTURE_2D)

	var texId uint32
	gl.GenTextures(1, &texId)

	texture := &Texture{texId, image.Pt(0, 0)}
	runtime.SetFinalizer(texture, func(t *Texture){
		t.Delete()
	})

	return texture
}

func NewTextureImg(img image.Image) (*Texture, error) {
	tex := NewTexture()

	err := tex.SetImg(img)

	return tex, err
}

func NewTextureImgFile(path string) (*Texture, error) {
	tex := NewTexture()

	err := tex.SetImgFile(path)

	return tex, err
}

func (t *Texture) ActiveTexture(i uint32) {
	if i < 0 {
		i = 0
	} else if i > MaxTextureUnits() - 1 {
		i = MaxTextureUnits() - 1
	}

	gl.ActiveTexture(gl.TEXTURE0 + i)
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t *Texture) Uniform(sp *ShaderProgram, name string) {
	uniformLocation := sp.GetUniformLocation(name)

	lastActiveTexture++;
	if lastActiveTexture > MaxTextureUnits() - 1 {
		lastActiveTexture = 0
	}

	t.ActiveTexture(lastActiveTexture)
	gl.Uniform1i(uniformLocation, int32(lastActiveTexture))
}

func (t *Texture) Bind() {
	if t.id != 0 {
		gl.BindTexture(gl.TEXTURE_2D, t.id)
	}
}

func (t *Texture) Delete() {
	if t.id != 0 {
		gl.DeleteTextures(1, &t.id)
		t.id = 0
		t.lastSize = image.Pt(0, 0)
	}
}

func (t Texture) SetImg(img image.Image) error {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X * 4 {
		return errors.New("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.BindTexture(gl.TEXTURE_2D, t.id)

	if t.lastSize == rgba.Rect.Size() {
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			int32(0),
			int32(0),
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)

		t.lastSize = rgba.Rect.Size()
	}

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return nil
}

func (t Texture) SetImgFile(path string) error {
	imgfile, err := os.Open(path)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(imgfile)
	if err != nil {
		return err
	}

	return t.SetImg(img)
}


func MaxTextureUnits() uint32 {
	var texUnits int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, &texUnits)

	if texUnits < 16 {
		log.Warning("GPU only supporting %v texture units (should be 16 min.)", texUnits)
	}

	return uint32(texUnits)
}
