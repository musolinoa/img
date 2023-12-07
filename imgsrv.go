package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

type YearIndexHandler struct {
	Idx *YearIdx
	Tpl *template.Template
}

func (h *YearIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "montage.jpg" {
		http.ServeFile(w, r, fmt.Sprintf("%s/montage.jpg", h.Idx.Path))
		return
	}

	type MonthTplData struct {
		Empty bool
		Number string
		Name string
	}

	type TplData struct {
		Title string
		Prev, Next string
		Curr string
		Months [12]MonthTplData
	}
	tplData := TplData{
		Title: fmt.Sprintf("Photos :: %s", path.Base(h.Idx.Path)),
		Next: h.Idx.Next(),
		Prev: h.Idx.Prev(),
	}
	for i := 0; i < 12; i++ {
		if h.Idx.Months[i] != nil {
			tplData.Months[i].Empty = false
			tplData.Months[i].Number = fmt.Sprintf("%02d", i+1)
			tplData.Months[i].Name = time.Month(i+1).String()
		} else {
			tplData.Months[i].Empty = true
		}
	}
	if err := h.Tpl.Execute(w, tplData); err != nil {
		log.Printf("error executing template: %v\n", err)
	}
}

type AlbumIndexHandler struct {
	Idx *AlbumIdx
	IndexTpl *template.Template
	ImageTpl *template.Template
	Tags *Tags
}

func (h *AlbumIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "":
		fallthrough
	case "index.html":
		type TplData struct {
			Title string
			Prev, Next string
			Images []string
		}
		tplData := TplData{
			Title: path.Base(h.Idx.Path),
			Images: h.Idx.Images,
		}
		if h.Idx.Year != 0 {
			tplData.Title = fmt.Sprintf("%s %d", time.Month(h.Idx.Month).String()[0:3], h.Idx.Year)
		}
		tplData.Title = fmt.Sprintf("Photos :: %s", tplData.Title)
		if h.Idx.Year != 0 {
			yearStr := fmt.Sprintf("%d", h.Idx.Year)
			if next := h.Idx.DB.nextMonth(yearStr, h.Idx.Month, +1); next != nil {
				tplData.Next = fmt.Sprintf("../../%d/%02d", next.Year, next.Month + 1)
			}
			if prev := h.Idx.DB.nextMonth(yearStr, h.Idx.Month, -1); prev != nil {
				tplData.Prev = fmt.Sprintf("../../%d/%02d", prev.Year, prev.Month + 1)
			}
		}
		if err := h.IndexTpl.Execute(w, tplData); err != nil {
			log.Printf("error executing template: %v\n", err)
		}
		return
	}
	if strings.HasSuffix(r.URL.Path, ".html") {
		type TplData struct {
			Title string
			Prev, Next string
			Image string
			ImgTags []string
			Tags []string
		}
		image, _ := strings.CutSuffix(r.URL.Path, ".html")
		tplData := TplData{
			Title: path.Base(h.Idx.Path),
			Next: h.Idx.Next(image, ".html"),
			Prev: h.Idx.Prev(image, ".html"),
			Image: image,
			ImgTags: h.Tags.TagsForImage(image),
			Tags: h.Tags.ShortList(),
		}
		log.Printf("%s has tags: %v\n", image, tplData.Tags)
		if h.Idx.Year != 0 {
			tplData.Title = fmt.Sprintf("%s %d", time.Month(h.Idx.Month).String()[0:3], h.Idx.Year)
		}
		tplData.Title = fmt.Sprintf("Photos :: %s :: %s", tplData.Title, image)
		if err := h.ImageTpl.Execute(w, tplData); err != nil {
			log.Printf("error executing template: %v\n", err)
		}
		return
	}
	if strings.HasSuffix(strings.ToLower(r.URL.Path), ".jpg") {
		http.ServeFile(w, r, fmt.Sprintf("%s/%s", h.Idx.Path, r.URL.Path))
		return
	}
	http.Error(w, "404 page not found", 404)
}

type MainIndexHandler struct {
	DB *ImgDB
	Tpl *template.Template
}

func (h *MainIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type TplData struct {
		Title string
		Years sort.StringSlice
		Albums sort.StringSlice
	}
	tplData := TplData{
		Title: "Photos",
		Years: make([]string, 0, len(h.DB.Years)),
		Albums: make([]string, 0, len(h.DB.Albums)),
	}
	for year := range h.DB.Years {
		tplData.Years = append(tplData.Years, year)
	}
	for album := range h.DB.Albums {
		tplData.Albums = append(tplData.Albums, album)
	}
	sort.Sort(sort.Reverse(tplData.Years))
	sort.Sort(tplData.Albums)
	if err := h.Tpl.Execute(w, tplData); err != nil {
		log.Printf("error executing template: %v\n", err)
	}
}

type AlbumIdx struct {
	DB *ImgDB
	Year int
	Month int
	Path string
	Images []string
}

func (a *AlbumIdx) indexOf(img string) int {
	for i, x := range a.Images {
		if x == img {
			return i
		}
	}
	return -1
}

func (a *AlbumIdx) next(img, suffix string, step int) string {
	i := a.indexOf(img)
	if i >= 0 {
		i += step
		if 0 <= i && i < len(a.Images) {
			return a.Images[i] + suffix
		}
	}
	if a.Year != 0 {
		if mIdx := a.DB.nextMonth(fmt.Sprintf("%d", a.Year), a.Month, step); mIdx != nil {
			i = 0
			if step < 0 {
				i = len(mIdx.Images)-1
			}
			return fmt.Sprintf("../../%d/%02d/%s.html", mIdx.Year, mIdx.Month + 1, mIdx.Images[i])
		}
	}
	return ""
}

func (a *AlbumIdx) Next(img, suffix string) string {
	return a.next(img, suffix, +1)
}

func (a *AlbumIdx) Prev(img, suffix string) string {
	return a.next(img, suffix, -1)
}

type YearIdx struct {
	DB *ImgDB
	Path string
	Months [12]*AlbumIdx
}

func (yIdx *YearIdx) next(step int) *YearIdx {
	var years []string
	for y := range yIdx.DB.Years {
		years = append(years, y)
	}
	sort.Strings(years)
	i := 0
	y0 := path.Base(yIdx.Path)
	for i < len(years) {
		if years[i] == y0 {
			break
		}
		i++
	}
	i += step
	if 0 <= i && i < len(years) {
		return yIdx.DB.Years[years[i]]
	}
	return nil
}

func (yIdx *YearIdx) Next() string {
	if next := yIdx.next(+1); next != nil {
		return path.Base(next.Path)
	}
	return ""
}

func (yIdx *YearIdx) Prev() string {
	if prev := yIdx.next(-1); prev != nil {
		return path.Base(prev.Path)
	}
	return ""
}

type ImgDB struct {
	Path string
	Years map[string]*YearIdx
	Albums map[string]*AlbumIdx
	Images map[string]struct{}
}

func (db *ImgDB) imageExists(img string) bool {
	_, exists := db.Images[img]
	return exists
}

func (db *ImgDB) nextMonth(y0 string, m0, step int) *AlbumIdx {
	var years []string
	for y := range db.Years {
		years = append(years, y)
	}
	sort.Strings(years)
	i0 := 0
	for i0 < len(years) {
		if years[i0] == y0 {
			break
		}
		i0++
	}
	for i := i0; 0 <= i && i < len(years); i += step {
		for m := m0 + step; 0 <= m && m < 12; m += step {
			if res := db.Years[years[i]].Months[m]; res != nil {
				return res
			}
		}
		if step > 0 {
			m0 = -1
		} else {
			m0 = 12
		}
	}
	return nil
}

type Templates struct {
	Main *template.Template
	Year *template.Template
	Album *template.Template
	Image *template.Template
}

func loadAlbum(db *ImgDB, year, month int, path string) (*AlbumIdx, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	albumIdx := &AlbumIdx{
		DB: db,
		Year: year,
		Month: month,
		Path: path,
		Images: make([]string, 0),
	}
	suffix := ".big.JPG"
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), suffix) {
			name, _ := strings.CutSuffix(e.Name(), suffix)
			albumIdx.Images = append(albumIdx.Images, name)
			db.Images[name] = struct{}{}
		}
	}
	return albumIdx, nil
}

func loadYear(db *ImgDB, year int, path string) (*YearIdx, error) {
	yearIdx := YearIdx{
		DB: db,
		Path: path,
	}
	for m := 0; m < 12; m++ {
		monthPath := fmt.Sprintf("%s/%02d", path, m + 1)
		if albumIdx, err := loadAlbum(db, year, m, monthPath); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("error loading album %s: %v\n", monthPath, err)
			}
		} else {
			yearIdx.Months[m] = albumIdx
		}
	}
	return &yearIdx, nil
}

func loadImageDatabase(path string, yearRanges []YearRange, albums []string) (*ImgDB, error) {
	db := &ImgDB{
		Path: path,
		Years: make(map[string]*YearIdx),
		Albums: make(map[string]*AlbumIdx),
		Images: make(map[string]struct{}),
	}
	for _, r := range yearRanges {
		curr := uint(time.Now().Year())
		for year := r.From; r.To == 0 && year <= curr || year < r.To; year++ {
			subdir := fmt.Sprintf("%s/%d", path, year)
			if _, err := os.Stat(subdir); err != nil && os.IsNotExist(err) {
				continue
			}
			if yearIdx, err := loadYear(db, int(year), subdir); err != nil {
				return nil, fmt.Errorf("loadYear: %v\n", err)
			} else {
				db.Years[fmt.Sprintf("%d", year)] = yearIdx
				log.Printf("loaded %s\n", subdir)
			}
		}
	}
	for _, album := range albums {
		subdir := fmt.Sprintf("%s/%s", path, album)
		if albumIdx, err := loadAlbum(db, 0, 0, subdir); err != nil {
			return nil, fmt.Errorf("loadAlbum: %v\n", err)
		} else {
			db.Albums[album] = albumIdx
			log.Printf("loaded %s\n", subdir)
		}
	}
	return db, nil
}

type StrLUT map[string]map[string]struct{}

func (lut StrLUT) Clr() {
	for k := range lut {
		delete(lut, k)
	}
}

func (lut StrLUT) Acc(other StrLUT) {
	for k1, obin := range other {
		bin := lut[k1]
		if bin == nil {
			bin = make(map[string]struct{})
			lut[k1] = bin
		}
		for k2 := range obin {
			bin[k2] = struct{}{}
		}
	}
}

func (lut StrLUT) Add(k, v string) {
	bin := lut[k]
	if bin == nil {
		bin = make(map[string]struct{})
	}
	bin[v] = struct{}{}
	lut[k] = bin
}

func (lut StrLUT) Del(k, v string) {
	if bin, ok := lut[k]; ok {
		delete(bin, v)
	}
}

func (lut StrLUT) Lookup(s string) []string {
	var res []string
	if bin := lut[s]; bin != nil {
		for k := range bin {
			res = append(res, k)
		}
	}
	return res
}

type Tags struct {
	sync.RWMutex
	TagLUT StrLUT
	ImgLUT StrLUT
	NewBorns StrLUT
	DeathRow StrLUT
}

func NewTags() *Tags {
	return &Tags{
		TagLUT: make(StrLUT),
		ImgLUT: make(StrLUT),
		NewBorns: make(StrLUT),
		DeathRow: make(StrLUT),
	}
}

func loadTags(tagDir, img string) (*Tags, error) {
	entries, err := os.ReadDir(fmt.Sprintf("%s/%s", tagDir, img))
	if err != nil {
		return nil, err
	}
	tags := NewTags()
	for _, e := range entries {
		if e.Type().IsRegular() {
			tags.Tag(img, e.Name())
		}
	}
	return tags, nil
}

func OpenTags(path string) (*Tags, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	tags := NewTags()
	for _, e := range entries {
		if e.IsDir() {
			if t, err := loadTags(path, e.Name()); err != nil {
				log.Printf("could not load tags for %s: %v\n", e.Name(), err)
			} else {
				tags.Acc(t)
			}
		}
	}
	return tags, nil
}

func (t *Tags) ShortList() []string {
	t.RLock()
	defer t.RUnlock()
	var tags []string
	for tag := range t.TagLUT {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func (t *Tags) Acc(u *Tags) {
	u.RLock()
	defer u.RUnlock()
	t.Lock()
	defer t.Unlock()
	t.TagLUT.Acc(u.TagLUT)
	t.ImgLUT.Acc(u.ImgLUT)
}

func (t *Tags) Tag(img, tag string) {
	t.Lock()
	defer t.Unlock()
	t.TagLUT.Add(tag, img)
	t.ImgLUT.Add(img, tag)
	t.NewBorns.Add(img, tag)
	t.DeathRow.Del(img, tag)
}

func (t *Tags) Untag(img, tag string) {
	t.Lock()
	defer t.Unlock()
	t.TagLUT.Del(tag, img)
	t.ImgLUT.Del(img, tag)
	t.NewBorns.Del(img, tag)
	t.DeathRow.Add(img, tag)
}

func (t *Tags) TagsForImage(img string) []string {
	t.RLock()
	defer t.RUnlock()
	tags := t.ImgLUT.Lookup(img)
	sort.Strings(tags)
	return tags
}

func (t *Tags) ImagesForTag(tag string) []string {
	t.RLock()
	defer t.RUnlock()
	return t.TagLUT.Lookup(tag)	
}

func (t *Tags) Flush(path string) error {
	t.Lock()
	defer t.Unlock()
	for img, tags := range t.NewBorns {
		for tag := range tags {
			imgDir := fmt.Sprintf("%s/%s", path, img)
			if err := os.MkdirAll(imgDir, 0755); err != nil {
				return err
			}
			f, err := os.Create(fmt.Sprintf("%s/%s", imgDir, tag))
			f.Close()
			if err != nil {
				return err
			}
		}
	}
	for img, tags := range t.DeathRow {
		for tag := range tags {
			if err := os.Remove(fmt.Sprintf("%s/%s/%s", path, img, tag)); err != nil {
				return err
			}
		}
	}
	return nil
}

type TagApiHandler struct {
	DB *ImgDB
	Tags *Tags
}

func (h *TagApiHandler) addTag(img, tag string) {

}

func (h *TagApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	img := r.FormValue("image")
	if !h.DB.imageExists(img) {
		http.Error(w, "not found", 404)
		return
	}
	_, delete := r.Form["delete"]
	action := "adding"
	if delete {
		action = "deleting"
	}
	for _, tag := range r.Form["tags"] {
		h.addTag(img, tag)
	}
	for _, v := range r.Form["tags"] {
		tags := strings.Split(v, " ")
		for _, tag := range tags {
			tag = strings.TrimSpace(strings.TrimPrefix(tag, "#"))
			if tag == "" {
				continue
			}
			log.Printf("%s %s tag for %s\n", action, tag, img)
			if delete {
				h.Tags.Untag(img, tag)
			} else {
				h.Tags.Tag(img, tag)
			}
			if err := h.Tags.Flush("tags"); err != nil {
				log.Printf("error updating tags dir: %v\n", err)
			}
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/%s/%s/%s.html", img[0:4], img[4:6], img), http.StatusSeeOther)
}

func loadTemplates(path string) (*Templates, error) {
	mainTpl, err := template.ParseFiles(fmt.Sprintf("%s/main.tpl", path))
	if err != nil {
		return nil, err
	}
	yearTpl, err := template.ParseFiles(fmt.Sprintf("%s/year.tpl", path))
	if err != nil {
		return nil, err
	}
	albumTpl, err := template.ParseFiles(fmt.Sprintf("%s/album.tpl", path))
	if err != nil {
		return nil, err
	}
	imageTpl, err := template.ParseFiles(fmt.Sprintf("%s/image.tpl", path))
	if err != nil {
		return nil, err
	}
	return &Templates{
		Main: mainTpl,
		Year: yearTpl,
		Album: albumTpl,
		Image: imageTpl,
	}, nil
}

type YearRange struct {
	From uint
	To uint
}

func main() {
	yearRanges := []YearRange{
		{From: 2008},
	}
	albums := []string{
		"misc",
	}
	templates, err := loadTemplates(".")
	if err != nil {
		log.Fatalf("could not load templates: %v\n", err)
	}
	db, err := loadImageDatabase(".", yearRanges, albums)
	if err != nil {
		log.Fatalf("could not load database: %v\n", err)
	}
	tags, err := OpenTags("tags")
	if err != nil {
		log.Fatalf("could not load tags: %v\n", err)
	}
	http.Handle("/api/tag", &TagApiHandler{db, tags})
	for y, yIdx := range db.Years {
		for m, mIdx := range yIdx.Months {
			if mIdx != nil {
				prefix := fmt.Sprintf("/%s/%02d/", y, m+1)
				http.Handle(prefix, http.StripPrefix(prefix, &AlbumIndexHandler{mIdx, templates.Album, templates.Image, tags}))
			}
		}
		prefix := fmt.Sprintf("/%s/", y)
		http.Handle(prefix, http.StripPrefix(prefix, &YearIndexHandler{yIdx, templates.Year}))
	}
	for album, idx := range db.Albums {
		prefix := fmt.Sprintf("/%s/", album)
		http.Handle(prefix, http.StripPrefix(prefix, &AlbumIndexHandler{idx, templates.Album, templates.Image, tags}))
	}
	http.Handle("/", &MainIndexHandler{db, templates.Main})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
