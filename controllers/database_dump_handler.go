package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/metainfo"
	"github.com/gin-gonic/gin"
)

const (
	// DatabaseDumpPath : Location of database dumps
	DatabaseDumpPath = "./public/dumps"
	// GPGPublicKeyPath : Location of the GPG key
	GPGPublicKeyPath = "./public/gpg/gpg.key"
)

// DatabaseDumpHandler : Controller for getting database dumps
func DatabaseDumpHandler(c *gin.Context) {
	// db params url
	// TODO Use config from cli
	files, _ := filepath.Glob(filepath.Join(DatabaseDumpPath, "*.torrent"))

	var dumpsJSON []model.DatabaseDumpJSON
	// TODO Filter *.torrent files
	for _, f := range files {
		// TODO Use config from cli
		file, err := os.Open(f)
		if err != nil {
			continue
		}
		var tf metainfo.TorrentFile
		err = tf.Decode(file)
		if err != nil {
			log.CheckError(err)
			fmt.Println(err)
			continue
		}
		dump := model.DatabaseDump{
			Date:        time.Now(),
			Filesize:    int64(tf.TotalSize()),
			Name:        tf.TorrentName(),
			TorrentLink: "/dbdumps/" + file.Name()}
		dumpsJSON = append(dumpsJSON, dump.ToJSON())
	}

	databaseDumpTemplate(c, dumpsJSON, "/gpg/gpg.pub")
}
