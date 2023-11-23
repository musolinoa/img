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
	"time"
)

type YearIndexHandler struct {
	Idx *YearIdx
	Tpl *template.Template
}

func (h *YearIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("YearIndexHandler.ServeHTTP: %s\n", r.URL.Path)
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
		PrevYear string
		CurrYear string
		NextYear string
		Months [12]MonthTplData
	}
	tplData := TplData{
		CurrYear: path.Base(h.Idx.Path),
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
}

func (h *AlbumIndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("AlbumIndexHandler.ServeHTTP: %s\n", r.URL.Path)
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
			yearStr := fmt.Sprintf("%d", h.Idx.Year)
			if next := h.Idx.DB.nextMonth(yearStr, h.Idx.Month, +1); next != "" {
				tplData.Next = "../../" + next
			}
			if prev := h.Idx.DB.nextMonth(yearStr, h.Idx.Month, -1); prev != "" {
				tplData.Prev = "../../" + prev
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
		}
		image, _ := strings.CutSuffix(r.URL.Path, ".html")
		tplData := TplData{
			Title: path.Base(h.Idx.Path),
			Next: h.Idx.Next(image, ".html"),
			Prev: h.Idx.Prev(image, ".html"),
			Image: image,
		}
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
	log.Printf("IndexHandler.ServeHTTP\n")
	type TplData struct {
		Years sort.StringSlice
		Albums sort.StringSlice
	}
	tplData := TplData{
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

type ImgDB struct {
	Path string
	Years map[string]*YearIdx
	Albums map[string]*AlbumIdx
}

func (db *ImgDB) nextMonth(y0 string, m0, step int) string {
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
		log.Printf("y=%s, m0=%d, i=%d\n", years[i], m0+1, i)
		for m := m0 + step; 0 <= m && m < 12; m += step {
			if db.Years[years[i]].Months[m] != nil {
				return fmt.Sprintf("%s/%02d", years[i], m + 1)
			}
		}
		if step > 0 {
			m0 = -1
		} else {
			m0 = 12
		}
	}
	return ""
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

func loadDatabase(path string, yearRanges []YearRange, albums []string) (*ImgDB, error) {
	db := &ImgDB{
		Path: path,
		Years: make(map[string]*YearIdx),
		Albums: make(map[string]*AlbumIdx),
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
	db, err := loadDatabase(".", yearRanges, albums)
	if err != nil {
		log.Fatalf("could not load database: %v\n", err)
	}
	for y, yIdx := range db.Years {
		for m, mIdx := range yIdx.Months {
			if mIdx != nil {
				prefix := fmt.Sprintf("/%s/%02d/", y, m+1)
				log.Printf("adding handler for %s\n", prefix)
				http.Handle(prefix, http.StripPrefix(prefix, &AlbumIndexHandler{mIdx, templates.Album, templates.Image}))
			}
		}
		prefix := fmt.Sprintf("/%s/", y)
		http.Handle(prefix, http.StripPrefix(prefix, &YearIndexHandler{yIdx, templates.Year}))
	}
	for album, idx := range db.Albums {
		prefix := fmt.Sprintf("/%s/", album)
		http.Handle(prefix, http.StripPrefix(prefix, &AlbumIndexHandler{idx, templates.Album, templates.Image}))
	}
	http.Handle("/", &MainIndexHandler{db, templates.Main})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
