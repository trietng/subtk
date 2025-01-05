package download

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"trietng/subtk/common"
	"trietng/subtk/common/utils"
	"trietng/subtk/download/errmsg"
	"trietng/subtk/match"
)

type SubtitleDownloader struct {
	Dir string
	ExtractArchive bool
	Url            string
}

// only support zip for now
func (d *SubtitleDownloader) extractArchive(r io.Reader, mis []common.MediaInfo) error {
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, r)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buff.Bytes())
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		if err == zip.ErrFormat {
			fmt.Println(errmsg.WarnInvalidFileFormat)
			return nil
		}
		return err
	}
	zMis := make(map[utils.Pair[int, int]]*zip.File)
	for _, file := range zipReader.File {
		season, episode, err := match.GetSeasonAndEpisode(file.Name)
		if err != nil {
			continue
		}
		zMis[utils.Pair[int, int]{First: season, Second: episode}] = file
	}
	for _, mi := range mis {
		if file, ok := zMis[utils.Pair[int, int]{First: mi.Season, Second: mi.Episode}]; ok {
			// extract the file
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			ext := filepath.Ext(file.Name)
			subtitleFileName := fmt.Sprintf("%s%s", mi.Title, ext)
			fmt.Println(subtitleFileName)
			out, err := os.OpenFile(fmt.Sprintf("%s/%s", d.Dir, subtitleFileName), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *SubtitleDownloader) Download() error {
	resp, err := http.Get(d.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if d.ExtractArchive {
		matcher := match.SubtitleMatcher{}
		mis, err := matcher.MatchFiles(d.Dir)
		if err != nil {
			return err
		}
		return d.extractArchive(resp.Body, mis)
	}
	return nil
}